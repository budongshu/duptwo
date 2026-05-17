# duptwo 数据登记平台 - 部署指南
# duptwo Data Registry Platform - Deployment Guide

> **版本 / Version**: v1.0.0 | **技术栈 / Stack**: Go 1.23 + Vue 3 + Gin + GORM

---

## 目录 / Table of Contents

- [1. 平台功能 / Platform Features](#1-平台功能-platform-features)
- [2. 安全功能 / Security Features](#2-安全功能-security-features)
- [3. 环境选择 / Environment Selection](#3-环境选择-environment-selection)
- [4. SQLite 开发/轻量部署 / SQLite Dev Setup](#4-sqlite-开发轻量部署-sqlite-dev-setup)
- [5. MySQL 生产部署 / MySQL Production](#5-mysql-生产部署-mysql-production)
- [6. Docker 部署 / Docker Deployment](#6-docker-部署-docker-deployment)
- [7. Kubernetes 部署 / Kubernetes Deployment](#7-kubernetes-部署-kubernetes-deployment)
- [8. Nginx 配置 / Nginx Configuration](#8-nginx-配置-nginx-configuration)

---

## 1. 平台功能 / Platform Features

| 模块 / Module | 功能说明 / Description |
|-------------|---------------------|
| **项目管理 / Project Management** | 看板/列表/网格视图、D3.js 网络图、项目周期预警 |
| **数据采集 / Data Collection** | API 推送/API 拉取、JSON/FormData/Excel 多格式、公开上传 |
| **磁盘标签 / Disk Labels** | 按磁盘标签分组统计、状态分布 |
| **数据统计 / Statistics** | 每日/每周/每月趋势图、磁盘标签分布、项目分布、成功率仪表 |
| **人员管理 / Personnel** | 多职位支持、人员负荷矩阵（卡片/散点图视图） |
| **字段配置 / Field Config** | 7 种字段类型（文本/数字/下拉/多选/日期/时间戳/多行文本） |
| **用户管理 / User Management** | RBAC 角色权限、MFA 双因素认证、批量管理、导入导出 |
| **数据审计 / Audit** | 操作日志、登录日志、数据变更追溯 |

---

## 2. 安全功能 / Security Features

### 密码策略 / Password Policy

| 配置项 / Config | 说明 / Description |
|-------------|------------------|
| 最小长度 / Min Length | 密码最小字符数（默认 8 / default 8） |
| 大写字母 / Uppercase | 是否要求 / Whether to require |
| 小写字母 / Lowercase | 是否要求 / Whether to require |
| 数字 / Digit | 是否要求 / Whether to require |
| 特殊字符 / Special Char | 是否要求（!@#$%^&*）/ Whether to require |

### 认证机制 / Authentication

| 功能 / Feature | 说明 / Description |
|--------------|------------------|
| JWT Token | 168 小时有效期 / 168h validity，支持刷新 / refresh supported |
| MFA 双因素 / MFA | TOTP 验证码 / TOTP code（Google Authenticator） |
| 登录失败锁定 / Lockout | N 次失败后锁定账户（可配置 / configurable） |
| 会话管理 / Session | 同时在线会话数限制 / concurrent session limit |

### 安全配置示例 / Security Config Example

```yaml
# conf/app.yaml
security:
  password_policy:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_digit: true
    require_special: false

  login:
    max_failed_attempts: 5
    lock_duration_minutes: 30
    max_sessions: 5
```

---

## 3. 环境选择 / Environment Selection

| 数据库 / DB | 推荐场景 / Use Case | 数据量 / Data | 运维 / Ops |
|------------|-------------------|--------------|-----------|
| **SQLite** | 开发、个人/小团队 / Dev, small team | < 10万条 | ⭐ 极简 |
| **MySQL** | 生产、多部门协作 / Production | 100万+ | ⭐⭐ 中等 |
| **PostgreSQL** | 高并发需求 / High concurrency | 无限制 | ⭐⭐ 中等 |

> 💡 切换数据库：修改 `conf/app.yaml` 中的 `database.type`，程序自动迁移 / auto-migrate.

---

## 4. SQLite 开发/轻量部署 / SQLite Dev Setup

适用于：开发调试、个人使用、小团队内部协作
Use for: dev/debug, personal use, small team internal use

### 方式一：Docker Compose（推荐 / Recommended）

```bash
# 进入 docker 目录
cd deploy/docker

# 启动 SQLite 版本
docker compose -f docker-compose.sqlite.yml up -d

# 查看日志
docker compose -f docker-compose.sqlite.yml logs -f

# 停止
docker compose -f docker-compose.sqlite.yml down
```

**访问 / Access**: http://localhost:18421
**默认账号 / Login**: `admin` / `admin123`

### 方式二：二进制直接运行 / Binary

```bash
# 1. 构建或下载二进制
cd backend

# 2. 创建目录 / Create dirs
mkdir -p data logs

# 3. 修改配置（可选，SQLite 是默认配置）
vim conf/app.sqlite.yaml

# 4. 启动
./datauptwo --config conf/app.yaml

# 或指定端口 / with port
./datauptwo --config conf/app.yaml --port 8080
```

### 备份与恢复 / Backup & Restore

```bash
# 备份（Docker）
docker cp duptwo-app:/app/data/registry.db ./backup/

# 备份（二进制 / Binary）
cp -r data data.backup.$(date +%Y%m%d)

# 恢复 / Restore
cp ./backup/registry.db data/
```

---

## 5. MySQL 生产部署 / MySQL Production

适用于：多用户并发访问、数据量大、需要定时备份
Use for: multi-user, large data, scheduled backup

### 方式一：Docker Compose（推荐 / Recommended）

```bash
cd deploy/docker

# 1. 启动 MySQL + 应用
docker compose -f docker-compose.mysql.yml up -d

# 2. 查看日志
docker compose -f docker-compose.mysql.yml logs -f

# 3. 停止
docker compose -f docker-compose.mysql.yml down
```

**访问 / Access**: http://localhost:80
**默认账号 / Login**: `admin` / `admin123`

### 方式二：二进制 + MySQL / Binary + MySQL

**Step 1: 安装 MySQL 8.0 / Install MySQL**

```bash
# Ubuntu
sudo apt update && sudo apt install mysql-server

# CentOS
sudo yum install mysql-server

# macOS
brew install mysql
```

**Step 2: 创建数据库 / Create Database**

```bash
mysql -u root -p
```

```sql
CREATE DATABASE duptwo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'duptwo'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON duptwo.* TO 'duptwo'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

**Step 3: 修改配置 / Edit Config**

```bash
vim conf/app.mysql.yaml
# 修改 db.host, db.port, db.user, db.password
```

**Step 4: 启动 / Start**

```bash
./datauptwo --config conf/app.mysql.yaml
```

### 备份与恢复 / Backup & Restore

```bash
# 备份（Docker）
docker exec duptwo-mysql mysqldump -u root -p'duptwo_root_pass' duptwo > backup_$(date +%Y%m%d).sql

# 备份（二进制 / Binary）
mysqldump -u duptwo -p duptwo > backup_$(date +%Y%m%d).sql

# 恢复 / Restore
docker exec -i duptwo-mysql mysql -u root -p'duptwo_root_pass' duptwo < backup_20260515.sql
```

---

## 6. Docker 部署 / Docker Deployment

### 目录结构 / Directory Structure

```
deploy/docker/
├── Dockerfile                    # 构建镜像 / Build image
├── docker-compose.yml            # 默认 MySQL 版 / Default MySQL
├── docker-compose.mysql.yml      # MySQL 生产版 / MySQL production
├── docker-compose.sqlite.yml     # SQLite 轻量版 / SQLite lightweight
├── docker-compose.postgres.yml   # PostgreSQL 高并发版 / PostgreSQL
├── conf/
│   ├── app.yaml                  # 默认配置 / Default
│   ├── app.mysql.yaml            # MySQL 配置 / MySQL config
│   ├── app.sqlite.yaml           # SQLite 配置 / SQLite config
│   └── app.postgres.yaml         # PostgreSQL 配置 / PostgreSQL config
└── nginx/
    └── nginx.conf                # Nginx 反向代理 / Reverse proxy
```

### 构建自定义镜像 / Build Custom Image

```bash
cd deploy/docker

# 构建镜像
docker build -t duptwo:latest .

# 推送到镜像仓库 / Push to registry
docker tag duptwo:latest your-registry.com/duptwo:latest
docker push your-registry.com/duptwo:latest
```

### 常用命令 / Common Commands

```bash
cd deploy/docker

# SQLite 开发版
docker compose -f docker-compose.sqlite.yml up -d

# MySQL 生产版
docker compose -f docker-compose.mysql.yml up -d

# PostgreSQL 高并发版
docker compose -f docker-compose.postgres.yml up -d

# 查看状态 / Check status
docker compose ps

# 查看日志 / View logs
docker compose logs -f

# 重启 / Restart
docker compose restart

# 停止并删除 / Stop and remove
docker compose down
```

### 数据卷 / Volumes

| 路径 / Path | 说明 / Description |
|-----------|------------------|
| `/app/data` | SQLite 数据库文件 / SQLite database file |
| `/app/logs` | 日志文件 / Log files |

---

## 7. Kubernetes 部署 / Kubernetes Deployment

### 前置要求 / Prerequisites

- Kubernetes 1.21+
- Helm 3.x
- 持久化存储（PVC）/ Persistent storage

### 方式一：YAML 清单 / YAML Manifests

```bash
cd deploy/kubernetes

# 修改镜像地址 / Edit image
vim 01-manifests.yaml

# 部署 MySQL 版 / Deploy MySQL
kubectl apply -f deployment-mysql.yaml

# 查看 / Check
kubectl get pods -n duptwo
kubectl get svc -n duptwo
```

### 方式二：Helm Chart（推荐 / Recommended）

```bash
cd deploy/kubernetes/helm-duptwo

# 安装 / Install
helm install duptwo . -n duptwo --create-namespace

# 自定义配置 / Custom config
helm install duptwo . -n duptwo \
  --set image.repository=your-registry.com/duptwo \
  --set image.tag=v1.0.0 \
  --set database.type=mysql \
  --set mysql.externalHost=your-mysql-host \
  --set mysql.password=xxx
```

### values.yaml 关键配置 / Key Config

```yaml
replicaCount: 2

image:
  repository: your-registry.com/duptwo
  tag: v1.0.0

database:
  type: mysql  # sqlite / mysql / postgres

mysql:
  externalHost: your-mysql-host
  database: duptwo
  username: duptwo
  password: "xxx"  # 推荐使用 Secret / Use Secret

ingress:
  enabled: true
  className: nginx
  host: duptwo.example.com
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "500m"

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 100m
    memory: 256Mi
```

### 运维操作 / Operations

```bash
# 查看 Pod / Check pods
kubectl get pods -n duptwo

# 查看日志 / View logs
kubectl logs -n duptwo -l app=duptwo -f

# 升级 / Upgrade
helm upgrade duptwo . -n duptwo --set image.tag=v1.1.0

# 回滚 / Rollback
helm rollback duptwo -n duptwo

# 健康检查 / Health check
kubectl exec -it -n duptwo deploy/duptwo -- wget -qO- http://localhost:18421/health

# 卸载 / Uninstall
helm uninstall duptwo -n duptwo
```

---

## 8. Nginx 配置 / Nginx Configuration

### 反向代理 / Reverse Proxy

```nginx
server {
    listen 80;
    server_name duptwo.example.com;

    client_max_body_size 500m;

    location / {
        proxy_pass http://127.0.0.1:18421;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 大文件上传需要增加超时 / Increase timeout for large uploads
        proxy_connect_timeout 60s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }

    location /swagger/ {
        proxy_pass http://127.0.0.1:18421/swagger/;
        proxy_set_header Host $host;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:18421/api/;
        proxy_set_header X-Real-IP $remote_addr;
        client_max_body_size 500m;
    }
}
```

### HTTPS 配置 / HTTPS Config

```nginx
server {
    listen 443 ssl http2;
    server_name duptwo.example.com;

    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;

    add_header Strict-Transport-Security "max-age=63072000" always;

    location / {
        proxy_pass http://127.0.0.1:18421;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        client_max_body_size 500m;
    }
}

# HTTP 跳转 HTTPS
server {
    listen 80;
    server_name duptwo.example.com;
    return 301 https://$server_name$request_uri;
}
```

### Let's Encrypt 免费证书 / Free SSL Certificate

```bash
# 安装 certbot
sudo apt update && sudo apt install certbot python3-certbot-nginx

# 申请证书 / Get certificate
sudo certbot --nginx -d duptwo.example.com

# 自动续期（certbot 自动添加）/ Auto-renew (certbot auto-adds)
sudo certbot renew --dry-run
```

---

## 默认账号 / Default Account

| 角色 / Role | 用户名 / Username | 密码 / Password |
|-----------|-----------------|---------------|
| 管理员 / Admin | admin | admin123 |

> ⚠️ **生产环境务必修改默认密码！/ Change default password in production!**

```bash
# 重置管理员密码（二进制 / Binary）
./datauptwo reset-admin your_new_password

# 重置管理员密码（Docker）
docker exec duptwo-app ./datauptwo reset-admin your_new_password
```

---

## 快速命令汇总 / Quick Command Reference

| 操作 / Operation | Docker | Kubernetes |
|----------------|--------|-----------|
| 启动 / Start | `docker compose up -d` | `helm install duptwo . -n duptwo` |
| 停止 / Stop | `docker compose down` | `helm uninstall duptwo -n duptwo` |
| 查看状态 / Status | `docker compose ps` | `kubectl get pods -n duptwo` |
| 查看日志 / Logs | `docker compose logs -f` | `kubectl logs -n duptwo -l app=duptwo -f` |
| 重启 / Restart | `docker compose restart` | `kubectl rollout restart deploy/duptwo -n duptwo` |
| 备份 / Backup | `docker cp duptwo-app:/app/data/registry.db ./` | mysqldump |
| 更新版本 / Update | `docker pull && docker compose up -d` | `helm upgrade duptwo . -n duptwo` |

---

## 技术支持 / Support

- **GitHub**: https://github.com/budongshu/duptwo
- **Issues**: https://github.com/budongshu/duptwo/issues
- **API 文档 / API Docs**: http://localhost:18421/swagger