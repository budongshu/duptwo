# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**duptwo** (数据登记平台) is a full-stack data management platform for multi-channel data collection, backend management, statistics, and report export. Supports JSON/FormData/Excel formats and provides flexible API interfaces.

Project internal name: `duptwo` (backend Go package `datauptwo`)
Frontend product name: **数据登记平台**
English name: **duptwo**

## Tech Stack

- **Backend**: Go 1.23 + Gin + GORM (layered architecture: API → Service → Repository)
- **Frontend**: Vue 3 + TypeScript + Vite + Pinia + Element Plus
- **Database**: SQLite (default), MySQL, PostgreSQL supported via GORM
- **Deployment**: Docker + Kubernetes, monolith serving frontend from backend
- **Auth**: JWT + MFA (TOTP) support

## Build Commands

### Full project (from backend/ directory)
```bash
make build-all      # Build frontend + backend binary
make build-docker   # Build Docker image (debian-slim based)
make docker-up     # Start with MySQL (docker compose)
make docker-up-sqlite  # Start SQLite version
make dev           # Dev mode
```

### Frontend-only (from frontend/)
```bash
npm run dev        # Dev server on port 4004
npm run build      # Production build → ../backend/cmd/server/web
```

### Backend-only (from backend/)
```bash
go run cmd/server/main.go  # Dev server on port 18421
make build                # Build binary
make build-docker        # Build Docker image
```

## Architecture

### Project Structure
```
DataRegistry/
├── backend/          # Go backend (package: datauptwo)
│   ├── cmd/server/  # Entry point, serves static web files
│   ├── app/api/     # API handlers (v1)
│   ├── app/service/ # Business logic
│   ├── app/repo/    # Data access layer
│   ├── app/model/   # GORM database models
│   ├── app/dto/     # Request/response DTOs
│   ├── conf/        # Configuration (YAML)
│   └── Makefile     # Build commands
├── frontend/        # Vue 3 + TypeScript frontend
├── deploy/         # Deployment files
│   ├── DEPLOY.md   # Full deployment guide
│   ├── deploy.sh   # One-click deploy script
│   ├── docker/     # Docker + compose files
│   ├── kubernetes/ # K8s manifests
│   └── clients/     # Client API examples
└── is-tmp/         # Old/ref files (ignored)
```

### Database Models
- **UploadRecord**: SerialNo (8-char UUID), DataType, ProjectName, FilePath, FileSize, Uploader, Status, Remark, Data (JSON)
- **FieldConfig**: Name, Code, Type, Options (JSON), Enabled
- **User/Role/UserGroup**: RBAC support with MFA
- **OperationLog/LoginLog**: Audit trail

### API Endpoints Summary
| Category | Path | Auth | Description |
|----------|------|------|-------------|
| Public | `POST /api/auth/login` | No | 用户登录 |
| Public | `POST /api/auth/register` | No | 用户注册 |
| Public | `POST /api/auth/mfa/verify` | No | MFA 验证 |
| Auth | `GET /api/auth/current` | JWT | 获取当前用户 |
| Auth | `PUT /api/auth/profile` | JWT | 更新个人资料 |
| Auth | `POST /api/auth/change-password` | JWT | 修改密码 |
| Auth | `GET /api/auth/mfa/status` | JWT | MFA 状态 |
| Auth | `POST /api/auth/mfa/enable` | JWT+Admin | 启用 MFA |
| Auth | `POST /api/auth/mfa/disable` | JWT+Admin | 禁用 MFA |
| Upload | `/api/upload-records` | JWT | 上传记录 CRUD |
| Project | `/api/projects` | JWT | 项目管理 |
| Personnel | `/api/personnels` | JWT | 人员管理 |
| User | `/api/users` | JWT | 用户管理 |
| Role | `/api/roles` | JWT | 角色管理 |
| UserGroup | `/api/user-groups` | JWT | 用户组管理 |
| FieldConfig | `/api/field-configs` | JWT | 字段配置 |
| Audit | `/api/audit/*` | JWT | 日志审计 |
| System | `/api/admin/*` | JWT | 系统配置 |
| Public | `POST /public/upload-records` | No | 公开上传 |
| Health | `GET /health` | No | 健康检查 |
| Swagger | `GET /swagger/` | No | API 文档 |

### MFA Authentication Flow
1. 用户登录 → `POST /api/auth/login`
2. 若用户已启用 MFA，服务端返回 `mfaRequired: true` + 临时 token
3. 前端跳转 MFA 验证页 → `POST /api/auth/mfa/verify` 输入验证码
4. 验证通过 → 返回正式 JWT token → 跳转首页

### Key Config
- Backend port: 18421 (dev), 8080 (prod/Docker)
- Frontend proxies to backend at localhost:18421 in dev mode
- Frontend build output: `../backend/cmd/server/web` (served as static files)
- Database: SQLite at `./data/registry.db` (dev), MySQL (prod)
- Config: `backend/conf/app.yaml` (dev), `deploy/docker/conf/app.prod.yaml` (prod)
- Backend package name: `datauptwo` (important for Go imports)

## Deployment

See `deploy/DEPLOY.md` for full deployment instructions.

### Quick Deploy
```bash
cd deploy && ./deploy.sh deploy        # Docker + MySQL
cd deploy && ./deploy.sh deploy:sqlite # Docker + SQLite
```

### Production Checklist
- [ ] Change JWT secret in config
- [ ] Change admin password
- [ ] Configure CORS allowed origins
- [ ] Enable HTTPS via Nginx
- [ ] Set `log.level: info`
- [ ] Configure database backup
