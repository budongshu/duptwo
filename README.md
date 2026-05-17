# duptwo 数据登记平台 / Data Registry Platform

> **English documentation below**

---

## 产品简介 / Product Overview

**duptwo** 是一套完整的**多渠道数据采集与管理解决方案**。通过灵活的 API 接口实现数据上报，支持 JSON/FormData/Excel 多格式，提供项目管理、数据采集、人员管理、字段配置、报表统计、数据审计等核心功能。适用于企业级数据登记、政务数据采集、项目过程管理等多种场景。

**duptwo** is a full-stack data collection and management platform with flexible API support, multi-format upload (JSON/FormData/Excel), project management, D3.js network visualization, RBAC, MFA, and bilingual Chinese/English UI.

---

## 平台功能 / Platform Features

| 模块 / Module | 功能说明 / Description |
|-------------|---------------------|
| **项目管理 / Project Management** | 看板/列表/网格视图、D3.js 网络图、项目周期预警 |
| **数据采集 / Data Collection** | API 推送/API 拉取、JSON/FormData/Excel 多格式、公开上传 |
| **磁盘标签 / Disk Labels** | 按磁盘标签分组统计、状态分布（完成/失败/混合/待处理） |
| **数据统计 / Statistics** | 每日/每周/每月趋势图、磁盘标签分布、项目分布 |
| **人员管理 / Personnel** | 多职位支持、人员负荷矩阵、人员统计 |
| **字段配置 / Field Config** | 7 种字段类型（文本/数字/下拉/多选/日期/时间戳/多行文本） |
| **用户管理 / User Management** | RBAC 角色权限、MFA 双因素认证、批量管理、导入导出 |
| **数据审计 / Audit** | 操作日志、登录日志、数据变更追溯 |

### 技术特性 / Technical Features

- **多数据库 / Multi-database**: SQLite / MySQL 8.0 / PostgreSQL 16 一键切换
- **中英双语 / Bilingual**: 全系统中文/English 切换
- **API 文档 / API Docs**: 内置 Swagger（`/swagger`）
- **健康检查 / Health Check**: `GET /health`

---

## 安全功能 / Security Features

| 功能 / Feature | 说明 / Description |
|--------------|-------------------|
| **密码策略 / Password Policy** | 长度、大小写、数字、特殊字符可配置 |
| **JWT Token** | 168 小时有效期，支持刷新 |
| **MFA 双因素 / MFA** | TOTP 验证码（Google Authenticator） |
| **登录锁定 / Login Lockout** | N 次失败后锁定账户（可配置） |
| **会话管理 / Session** | 同时在线会话数限制、强制登出 |

---

## 环境选择 / Environment Selection

| 数据库 / DB | 推荐场景 / Use Case | 数据量 / Data Volume | 运维 / Ops |
|------------|-------------------|--------------------|-----------|
| **SQLite** | 开发、个人/小团队 / Dev, small team | < 10万条 / < 100K rows | ⭐ 极简 / Minimal |
| **MySQL** | 生产环境、多部门协作 / Production | 100万+ / 1M+ rows | ⭐⭐ 中等 / Medium |
| **PostgreSQL** | 高并发、强一致性 / High concurrency | 无限制 / Unlimited | ⭐⭐ 中等 / Medium |

> 💡 切换数据库只需修改 `conf/app.yaml` 中的 `database.type`，程序自动迁移表结构。

---

## 快速开始 / Quick Start

### 开发/轻量部署 SQLite（推荐个人/小团队）

```bash
# 方式一：Docker Compose（推荐）
git clone https://github.com/budongshu/duptwo.git
cd duptwo/deploy/docker
docker compose -f docker-compose.sqlite.yml up -d

# 方式二：二进制直接运行
cd backend
mkdir -p data logs
./datauptwo --config conf/app.sqlite.yaml
```

访问 http://localhost:18421 | 默认账号：`admin` / `admin123`

### 生产环境部署 MySQL

```bash
# Docker Compose
cd deploy/docker
docker compose -f docker-compose.mysql.yml up -d

# 或 Helm（Kubernetes）
cd deploy/kubernetes/helm-duptwo
helm install duptwo . -n duptwo --create-namespace \
  --set mysql.externalHost=your-mysql-host \
  --set mysql.password=your_password
```

访问 http://localhost:80 | 默认账号：`admin` / `admin123`

---

## 详细部署文档 / Detailed Deployment Guides

### [1. SQLite 开发/轻量部署 →](DEPLOY.md#开发轻量部署sqlite)

适用于：开发调试、个人使用、小团队内部协作

```bash
# Docker
docker compose -f docker-compose.sqlite.yml up -d

# 二进制
cd backend && ./datauptwo --config conf/app.sqlite.yaml
```

### [2. MySQL 生产部署 →](DEPLOY.md#生产环境部署mysql)

适用于：多用户并发访问、数据量大、需要定时备份

```bash
# Docker
docker compose -f docker-compose.mysql.yml up -d

# 二进制 + MySQL
mysql -u root -p -e "CREATE DATABASE duptwo DEFAULT CHARSET utf8mb4;"
./datauptwo --config conf/app.mysql.yaml
```

### [3. Docker 完整部署 →](DEPLOY.md#docker-部署)

```bash
cd deploy/docker

# 构建自定义镜像
docker build -t duptwo:latest .

# 推送
docker tag duptwo:latest your-registry.com/duptwo:latest
docker push your-registry.com/duptwo:latest

# SQLite
docker compose -f docker-compose.sqlite.yml up -d

# MySQL
docker compose -f docker-compose.mysql.yml up -d

# PostgreSQL
docker compose -f docker-compose.postgres.yml up -d

# 查看日志
docker compose logs -f

# 停止
docker compose down
```

### [4. Kubernetes 部署 →](DEPLOY.md#kubernetes-部署)

```bash
# Helm Chart 安装
cd deploy/kubernetes/helm-duptwo
helm install duptwo . -n duptwo --create-namespace \
  --set image.repository=your-registry.com/duptwo \
  --set database.type=mysql

# 查看状态
kubectl get pods -n duptwo

# 升级
helm upgrade duptwo . -n duptwo --set image.tag=v1.1.0

# 卸载
helm uninstall duptwo -n duptwo
```

### [5. Nginx 配置 →](DEPLOY.md#nginx-配置)

**反向代理（HTTP）：**
```nginx
server {
    listen 80;
    server_name duptwo.example.com;

    client_max_body_size 100m;

    location / {
        proxy_pass http://127.0.0.1:18421;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_connect_timeout 60s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }
}
```

**HTTPS 配置 + Let's Encrypt 证书：**
```bash
# 申请证书
sudo certbot --nginx -d duptwo.example.com

# 自动续期
sudo certbot renew --dry-run
```

**关键配置说明：**
- `client_max_body_size 500m;` — 允许上传大文件
- `proxy_read_timeout 300s;` — 大文件上传需要延长超时
- `proxy_buffering off;` — 流式响应（如 Swagger）建议关闭缓冲

---

## 技术栈 / Tech Stack

| 层级 / Layer | 技术 / Technology |
|------------|------------------|
| **前端 / Frontend** | Vue 3 + TypeScript + Vite + Pinia + Element Plus |
| **后端 / Backend** | Go 1.23 + Gin + GORM |
| **数据库 / Database** | SQLite / MySQL 8.0 / PostgreSQL 16 |
| **部署 / Deployment** | Docker + Docker Compose + Kubernetes + Helm |

---

## 环境要求 / Requirements

| 场景 / Scenario | CPU | 内存 / Memory | 磁盘 / Disk |
|----------------|-----|--------------|-----------|
| 开发/测试 / Dev | 2核 | 4GB | 20GB |
| 小规模 SQLite / Small | 2核 | 4GB | 50GB |
| 中等规模 MySQL / Medium | 4核 | 8GB | 100GB |
| 大规模生产 / Production | 8核+ | 16GB+ | 200GB+ |

---

## 常见问题 / FAQ

**Q: 端口被占用 / Port in use**
```bash
lsof -i :18421
pkill duptwo
```

**Q: 更新版本 / Update version**
```bash
cp -r data data.backup
pkill duptwo
./datauptwo --config conf/app.yaml
```

**Q: 修改默认密码 / Change default password**
```bash
./datauptwo reset-admin your_new_password
```

---

## 技术支持 / Support

- **GitHub**: https://github.com/budongshu/duptwo
- **Issues**: https://github.com/budongshu/duptwo/issues
- **API Docs**: http://localhost:18421/swagger

---

## 构建 / Build

```bash
# 从源码构建（前端 + 后端 + Docker 镜像）
make build-all

# Docker 镜像
make build-docker

# 二进制运行
make run-dev
```

---

**License**: MIT