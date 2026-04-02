# duptwo 数据登记平台 - 部署指南

> **版本**: v1.0.0
> **最后更新**: 2026-04-01

## 目录

- [快速开始](#快速开始)
- [系统架构](#系统架构)
- [环境要求](#环境要求)
- [Docker 部署](#docker-部署)
  - [SQLite 轻量版](#sqlite-轻量版)
  - [MySQL 生产版](#mysql-生产版)
  - [PostgreSQL 生产版](#postgresql-生产版)
- [二进制部署](#二进制部署)
- [Kubernetes 部署](#kubernetes-部署)
- [配置说明](#配置说明)
- [数据库支持](#数据库支持)
- [数据备份](#数据备份)
- [安全配置](#安全配置)
- [常见问题](#常见问题)

---

## 快速开始

### 一键部署（推荐）

```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo

# SQLite 版（开发/小规模）
./deploy/deploy.sh deploy:sqlite

# MySQL 版（生产，推荐）
./deploy/deploy.sh deploy
```

启动后访问：
| 服务 | 地址 |
|------|------|
| **Web UI** | http://localhost:80 |
| **API 直连** | http://localhost:18421 |
| **健康检查** | http://localhost:18421/health |
| **Swagger 文档** | http://localhost:18421/swagger |
| **默认账号** | `admin` / `admin123` |

---

## 系统架构

```
                           ┌─────────────────┐
                           │   用户浏览器      │
                           └────────┬────────┘
                                    │ HTTP :80 / HTTPS :443
                           ┌────────▼────────┐
                           │     Nginx       │
                           │  反向代理 + SSL  │
                           └────────┬────────┘
                                    │
                           ┌────────▼────────┐
                           │   Go API        │
                           │   :18421        │
                           │                 │
                           │  /api/*         │
                           │  /public/*      │
                           │  /health        │
                           │  /              │ ← 内置前端
                           └────────┬────────┘
                                    │
        ┌────────────────────────────┼────────────────────────────┐
        │                            │                            │
┌───────▼───────┐          ┌───────▼───────┐          ┌───────▼───────┐
│    SQLite     │          │    MySQL      │          │  PostgreSQL   │
│  ./data/*.db  │          │    :3306      │          │    :5432      │
└───────────────┘          └───────────────┘          └───────────────┘
```

**架构特点**：
- Go 服务内置静态文件服务，前端直接嵌入二进制
- 支持 SQLite、MySQL、PostgreSQL 三种数据库，修改配置即可切换
- Nginx 用于 HTTPS 终结和负载均衡（可选）
- 所有数据模型自动迁移（Auto Migrate）

---

## 环境要求

### 硬件要求

| 场景 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发/测试 | 2核 | 4GB | 20GB |
| 小规模（SQLite） | 2核 | 4GB | 50GB |
| 中等规模（MySQL/PGSQL） | 4核 | 8GB | 100GB |
| 大规模生产 | 8核+ | 16GB+ | 200GB+ |

### 软件要求

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Docker | 20.10 | 24.0+ |
| Docker Compose | 2.0 | 2.20+ |
| MySQL（可选） | 8.0 | 8.0 LTS |
| PostgreSQL（可选） | 13 | 16 LTS |
| Go（源码构建） | 1.21 | 1.23 |
| Node.js（前端构建） | 18.0 | 20 LTS |

---

## Docker 部署

### 准备工作

```bash
# 检查 Docker 是否安装
docker --version
docker compose version

# 如果未安装
# Ubuntu/Debian:
curl -fsSL https://get.docker.com | sh
# 或
apt update && apt install -y docker.io docker-compose
```

### SQLite 轻量版

适用于：**开发、个人/团队内部系统、小规模使用（< 100 用户）**

```bash
cd deploy/docker

# 启动
docker compose -f docker-compose.sqlite.yml up -d

# 查看状态
docker compose -f docker-compose.sqlite.yml ps

# 查看日志
docker compose -f docker-compose.sqlite.yml logs -f

# 停止
docker compose -f docker-compose.sqlite.yml down

# 完全清理（包括数据）
docker compose -f docker-compose.sqlite.yml down -v
```

### MySQL 生产版

适用于：**生产环境、多部门协作、多用户（100+ 用户）**

```bash
cd deploy/docker

# 启动（首次自动创建数据库）
docker compose -f docker-compose.mysql.yml up -d

# 查看日志
docker compose -f docker-compose.mysql.yml logs -f

# 停止
docker compose -f docker-compose.mysql.yml down
```

### PostgreSQL 生产版

适用于：**生产环境、高并发、强一致性需求**

```bash
cd deploy/docker

# 启动
docker compose -f docker-compose.postgres.yml up -d

# 查看日志
docker compose -f docker-compose.postgres.yml logs -f

# 停止
docker compose -f docker-compose.postgres.yml down
```

### 自定义端口

修改 `docker-compose.*.yml` 中的端口映射：

```yaml
services:
  app:
    ports:
      - "28021:18421"   # 修改为 28021:18421
  nginx:
    ports:
      - "8080:80"       # 修改为 8080:80
      - "8443:443"      # 修改为 8443:443
```

---

## 二进制部署

### 下载 Release

```bash
# 下载二进制包
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64
```

### 从源码构建

```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo

# 依赖安装
cd backend && go mod download && cd ..
cd frontend && npm ci && cd ..

# 构建前端
cd frontend && npm run build && cd ..

# 构建二进制（Linux amd64，CGO 启用支持 SQLite）
cd backend
CGO_ENABLED=1 go build -ldflags="-s -w -X datauptwo/global.VERSION=v1.0.0" -o duptwo ./cmd/server/main.go
```

### 启动

> ⚠️ **必须从 backend/ 目录运行**，二进制依赖 `./cmd/server/web` 读取前端文件。

```bash
cd backend

# 修改配置（必须修改 JWT Secret）
vim conf/app.yaml

# 启动（SQLite 模式）
./duptwo --config conf/app.yaml

# 或指定端口
./duptwo --config conf/app.yaml --port 8080
```

### 初始化管理员密码

```bash
# 重置 admin 密码
./duptwo reset-admin your_new_password

# 不带参数启动
./duptwo --config conf/app.yaml
```

### systemd 服务（Linux）

```ini
# /etc/systemd/system/duptwo.service
[Unit]
Description=duptwo Data Registry Platform
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/duptwo
ExecStart=/opt/duptwo/duptwo --config /opt/duptwo/conf/app.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# 启用服务
cp duptwo /opt/duptwo/
cp -r deploy/docker/conf /opt/duptwo/conf
systemctl daemon-reload
systemctl enable duptwo
systemctl start duptwo
```

---

## Kubernetes 部署

### 前置条件

- Kubernetes 1.20+
- kubectl 已配置
- Ingress Controller（如 nginx-ingress）

### 部署步骤

```bash
# 修改镜像地址（替换为你的镜像仓库）
# 编辑 deploy/kubernetes/manifests.yaml 中的镜像地址

# 部署
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/manifests.yaml

# 查看状态
kubectl get pods -n duptwo
kubectl get svc -n duptwo

# 查看日志
kubectl logs -l app=duptwo -n duptwo -f
```

### 挂载持久化存储

修改 `manifests.yaml` 中的 PVC 配置，使用云存储或 NFS：

```yaml
volumes:
  - name: data
    persistentVolumeClaim:
      claimName: duptwo-data-pvc
```

---

## 配置说明

### 配置文件位置

| 文件 | 用途 |
|------|------|
| `deploy/docker/conf/app.sqlite.yaml` | SQLite 配置 |
| `deploy/docker/conf/app.mysql.yaml` | MySQL 配置 |
| `deploy/docker/conf/app.postgres.yaml` | PostgreSQL 配置 |
| `backend/conf/app.yaml` | 开发配置 |

### SQLite 配置

```yaml
base:
  mode: prod
  port: 18421
  serve_web: true

database:
  type: sqlite
  path: ./data/registry.db

jwt:
  secret: YOUR_32_CHAR_SECRET_HERE  # ⚠️ 必须修改
  expire_hours: 168  # 7天

session:
  secret: YOUR_SESSION_SECRET_HERE  # ⚠️ 必须修改
```

### MySQL 配置

```yaml
database:
  type: mysql
  host: localhost       # Docker: mysql
  port: 3306
  user: duptwo
  pass: your_password   # ⚠️ 修改密码
  name: duptwo
```

### PostgreSQL 配置

```yaml
database:
  type: postgres
  host: localhost       # Docker: postgres
  port: 5432
  user: duptwo
  pass: your_password   # ⚠️ 修改密码
  name: duptwo
  sslmode: disable      # 生产改为 require
```

### 环境变量覆盖

Docker 部署时可通过环境变量覆盖配置：

```bash
docker run -d \
  -e DATABASE_TYPE=mysql \
  -e DB_HOST=mysql \
  -e DB_PORT=3306 \
  -e DB_USER=duptwo \
  -e DB_PASSWORD=your_pass \
  -e DB_NAME=duptwo \
  -p 18421:18421 \
  duptwo:latest
```

---

## 数据库支持

| 数据库 | type 值 | 适用场景 | 数据存储 | 备份方式 |
|--------|---------|---------|---------|---------|
| **SQLite** | `sqlite` | 开发、小规模 | `./data/registry.db` | `cp` 打包 |
| **MySQL** | `mysql` | 生产、多用户 | MySQL 数据库 | `mysqldump` |
| **PostgreSQL** | `postgres` | 高并发、强一致性 | PGSQL 数据库 | `pg_dump` |

### 切换数据库

修改 `database.type` 后重启服务，程序自动迁移表结构：

```yaml
# 切换到 MySQL
database:
  type: mysql

# 切换到 PostgreSQL
database:
  type: postgres
```

> ⚠️ 切换前建议备份数据

---

## 数据备份

### SQLite

```bash
# Docker 方式
docker cp duptwo-app:/app/data/registry.db ./backup/registry.db.$(date +%Y%m%d)

# 定时备份脚本
0 2 * * * docker cp duptwo-app:/app/data/registry.db /backup/registry.db.$(date +\%Y\%m\%d)
```

### MySQL

```bash
# 备份
docker exec duptwo-mysql mysqldump \
  -u root -p'duptwo_root_pass' \
  --single-transaction --quick \
  duptwo > backup_$(date +%Y%m%d).sql

# 压缩备份
docker exec duptwo-mysql mysqldump \
  -u root -p'duptwo_root_pass' duptwo | gzip > backup_$(date +%Y%m%d).sql.gz

# 恢复
docker exec -i duptwo-mysql mysql \
  -u root -p'duptwo_root_pass' duptwo < backup_20260401.sql
```

### PostgreSQL

```bash
# 备份
docker exec duptwo-postgres pg_dump \
  -U duptwo -d duptwo \
  --format=plain --compress=6 \
  > backup_$(date +%Y%m%d).sql

# 恢复
cat backup_20260401.sql | docker exec -i duptwo-postgres psql -U duptwo -d duptwo
```

---

## 安全配置

### ⚠️ 生产环境必做清单

- [ ] JWT Secret 修改为 32+ 位随机字符串
- [ ] Session Secret 修改为随机字符串
- [ ] Admin 默认密码已修改
- [ ] MySQL/PostgreSQL 密码已修改
- [ ] CORS 白名单已配置
- [ ] 生产环境设置 `mode: prod`
- [ ] 配置了 HTTPS（Let's Encrypt）
- [ ] 数据库端口（3306/5432）仅内网访问
- [ ] 配置了定时备份策略
- [ ] 备份存储与生产分离

### 生成随机密钥

```bash
# Linux/macOS
openssl rand -base64 32

# 或使用 tr
tr -dc 'A-Za-z0-9!@#$%^&*' < /dev/urandom | head -c 32
```

### 配置 HTTPS（Nginx）

```bash
# 使用 Let's Encrypt（需域名）
certbot --nginx -d your-domain.com

# 自签名证书（测试用）
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/server.key -out nginx/ssl/server.crt
```

---

## 常见问题

### Q: 启动报 `port already in use`

```bash
# 查看占用端口的进程
lsof -i :18421
lsof -i :80

# 杀掉占用进程或修改配置中的端口
```

### Q: MySQL 连接失败

```bash
# 检查容器网络
docker network ls
docker compose -f docker-compose.mysql.yml exec mysql mysql -u duptwo -pduptwo_secure_pass

# 检查连接
docker compose -f docker-compose.mysql.yml exec app curl http://mysql:3306
```

### Q: PostgreSQL 连接失败

```bash
# 检查连接
docker compose -f docker-compose.postgres.yml exec postgres psql -U duptwo -d duptwo

# 查看日志
docker compose -f docker-compose.postgres.yml logs -f postgres
```

### Q: 如何完全重置系统？

```bash
# SQLite
docker compose -f docker-compose.sqlite.yml down -v
docker compose -f docker-compose.sqlite.yml up -d

# MySQL
docker compose -f docker-compose.mysql.yml down -v
docker compose -f docker-compose.mysql.yml up -d
```

### Q: 二进制部署如何更新？

```bash
# 1. 备份数据
cp -r data data.backup.$(date +%Y%m%d)

# 2. 停止服务
pkill duptwo

# 3. 替换二进制
mv duptwo duptwo.old
cp new-duptwo duptwo

# 4. 启动
./duptwo --config conf/app.yaml
```

### Q: ARM 架构（Apple Silicon / 树莓派）？

```bash
cd backend
CGO_ENABLED=1 go build -ldflags="-s -w" -o duptwo-arm64 ./cmd/server/main.go
```

---

## 技术支持

- GitHub Issues: https://github.com/budongshu/duptwo/issues
- GitHub Releases: https://github.com/budongshu/duptwo/releases
