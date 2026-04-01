# Changelog

All notable changes for duptwo v1.0.0 will be documented in this file.

## [v1.0.0] - 2026-04-01

### Added

- **三数据库支持**：SQLite / MySQL / PostgreSQL，修改配置即可切换
- **Docker Compose 部署**：
  - `docker-compose.sqlite.yml` — SQLite 轻量版（开发/小规模）
  - `docker-compose.mysql.yml` — MySQL 生产版
  - `docker-compose.postgres.yml` — PostgreSQL 生产版
- **二进制部署**：Linux amd64 可执行文件，支持本地部署
- **一键部署脚本**：`./deploy/deploy.sh` 支持所有部署方式
- **GitHub Actions CI/CD**：自动构建 Docker 镜像和二进制 Release
- **RBAC 角色权限管理**：5 个默认角色（管理员、数据操作员、项目经理、审计员、只读）
- **MFA 双因素认证**：TOTP 身份验证器支持
- **JWT 会话管理**：可配置过期时间
- **操作日志审计**：操作日志 + 登录日志
- **数据备份**：SQLite 文件复制、MySQL mysqldump、PostgreSQL pg_dump
- **systemd 服务配置**：Linux 系统服务部署
- **Kubernetes 部署**：Namespace + Deployment + Service + Ingress + PVC

### Modules

| 模块 | 功能 |
|------|------|
| 项目管理 | D3.js 可视化网络图、看板/列表/网格视图、项目 Logo |
| 数据采集 | API 推送/API 拉取、JSON/FormData/Excel |
| 人员管理 | 多职位支持（16 种职位）、人员统计、系统用户关联 |
| 字段配置 | 7 种字段类型、必填、默认值、占位符 |
| 用户管理 | 批量创建/编辑/删除、重置密码、MFA 管理 |
| 角色管理 | 角色 CRUD、权限分配 |
| 用户组管理 | 用户组 CRUD |
| 安全设置 | 密码策略、会话超时、CORS 配置 |
| 中英双语 | 全系统 i18n 支持 |

### Default Credentials

> ⚠️ **生产环境必须修改默认密码**

- 用户名: `admin`
- 密码: `admin123`

### Database

| 字段 | SQLite | MySQL | PostgreSQL |
|------|--------|-------|------------|
| `type` | `sqlite` | `mysql` | `postgres` |
| 路径/连接 | `./data/registry.db` | `host:port` | `host:port` |
| 数据存储 | 单文件 | 数据库 | 数据库 |
| 并发能力 | 低 | 高 | 极高 |
| 适用规模 | < 100 用户 | 100 万+ | 无限制 |
| 运维难度 | ⭐ | ⭐⭐ | ⭐⭐ |

### Security Checklist

- [x] JWT Secret 可配置
- [x] Session Secret 可配置
- [x] CORS 白名单可配置
- [x] 密码策略可配置（最小长度、大小写、数字、特殊字符）
- [x] 会话超时可配置
- [x] HTTPS 支持（Nginx）
- [x] 非 root 用户运行（Docker/二进制）

### Known Issues

None at v1.0.0 release.
