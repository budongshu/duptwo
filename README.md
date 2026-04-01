# duptwo 数据登记平台

> **English documentation below / 英文说明往下**

---

## 产品简介

**duptwo**（数据登记平台）是一套完整的**多渠道数据采集与管理解决方案**。通过灵活高效的 API 接口实现数据上报，支持 JSON/FormData/Excel 多格式，提供强大的后台管理能力，包括项目管理、人员管理、字段配置、报表导出等核心功能。

适用于企业级数据登记、政务数据采集、项目过程管理等多种场景。

---

## 核心功能

| 模块 | 功能说明 |
|------|---------|
| **项目管理** | D3.js 可视化网络图、看板/列表/网格视图、项目 Logo |
| **数据采集** | API 推送/API 拉取、JSON/FormData/Excel 多格式、文件上传 |
| **人员管理** | 多职位支持、人员统计、人员关联系统用户 |
| **字段配置** | 7 种字段类型（文本/数字/下拉/多选/日期/时间戳/多行文本） |
| **用户管理** | RBAC 角色权限、MFA 双因素认证、批量管理 |
| **数据审计** | 操作日志、登录日志、数据变更追溯 |
| **安全设置** | 密码策略、会话管理、CORS 配置 |
| **中英双语** | 全系统中英文切换支持 |

---

## 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + TypeScript + Vite + Pinia + Element Plus |
| **后端** | Go 1.23 + Gin + GORM |
| **数据库** | SQLite / MySQL 8.0 / PostgreSQL 16 |
| **部署** | Docker + Docker Compose + Kubernetes |

---

## 快速开始

### 🐳 Docker 部署（推荐）

**SQLite 轻量版（开发/小规模）：**
```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo
./deploy/deploy.sh deploy:sqlite
```

**MySQL 生产版：**
```bash
./deploy/deploy.sh deploy
```

**PostgreSQL 生产版：**
```bash
./deploy/deploy.sh build
cd deploy/docker
docker compose -f docker-compose.postgres.yml up -d
```

> 访问 http://localhost:80
> 默认账号：`admin` / `admin123`

### ⚙️ 二进制部署

```bash
# 下载 release 包
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64

# 修改配置
vim deploy/docker/conf/app.sqlite.yaml

# 启动
./duptwo --config deploy/docker/conf/app.sqlite.yaml
```

### 🔨 从源码构建

```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo

# 构建（前端 + 后端 + Docker 镜像）
make build-all

# Docker 部署（MySQL）
make docker-up

# 二进制运行
cd backend && ./duptwo --config conf/app.yaml
```

---

## 数据库选择指南

| 数据库 | 推荐场景 | 并发能力 | 数据量 | 运维复杂度 |
|--------|---------|---------|--------|-----------|
| **SQLite** | 开发、个人/团队内部 | 低 | < 10万条 | ⭐ 极简 |
| **MySQL** | 生产环境、多部门协作 | 高 | 100万+ | ⭐⭐ 中等 |
| **PostgreSQL** | 高并发、强一致性需求 | 极高 | 无限制 | ⭐⭐ 中等 |

> 💡 从 SQLite 切换到 MySQL/PostgreSQL：修改 `database.type`，程序自动迁移表结构。

---

## 目录结构

```
duptwo/
├── backend/                  # Go 后端（package: datauptwo）
│   ├── cmd/server/          # 入口，嵌入前端静态文件
│   ├── app/api/            # API 处理器
│   ├── app/service/         # 业务逻辑
│   ├── app/repo/            # 数据访问层
│   ├── app/model/           # GORM 模型
│   ├── app/dto/             # 请求/响应 DTO
│   ├── conf/                # 配置文件
│   ├── Makefile             # 构建命令
│   └── go.mod
├── frontend/                # Vue 3 前端
│   ├── src/
│   ├── vite.config.ts       # 构建到 backend/cmd/server/web
│   └── package.json
├── deploy/                  # 部署配置
│   ├── docker/              # Docker 部署
│   │   ├── Dockerfile
│   │   ├── docker-compose.sqlite.yml
│   │   ├── docker-compose.mysql.yml
│   │   ├── docker-compose.postgres.yml
│   │   └── nginx/
│   ├── kubernetes/           # K8s 部署
│   └── deploy.sh             # 一键部署脚本
├── .github/workflows/        # CI/CD
└── SPEC.md                  # 详细规格说明
```

---

## API 接口

| 类别 | 路径 | 说明 |
|------|------|------|
| 认证 | `POST /api/auth/login` | 用户登录 |
| 认证 | `POST /api/auth/register` | 用户注册 |
| 认证 | `POST /api/auth/mfa/verify` | MFA 验证 |
| 上传 | `POST /public/upload-records` | 公开上传（无需认证）|
| 上传 | `GET/POST/PUT/DELETE /api/upload-records` | 上传记录 CRUD |
| 项目 | `GET/POST/PUT/DELETE /api/projects` | 项目管理 |
| 人员 | `GET/POST/PUT/DELETE /api/personnels` | 人员管理 |
| 用户 | `GET/POST/PUT/DELETE /api/users` | 用户管理 |
| 角色 | `GET/POST/PUT/DELETE /api/roles` | 角色管理 |
| 用户组 | `GET/POST/PUT/DELETE /api/user-groups` | 用户组管理 |
| 字段 | `GET/POST/PUT/DELETE /api/field-configs` | 字段配置 |
| 审计 | `GET /api/audit/*` | 日志审计 |
| 系统 | `GET/PUT /api/admin/*` | 系统配置 |
| 健康 | `GET /health` | 健康检查 |

> Swagger 文档：http://localhost:18421/swagger

---

## 安全配置

### ⚠️ 生产环境必做

1. **修改 JWT Secret** — 生成 32+ 位随机字符串
   ```bash
   # 启动后重置管理员密码
   ./duptwo reset-admin your_new_password
   ```

2. **配置 HTTPS** — 使用 Let's Encrypt 或商业证书

3. **配置 CORS** — 仅允许受信任的域名
   ```yaml
   cors:
     allow_origins:
       - "https://your-domain.com"
   ```

4. **数据库密码** — MySQL/PostgreSQL 生产环境设置强密码

---

## 数据备份

**SQLite：**
```bash
# 方式1：Docker 卷
docker cp duptwo-app:/app/data/registry.db ./backup/

# 方式2：直接复制
cp -r data data.backup.$(date +%Y%m%d)
```

**MySQL：**
```bash
docker exec duptwo-mysql mysqldump -u root -p'duptwo_root_pass' duptwo > backup_$(date +%Y%m%d).sql

# 恢复
docker exec -i duptwo-mysql mysql -u root -p'duptwo_root_pass' duptwo < backup_20260401.sql
```

**PostgreSQL：**
```bash
docker exec duptwo-postgres pg_dump -U duptwo duptwo > backup_$(date +%Y%m%d).sql

# 恢复
docker exec -i duptwo-postgres psql -U duptwo duptwo < backup_20260401.sql
```

---

## 环境要求

| 场景 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发/测试 | 2核 | 4GB | 20GB |
| 小规模（SQLite） | 2核 | 4GB | 50GB |
| 中等规模（MySQL/PGSQL） | 4核 | 8GB | 100GB |
| 大规模生产 | 8核+ | 16GB+ | 200GB+ |

---

## 常见问题

**Q: 启动报 `port already in use`**
```bash
lsof -i :18421
pkill duptwo
```

**Q: 二进制部署如何更新？**
```bash
cp -r data data.backup
pkill duptwo
mv duptwo duptwo.old
# 替换新二进制后重启
./duptwo --config conf/app.yaml
```

**Q: 支持 ARM 架构吗？**
```bash
# 在 ARM 服务器上直接构建
cd backend && CGO_ENABLED=1 go build -o duptwo ./cmd/server/main.go
```

---

## 贡献

欢迎提交 Issue 和 Pull Request！

---

## License

MIT

---

---

# duptwo Data Registry Platform

A full-stack data collection and management platform with flexible API support, multi-database backend (SQLite/MySQL/PostgreSQL), and comprehensive admin capabilities.

### Quick Start

```bash
# SQLite (dev/lightweight)
git clone https://github.com/budongshu/duptwo.git && cd duptwo
./deploy/deploy.sh deploy:sqlite

# MySQL (production)
./deploy/deploy.sh deploy

# PostgreSQL (production, high-concurrency)
cd deploy/docker && docker compose -f docker-compose.postgres.yml up -d
```

Visit http://localhost:80 — Login: `admin` / `admin123`

### Features

- **Multi-format upload**: JSON / FormData / Excel
- **Project management**: D3.js network graph, kanban/list/grid views
- **Personnel tracking**: Multi-position support, system user linking
- **RBAC**: Roles, user groups, granular permissions
- **MFA**: TOTP-based two-factor authentication
- **Bilingual**: Full Chinese/English i18n support
- **Audit logs**: Operation and login trail

### Architecture

```
User → Nginx (80/443) → Go API (:18421) → SQLite / MySQL / PostgreSQL
```

### Binary Build

```bash
make build-all    # Frontend + Backend binary
make build-docker # Docker image
```

### Docs

- [Deployment Guide](DEPLOY.md)
- [API Documentation](http://localhost:18421/swagger)
