# duptwo 数据登记平台 - 部署指南

> **产品名称**: duptwo / 数据登记平台
> **版本**: v2.0
> **最后更新**: 2026-03-27

## 目录

- [快速开始](#快速开始)
- [系统架构](#系统架构)
- [环境要求](#环境要求)
- [配置说明](#配置说明)
- [单机部署](#单机部署)
- [Docker 部署](#docker-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [Nginx 配置](#nginx-配置)
- [数据库迁移](#数据库迁移)
- [生产环境检查清单](#生产环境检查清单)
- [常见问题](#常见问题)

---

## 快速开始

### 一键启动（Docker，推荐）

```bash
# 进入部署目录
cd deploy

# 方式1：MySQL 版（推荐生产环境）
./deploy.sh deploy

# 方式2：SQLite 版（轻量，适合小规模 / 开发测试）
./deploy.sh deploy:sqlite
```

启动后访问：
- **Web UI**: http://localhost:8080
- **Swagger API 文档**: http://localhost:8080/swagger
- **健康检查**: http://localhost:8080/health
- **默认账号**: `admin` / `admin123`

### 开发模式

```bash
# 后端（端口 18421）
cd backend && go run cmd/server/main.go

# 前端（端口 4004，自动代理 API 到 18421）
cd frontend && npm run dev
```

---

## 系统架构

```
                           ┌─────────────────┐
                           │   用户浏览器      │
                           └────────┬────────┘
                                    │ HTTPS :443 / HTTP :80
                           ┌────────▼────────┐
                           │     Nginx       │  ← 反向代理 + SSL + 静态缓存
                           │    (可选)        │    代理到 :8080
                           └────────┬────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
        │                           │                           │
┌──────▼───────┐          ┌──────▼───────┐          ┌─────▼─────┐
│   前端静态文件  │          │   Go API      │          │   MySQL   │
│   (内置于Go)   │          │   :8080       │          │   :3306   │
│               │          │               │          │           │
│  /            │          │  /api/*       │          └───────────┘
│  /swagger     │          │  /public/*    │
│  (SPA路由兜底) │          │  /health      │
└───────────────┘          └───────────────┘
```

**架构说明**：
- Go 服务内置静态文件服务（`serve_web: true`），可直接提供前端页面
- Nginx 属于可选项：用于 HTTPS 终结、负载均衡、静态资源缓存
- SQLite 版无需 MySQL，所有数据存储在单个文件 `./data/registry.db`

---

## 环境要求

### 基础环境

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Docker | 20.10 | 24.0+ |
| Docker Compose | 2.0 | 2.20+ |
| Go | 1.21 | 1.23 |
| Node.js | 18.0 | 20.0 LTS |
| MySQL | 8.0 | 8.0 LTS |
| Python3 | 3.9 | 3.11 |

### 硬件要求

| 场景 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发/测试 | 2核 | 4GB | 20GB |
| 小规模生产 | 2核 | 4GB | 50GB |
| 中等规模生产 | 4核 | 8GB | 100GB |
| 大规模生产 | 8核+ | 16GB+ | 200GB+ |

---

## 配置说明

### 配置文件位置

- 开发配置: `backend/conf/app.yaml`
- 生产配置: `deploy/docker/conf/app.prod.yaml`
- Docker 环境变量: `deploy/docker/.env`
- 一键部署配置: `deploy/config.yaml`

### 基础配置 (`conf/app.yaml`)

```yaml
base:
  mode: prod                # dev | prod
  port: 8080                # 监听端口
  serve_web: true          # true=提供前端，false=仅API
  log_level: info          # debug | info | warn | error

database:
  type: mysql              # mysql | sqlite | postgres
  host: localhost
  port: 3306
  user: root
  password: your_password
  name: duptwo
  max_open_conns: 100
  max_idle_conns: 10

jwt:
  secret: CHANGE_THIS_SECRET_KEY_HERE  # ⚠️ 必须修改
  expire_hours: 168             # 7天
  refresh_expire_hours: 720     # 30天

ldap:
  enabled: false
  host: ldap.example.com
  port: 389
  base_dn: dc=example,dc=com
  bind_dn: cn=admin,dc=example,dc=com
  bind_password: ""
  user_filter: "(sAMAccountName=%s)"
  mappings:
    username: sAMAccountName
    nickname: displayName
    email: mail
    phone: telephoneNumber

security:
  mfa:
    enabled: true           # 启用 MFA
    issuer: duptwo           # TOTP 发行方名称
  password:
    min_length: 6
    require_uppercase: false
    require_numbers: false
    require_special: false
  lockout:
    enabled: true
    max_failed_attempts: 5
    lockout_duration: 30m
```

### Docker 环境变量 (`.env`)

```bash
# 服务配置
IMAGE_NAME=duptwo
IMAGE_TAG=latest
API_PORT=8080
VOLUME_PATH=./data

# 数据库配置（MySQL 版）
DB_TYPE=mysql
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_secure_password
DB_NAME=duptwo

# JWT 配置
JWT_SECRET=your_very_long_random_secret_here

# 可选：反向代理
ENABLE_NGINX=true
SSL_EMAIL=admin@example.com
```

---

## 单机部署

### 方式一：二进制部署

```bash
# 1. 构建后端
cd backend
make build

# 2. 初始化数据库（自动创建表）
./server init-db

# 3. 启动服务
./server --config=../conf/app.yaml

# 4. 构建前端并部署到 Go 服务
cd ../frontend
npm run build   # 输出到 ../backend/cmd/server/web
```

### 方式二：纯 Docker 单机

```bash
cd deploy

# 快速启动（SQLite，无需 MySQL）
docker compose -f docker/docker-compose.sqlite.yml up -d

# 或 MySQL 版
docker compose -f docker/docker-compose.yml up -d
```

---

## Docker 部署

### 项目结构

```
deploy/
├── docker/
│   ├── Dockerfile           # 多阶段构建（Go + 前端）
│   ├── docker-compose.yml   # MySQL 版编排
│   ├── docker-compose.sqlite.yml  # SQLite 版编排
│   ├── nginx/
│   │   └── nginx.conf      # Nginx 配置
│   └── conf/
│       └── app.prod.yaml   # 生产配置
├── kubernetes/
│   ├── namespace.yaml
│   └── manifests.yaml      # 所有 K8s 资源
├── config.yaml             # 一键部署配置（统一入口）
├── config_gen.py           # 配置生成器
└── deploy.sh               # 一键部署脚本
```

### 使用一键部署脚本

```bash
cd deploy

# 完整部署（构建 + 启动 MySQL + 启动服务）
./deploy.sh deploy

# 仅启动已构建的服务
./deploy.sh start

# 停止服务
./deploy.sh stop

# 查看日志
./deploy.sh logs

# 进入容器
./deploy.sh shell

# 重建并部署
./deploy.sh rebuild

# 完全清理（包括数据卷）
./deploy.sh clean
```

### 自定义构建

```bash
# 1. 编辑配置
vim config.yaml

# 2. 生成各环境配置
python3 config_gen.py --all

# 3. 构建镜像
./deploy.sh build

# 4. 启动
./deploy.sh deploy
```

### MySQL 版 (`docker-compose.yml`)

包含服务：
- `app`: Go 主服务 (:8080)
- `mysql`: MySQL 8.0 (:3306)
- `nginx`: 反向代理 (:80, :443，可选)

```bash
# 启动
docker compose -f docker/docker-compose.yml up -d

# 查看状态
docker compose -f docker/docker-compose.yml ps

# 查看日志
docker compose -f docker/docker-compose.yml logs -f app
```

### SQLite 版 (`docker-compose.sqlite.yml`)

轻量级，单容器，无需 MySQL。数据存储在 `data/registry.db`。

```bash
# 启动
docker compose -f docker/docker-compose.sqlite.yml up -d

# 数据备份（宿主机执行）
cp data/registry.db data/registry.db.$(date +%Y%m%d).bak
```

---

## Kubernetes 部署

### 前置要求

- Kubernetes 1.25+
- kubectl 已配置集群
- PV/PVC 支持（建议使用云存储）
- Ingress Controller（如 Nginx Ingress）

### 快速部署

```bash
cd deploy/kubernetes

# 编辑镜像地址（替换为你的镜像仓库）
sed -i 's|image: your-registry.com/duptwo:v2.0|your-registry.com/duptwo:latest|g' manifests.yaml

# 应用清单
kubectl apply -f namespace.yaml
kubectl apply -f manifests.yaml

# 检查状态
kubectl -n duptwo get all
kubectl -n duptwo get ingress
```

### 手动分步部署

```bash
# 1. 创建命名空间
kubectl apply -f namespace.yaml

# 2. 创建 ConfigMap（配置）
kubectl -n duptwo create configmap duptwo-config \
  --from-file=app.yaml=./conf/app.prod.yaml

# 3. 创建 Secret（密钥，请先 base64 编码）
echo -n 'your-jwt-secret' | base64
kubectl -n duptwo create secret generic duptwo-secrets \
  --from-literal=jwt-secret=your_base64_encoded_secret \
  --from-literal=db-password=your_db_password

# 4. 部署数据库（可选，使用云数据库时跳过）
kubectl apply -f mysql.yaml

# 5. 部署应用
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

### Kubernetes 资源说明 (`manifests.yaml`)

| 资源类型 | 名称 | 说明 |
|---------|------|------|
| Namespace | duptwo | 独立命名空间 |
| ConfigMap | duptwo-config | 应用配置 |
| Secret | duptwo-secrets | 密钥（JWT、DB 密码等） |
| Deployment | duptwo-app | 应用副本集 |
| Service | duptwo-svc | ClusterIP 服务 |
| Ingress | duptwo-ingress | HTTP 入口 |
| HorizontalPodAutoscaler | duptwo-hpa | 自动扩缩容（可选） |

### 资源配额建议

```yaml
resources:
  requests:
    cpu: 100m
    memory: 256Mi
  limits:
    cpu: 2000m
    memory: 2Gi
```

### 健康检查

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 20

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

---

## Nginx 配置

### 基础反向代理

```nginx
upstream duptwo_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name duptwo.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name duptwo.example.com;

    ssl_certificate     /etc/ssl/certs/duptwo.crt;
    ssl_certificate_key /etc/ssl/private/duptwo.key;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # API 代理
    location /api/ {
        proxy_pass         http://duptwo_backend;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 60s;
    }

    # 静态文件（前端已内置在 Go 服务中）
    location / {
        proxy_pass http://duptwo_backend;
        proxy_set_header Host $host;
    }

    # Swagger 文档
    location /swagger/ {
        proxy_pass http://duptwo_backend;
        proxy_set_header Host $host;
    }
}
```

### 高可用 Nginx 配置

使用 `nginx.conf` 中的 `upstream` 块定义多后端实例，配合 keepalived 实现 VIP 漂移。

---

## 数据库迁移

### 首次启动

服务启动时自动执行 GORM AutoMigrate，无需手动迁移。

### 生产数据库初始化（MySQL）

```bash
# 1. 创建数据库
mysql -h localhost -u root -p -e "CREATE DATABASE IF NOT EXISTS duptwo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. 创建专用用户（推荐）
mysql -h localhost -u root -p -e "
CREATE USER IF NOT EXISTS 'duptwo'@'%' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON duptwo.* TO 'duptwo'@'%';
FLUSH PRIVILEGES;
"

# 3. 启动服务（GORM 自动建表）
docker compose -f docker/docker-compose.yml up -d app
```

### 数据备份

**MySQL:**
```bash
mysqldump -h localhost -u duptwo -p duptwo > backup_$(date +%Y%m%d).sql

# 恢复
mysql -h localhost -u duptwo -p duptwo < backup_20260327.sql
```

**SQLite:**
```bash
# 备份
cp data/registry.db registry.db.$(date +%Y%m%d).bak

# 压缩备份
tar czf duptwo_backup_$(date +%Y%m%d).tar.gz data/registry.db
```

### 定时备份（Cron）

```bash
# 每天凌晨 3 点备份
0 3 * * * cd /opt/duptwo && docker exec duptwo-mysql mysqldump -u root -p'pass' duptwo > backup/$(date +\%Y\%m\%d).sql

# 每周一清理 30 天前的备份
0 4 * * 1 find /opt/duptwo/backups -name "*.sql" -mtime +30 -delete
```

---

## 生产环境检查清单

部署前请逐项确认：

### 安全配置
- [ ] JWT Secret 已修改为强随机字符串（32+ 字符）
- [ ] Admin 默认密码已修改
- [ ] MySQL Root 密码已修改
- [ ] CORS 白名单已配置（`conf/app.yaml` → `cors.allowed_origins`）
- [ ] 数据库用户权限已限制（不使用 root 远程连接）
- [ ] 生产环境关闭 `mode: prod`

### 网络配置
- [ ] 配置了 HTTPS（Let's Encrypt 或商业证书）
- [ ] Nginx 反向代理已配置
- [ ] 防火墙已开放必要端口（80, 443）
- [ ] 数据库端口（3306）仅内网访问

### 监控告警
- [ ] 配置了日志收集（ELK / Loki）
- [ ] 配置了指标监控（Prometheus + Grafana）
- [ ] 配置了健康检查告警
- [ ] 配置了磁盘空间告警

### 数据安全
- [ ] 配置了定时备份策略
- [ ] 备份存储与生产分离
- [ ] 数据库开启慢查询日志

### 性能优化
- [ ] 启用 Nginx Gzip 压缩
- [ ] 数据库连接池已调优（`max_open_conns`）
- [ ] 前端静态资源设置了缓存头
- [ ] 大文件上传配置了合理的超时时间

---

## 常见问题

### Q: 启动报 `port already in use`

```bash
# 查看占用端口的进程
lsof -i :8080

# 杀掉占用进程或修改配置中的端口
```

### Q: MySQL 连接失败

```bash
# 检查容器网络
docker network ls
docker compose -f docker/docker-compose.yml exec mysql mysql -u root -p

# 检查连接字符串
docker compose -f docker/docker-compose.yml exec app env | grep DB_
```

### Q: 前端构建失败

```bash
cd frontend
npm install
npm run build
```

### Q: Swagger 文档空白

```bash
cd backend
swag init --parseInternal --parseDependency ./docs
make build
```

### Q: MFA 无法验证

确认系统时间同步：
```bash
# 容器内安装 ntpdate
docker compose exec app apk add --no-cache ntpdate
docker compose exec app ntpdate pool.ntp.org
```

### Q: AD 域登录失败

```bash
# 测试 AD 连接
docker compose exec app nc -zv ldap.example.com 389

# 检查 AD 配置
docker compose exec app cat /app/conf/app.yaml | grep -A 20 ldap
```

### Q: 如何查看详细日志？

```bash
# Docker
docker compose -f docker/docker-compose.yml logs -f --tail=100 app

# Kubernetes
kubectl -n duptwo logs -f deployment/duptwo-app --tail=100

# 直接运行
./server --log-level=debug
```
