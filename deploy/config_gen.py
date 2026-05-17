#!/usr/bin/env python3
"""
datauptwo 配置生成器
从 deploy/config.yaml 统一配置，生成各环境专用配置文件
密钥自动生成
"""

import os
import sys
import yaml
import secrets
import argparse
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent
ROOT_DIR = SCRIPT_DIR.parent


def gen_secret(length=32):
    """生成安全的随机密钥"""
    alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    return ''.join(secrets.choice(alphabet) for _ in range(length))


def load_config():
    path = SCRIPT_DIR / "config.yaml"
    if not path.exists():
        print(f"[ERROR] 配置文件不存在: {path}")
        sys.exit(1)
    with open(path, encoding="utf-8") as f:
        return yaml.safe_load(f)


def gen_backend_config(cfg):
    """生成后端 app.yaml"""
    app = cfg["app"]
    db = cfg["database"]
    sec = cfg["security"]
    app_name = app["name"]

    # 找到启用的数据库
    enabled_db = None
    for db_type in ["sqlite", "mysql", "postgres"]:
        if db.get(db_type, {}).get("enabled"):
            enabled_db = (db_type, db[db_type])
            break
    if enabled_db is None:
        print("[WARN] 未找到启用的数据库配置，默认为 sqlite")
        enabled_db = ("sqlite", db["sqlite"])

    db_type, db_cfg = enabled_db

    # 动态生成密钥
    jwt_secret = gen_secret()
    session_secret = gen_secret()

    lines = [
        f"# ============================================================",
        f"# {app_name} 后端配置（自动生成）",
        f"# 来源: deploy/config.yaml",
        f"# ============================================================",
        f"",
        f"base:",
        f"  mode: {app['mode']}",
        f"  port: {app['port']}",
        f"  username: admin",
        f"  password: {sec['admin_password']}",
        f"  version: {app['version']}",
        f"  language: zh",
        f"  install_dir: ./data",
        f"  serve_web: {str(app['serve_web']).lower()}",
        f"  web_root: ./web",
        f"",
        f"database:",
        f"  type: {db_type}",
    ]

    if db_type == "sqlite":
        lines.append(f"  path: {db_cfg['path']}")
    elif db_type == "mysql":
        lines.extend([
            f"  host: {db_cfg['host']}",
            f"  port: {db_cfg['port']}",
            f"  user: {db_cfg['user']}",
            f"  pass: {db_cfg['pass']}",
            f"  name: {db_cfg['name']}",
        ])
    elif db_type == "postgres":
        lines.extend([
            f"  host: {db_cfg['host']}",
            f"  port: {db_cfg['port']}",
            f"  user: {db_cfg['user']}",
            f"  pass: {db_cfg['pass']}",
            f"  name: {db_cfg['name']}",
        ])

    lines.extend([
        f"",
        f"log:",
        f"  level: {cfg['logging']['level']}",
        f"  time_zone: {app['timezone']}",
        f"  log_name: {cfg['logging']['log_name']}",
        f"  max_backup: {cfg['logging']['max_backup']}",
        f"  max_size: {cfg['logging']['max_size']}",
        f"",
        f"session:",
        f"  timeout: {sec['session_timeout']}",
        f"  secret: {session_secret}",
        f"",
        f"cors:",
        f"  allow_origins:",
    ])
    for origin in cfg["cors"]["allow_origins"]:
        lines.append(f'    - "{origin}"')

    lines.extend([
        f"  allow_methods:",
    ])
    for method in cfg["cors"]["allow_methods"]:
        lines.append(f"    - {method}")

    lines.extend([
        f"  allow_headers:",
    ])
    for header in cfg["cors"]["allow_headers"]:
        lines.append(f"    - {header}")

    lines.extend([
        f"  max_age: {cfg['cors']['max_age']}",
        f"",
        f"jwt:",
        f"  secret: {jwt_secret}",
    ])

    content = "\n".join(lines) + "\n"
    out = ROOT_DIR / "backend" / "conf" / "app.yaml"
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"[OK]   后端配置: {out}")
    print(f"[INFO] JWT_SECRET  = {jwt_secret}")
    print(f"[INFO] SESSION_SECRET = {session_secret}")


def gen_docker_compose(cfg):
    """生成 docker-compose.yml"""
    db = cfg["database"]
    dep = cfg["deploy"]
    app = cfg["app"]
    app_name = app["name"]

    # 找到启用的数据库
    enabled_db = None
    svc_db = {}
    for db_type in ["mysql", "postgres"]:
        if db.get(db_type, {}).get("enabled"):
            enabled_db = db_type
            svc_db = db[db_type]
            break

    # 服务定义
    app_svc = f"""  app:
    build:
      context: ..
      dockerfile: deploy/docker/Dockerfile
    container_name: {app_name}-app
    restart: unless-stopped
    ports:
      - "{dep['external_port']}:{app['port']}"
    environment:
      - TZ={app['timezone']}
    volumes:
      - app-data:/app/data
      - app-logs:/app/logs
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:{app['port']}/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - {app_name}
"""

    nginx_svc = f"""  nginx:
    image: nginx:1.27-alpine
    container_name: {app_name}-nginx
    restart: unless-stopped
    ports:
      - "{dep['http_port']}:80"
      - "{dep['https_port']}:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - app
    networks:
      - {app_name}"""

    mysql_svc = ""
    postgres_svc = ""
    app_depends = ""
    if enabled_db == "mysql":
        app_depends = "\n    depends_on:\n      mysql:\n        condition: service_healthy"
        mysql_svc = f"""
  mysql:
    image: mysql:8.0
    container_name: {app_name}-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: {svc_db['pass']}
      MYSQL_DATABASE: {svc_db['name']}
      MYSQL_USER: {svc_db['user']}
      MYSQL_PASSWORD: {svc_db['pass']}
      TZ: {app['timezone']}
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
    networks:
      - {app_name}
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
      - --default-authentication-plugin=mysql_native_password
      - --skip-name-resolve"""
    elif enabled_db == "postgres":
        app_depends = "\n    depends_on:\n      postgres:\n        condition: service_healthy"
        postgres_svc = f"""
  postgres:
    image: postgres:16-alpine
    container_name: {app_name}-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: {svc_db['user']}
      POSTGRES_PASSWORD: {svc_db['pass']}
      POSTGRES_DB: {svc_db['name']}
      TZ: {app['timezone']}
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {svc_db['user']}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - {app_name}"""

    # app 服务加上 depends
    if app_depends:
        app_svc = app_svc.replace(
            "    networks:\n      - " + app_name,
            app_depends + "\n    networks:\n      - " + app_name
        )

    volumes = f"""volumes:
  app-data:
    driver: local
  app-logs:
    driver: local
"""
    if enabled_db == "mysql":
        volumes += "  mysql-data:\n    driver: local\n"
    elif enabled_db == "postgres":
        volumes += "  postgres-data:\n    driver: local\n"

    content = f"""# ============================================================
# {app_name} Docker Compose（自动生成）
# 来源: deploy/config.yaml
# 数据库: {enabled_db or 'sqlite'}
# ============================================================

services:
{app_svc}
{nginx_svc}{mysql_svc}{postgres_svc}

{volumes}
networks:
  {app_name}:
    driver: bridge
"""

    out = SCRIPT_DIR / "docker" / "docker-compose.yml"
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"[OK]   Docker Compose: {out}")


def gen_nginx_config(cfg):
    """生成 Nginx 配置"""
    app = cfg["app"]
    dep = cfg["deploy"]

    content = f"""# ============================================================
# Nginx 反向代理配置（自动生成）
# ============================================================

worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {{
    worker_connections 1024;
    multi_accept on;
    use epoll;
}}

http {{
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for" '
                    'rt=$request_time';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_min_length 1024;
    gzip_types
        text/plain text/css text/xml text/javascript application/json
        application/javascript application/xml application/xml+rss;

    server_tokens off;

    upstream {app['name']} {{
        server app:{app['port']};
        keepalive 32;
    }}

    server {{
        listen 80;
        server_name _;

        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        location /health {{
            proxy_pass http://{app['name']}/health;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            access_log off;
        }}

        location /swagger {{
            proxy_pass http://{app['name']}/swagger;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }}

        location /api/ {{
            proxy_pass http://{app['name']}/api/;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            client_max_body_size 500m;
        }}

        location /public/ {{
            proxy_pass http://{app['name']}/public/;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            client_max_body_size 500m;
        }}

        location / {{
            proxy_pass http://{app['name']}/;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }}

        location ~* \\.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {{
            proxy_pass http://{app['name']};
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            expires 7d;
            add_header Cache-Control "public, immutable";
            access_log off;
        }}

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {{
            root /usr/share/nginx/html;
        }}
    }}
}}
"""

    out = SCRIPT_DIR / "docker" / "nginx" / "nginx.conf"
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"[OK]   Nginx 配置: {out}")


def gen_k8s_manifests(cfg):
    """生成 Kubernetes manifests"""
    app = cfg["app"]
    db = cfg["database"]
    sec = cfg["security"]
    dep = cfg["deploy"]
    app_name = app["name"]

    # 找到启用的数据库
    enabled_db = None
    svc_db = {}
    for db_type in ["mysql", "postgres", "sqlite"]:
        if db.get(db_type, {}).get("enabled"):
            enabled_db = db_type
            svc_db = db[db_type]
            break
    if enabled_db is None:
        enabled_db = "sqlite"

    db_host_map = {"mysql": "mysql-svc", "postgres": "postgres-svc", "sqlite": ""}
    db_host = db_host_map.get(enabled_db, "")
    db_pass = svc_db.get("pass", "") if enabled_db != "sqlite" else ""

    jwt_secret = gen_secret()
    session_secret = gen_secret()

    content = f"""# ============================================================
# {app_name} Kubernetes 部署清单（自动生成）
# 来源: deploy/config.yaml
# 数据库: {enabled_db}
# ============================================================

---
apiVersion: v1
kind: Namespace
metadata:
  name: {app_name}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {app_name}-config
  namespace: {app_name}
data:
  app.yaml: |
    base:
      mode: prod
      port: {app['port']}
      version: {app['version']}
      language: zh
      install_dir: ./data
      serve_web: {str(app['serve_web']).lower()}
      web_root: ./web

    database:
      type: {enabled_db}
"""
    if enabled_db == "sqlite":
        content += f"      path: ./data/registry.db\n"
    else:
        content += f"      host: {db_host}\n      port: {svc_db['port']}\n      user: {svc_db['user']}\n      pass: {svc_db['pass']}\n      name: {svc_db['name']}\n"

    content += f"""
    log:
      level: {cfg['logging']['level']}
      time_zone: {app['timezone']}
      log_name: {cfg['logging']['log_name']}
      max_backup: {cfg['logging']['max_backup']}
      max_size: {cfg['logging']['max_size']}

    session:
      timeout: {sec['session_timeout']}
      secret: {session_secret}

    cors:
      allow_origins:
"""
    for origin in cfg["cors"]["allow_origins"]:
        content += f'        - "{origin}"\n'

    content += f"""      allow_methods:
        - GET
        - POST
        - PUT
        - DELETE
        - OPTIONS
      allow_headers:
        - Origin
        - Content-Type
        - Authorization
      max_age: {cfg['cors']['max_age']}

    jwt:
      secret: {jwt_secret}

---
apiVersion: v1
kind: Secret
metadata:
  name: {app_name}-secret
  namespace: {app_name}
type: Opaque
stringData:
  DB_PASSWORD: {db_pass}
  SESSION_SECRET: {session_secret}
  JWT_SECRET: {jwt_secret}

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {app_name}-data-pvc
  namespace: {app_name}
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {dep['storage_size']}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {app_name}
  namespace: {app_name}
spec:
  replicas: {dep['replicas']}
  selector:
    matchLabels:
      app: {app_name}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: {app_name}
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        runAsGroup: 1000
        fsGroup: 1000
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - {app_name}
                topologyKey: kubernetes.io/hostname
      containers:
        - name: app
          image: {dep['image_registry']}/{dep['image_name']}:{app['version']}
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: {app['port']}
          env:
            - name: TZ
              value: {app['timezone']}
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: {app_name}-secret
                  key: DB_PASSWORD
            - name: SESSION_SECRET
              valueFrom:
                secretKeyRef:
                  name: {app_name}-secret
                  key: SESSION_SECRET
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: {app_name}-secret
                  key: JWT_SECRET
          volumeMounts:
            - name: config
              mountPath: /app/conf/app.yaml
              subPath: app.yaml
              readOnly: true
            - name: data
              mountPath: /app/data
            - name: logs
              mountPath: /app/logs
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              cpu: 1000m
              memory: 1Gi
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 15
            periodSeconds: 15
          readinessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
      volumes:
        - name: config
          configMap:
            name: {app_name}-config
        - name: data
          persistentVolumeClaim:
            claimName: {app_name}-data-pvc
        - name: logs
          emptyDir: {{}}
      terminationGracePeriodSeconds: 30

---
apiVersion: v1
kind: Service
metadata:
  name: {app_name}-svc
  namespace: {app_name}
spec:
  type: ClusterIP
  ports:
    - name: http
      port: {app['port']}
      targetPort: http
  selector:
    app: {app_name}

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {app_name}-ingress
  namespace: {app_name}
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
    nginx.ingress.kubernetes.io/proxy-body-size: "500m"
spec:
  ingressClassName: nginx
  rules:
    - host: {dep['domain']}
      http:
        paths:
          - path: /swagger
            pathType: Prefix
            backend:
              service:
                name: {app_name}-svc
                port:
                  number: {app['port']}
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: {app_name}-svc
                port:
                  number: {app['port']}
          - path: /public
            pathType: Prefix
            backend:
              service:
                name: {app_name}-svc
                port:
                  number: {app['port']}
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {app_name}-svc
                port:
                  number: {app['port']}
"""

    out = SCRIPT_DIR / "kubernetes" / "manifests.yaml"
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"[OK]   K8s 清单: {out}")


def gen_env_file(cfg):
    """生成 .env 文件"""
    db = cfg["database"]
    app = cfg["app"]

    mysql_cfg = db.get("mysql", {})
    postgres_cfg = db.get("postgres", {})

    lines = [
        f"# {app['name']} Docker 环境变量（自动生成）",
        f"TZ={app['timezone']}",
        f"VERSION={app['version']}",
    ]
    if mysql_cfg.get("enabled"):
        lines.append(f"MYSQL_ROOT_PASSWORD={mysql_cfg['pass']}")
        lines.append(f"MYSQL_DB_PASSWORD={mysql_cfg['pass']}")
    if postgres_cfg.get("enabled"):
        lines.append(f"POSTGRES_PASSWORD={postgres_cfg['pass']}")

    content = "\n".join(lines) + "\n"
    out = SCRIPT_DIR / "docker" / ".env"
    with open(out, "w", encoding="utf-8") as f:
        f.write(content)
    print(f"[OK]   .env: {out}")


def main():
    parser = argparse.ArgumentParser(description="datauptwo 配置生成器")
    parser.add_argument("--all", action="store_true", help="生成所有配置")
    parser.add_argument("--backend", action="store_true", help="生成后端配置")
    parser.add_argument("--docker", action="store_true", help="生成 Docker Compose")
    parser.add_argument("--nginx", action="store_true", help="生成 Nginx 配置")
    parser.add_argument("--k8s", action="store_true", help="生成 K8s 清单")
    parser.add_argument("--env", action="store_true", help="生成 .env 文件")
    args = parser.parse_args()

    if not any(vars(args).values()):
        args.all = True

    cfg = load_config()

    if args.all or args.backend:
        gen_backend_config(cfg)
    if args.all or args.docker:
        gen_docker_compose(cfg)
    if args.all or args.nginx:
        gen_nginx_config(cfg)
    if args.all or args.k8s:
        gen_k8s_manifests(cfg)
    if args.all or args.env:
        gen_env_file(cfg)

    print()
    print("=" * 50)
    print("配置生成完成！")
    print("=" * 50)


if __name__ == "__main__":
    main()
