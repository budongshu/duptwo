<!-- ============================================================
  语言切换 / Language Switch
  ============================================================ -->
<p align="center">
  <span>🌐</span>
  <a href="#中文"><button style="padding:6px 16px;border-radius:6px;border:1.5px solid #409eff;background:#409eff;color:#fff;font-size:13px;cursor:pointer">🇨🇳 中文</button></a>
  <a href="#english"><button style="padding:6px 16px;border-radius:6px;border:1.5px solid #67c23a;background:#fff;color:#67c23a;font-size:13px;cursor:pointer">🇺🇸 English</button></a>
</p>

---

<!-- ============================================================
  中文文档
  ============================================================ -->
<h1 id="中文">🇨🇳 duptwo 数据登记平台</h1>

<p align="center">
  <img src="frontend/src/assets/favicon.svg" alt="duptwo" width="64" height="64">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go">
  <img src="https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat-square&logo=vue.js">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=flat-square">
</p>

<p align="center">
  <strong>duptwo</strong> 是一套完整的<strong>多渠道数据采集与管理解决方案</strong>。
  支持 API 推送、Excel 导入、公开上传，提供项目管理、数据统计、人员管理、RBAC、MFA、数据审计等企业级功能。
</p>

---

## ✨ 平台功能

| 功能模块 | 说明 |
|---------|------|
| 📊 **[数据概览仪表盘](#-快速开始)** | KPI 看板、每日趋势图、磁盘状态分布、项目分布，全局筛选联动 |
| 📁 **项目管理** | 看板视图 / 列表视图 / 网格视图、D3.js 网络关系图、项目周期预警（距到期≤30天自动提示） |
| 📤 **数据采集** | API 推送（`POST /api/upload-records`）、Excel 批量导入、公开上传（`POST /public/upload-records`，无需登录） |
| 🔄 **数据同步** | Agent / Center 双向同步模式，支持多节点数据汇聚与分发 |
| 🏷️ **磁盘标签** | 按标签分组统计、状态分布（完成/失败/混合/待处理），支持快速筛选 |
| 👥 **人员管理** | 多职位支持（算法工程师/前端工程师/DBA 等）、人员负荷矩阵（卡片象限/散点图视图） |
| 🔍 **字段配置** | 7 种字段类型：文本、数字、下拉、多选、日期、时间戳、多行文本 |
| 👤 **用户管理** | RBAC 角色权限、MFA 双因素（Google Authenticator）、批量导入/导出、强制登出 |
| 🌐 **AD/LDAP 登录** | 企业 Active Directory / LDAP 单点登录，自动同步用户与部门 |
| 📋 **操作审计** | 操作日志、登录日志、数据变更追溯 |
| 🗂️ **系统设置** | 安全策略（密码强度、登录锁定）、CORS 配置、会话管理 |
| 🎮 **彩蛋** | 在任意页面按下 ↑↑↓↓←→←→BA 触发开发者终端（Konami Code） |

### 技术特性

| 特性 | 说明 |
|-----|------|
| 🌐 中英双语 | 全系统右上角一键切换 |
| 📡 RESTful API | 内置 Swagger 文档（`/swagger`） |
| ❤️ 健康检查 | `GET /health` 返回服务状态 |
| 🗄️ 多数据库 | SQLite / MySQL / PostgreSQL 一键切换，程序自动迁移 |
| 📦 多格式 | JSON / FormData / Excel 上传 |
| 🔐 安全认证 | JWT（168h有效期）+ MFA + 登录锁定 |

---

## 🔐 安全功能

| 功能 | 说明 |
|-----|------|
| **密码策略** | 最小长度、大写字母、小写字母、数字、特殊字符，可按需启用 |
| **JWT Token** | 168 小时有效期，支持刷新令牌 |
| **MFA 双因素** | TOTP 验证码（Google Authenticator / 1Password） |
| **登录失败锁定** | N 次失败后锁定账户（可配置时长） |
| **会话管理** | 同时在线会话数限制、管理员强制登出 |
| **CORS 控制** | 白名单域名单独配置 |

---

## 🚀 快速开始

### SQLite 轻量部署（开发/个人/小团队）

```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo/deploy/docker
docker compose -f docker-compose.sqlite.yml up -d
```

访问 **http://localhost:18421** | 账号：`admin` / `admin123`

### MySQL 生产部署

```bash
cd duptwo/deploy/docker
docker compose -f docker-compose.mysql.yml up -d
```

访问 **http://localhost:80** | 账号：`admin` / `admin123`

---

## 📖 详细部署文档

> 点击标题即可跳转到完整部署指南 👇

| 部署方式 | 说明 | 链接 |
|---------|------|------|
| 🦪 **SQLite 轻量部署** | 开发调试、个人使用 | [→ DEPLOY.md - SQLite](DEPLOY.md#4-sqlite-开发轻量部署-sqlite-dev-setup) |
| 🗄️ **MySQL 生产部署** | 多用户并发、数据量大 | [→ DEPLOY.md - MySQL](DEPLOY.md#5-mysql-生产部署-mysql-production) |
| 🐳 **Docker 完整部署** | 自定义镜像构建、Docker Compose | [→ DEPLOY.md - Docker](DEPLOY.md#6-docker-部署-docker-deployment) |
| ☸️ **Kubernetes 部署** | Helm Chart、YAML 清单 | [→ DEPLOY.md - K8s](DEPLOY.md#7-kubernetes-部署-kubernetes-deployment) |
| 🌐 **Nginx 配置** | 反向代理、HTTPS、Let's Encrypt | [→ DEPLOY.md - Nginx](DEPLOY.md#8-nginx-配置-nginx-configuration) |

---

## 🛠️ 开发

```bash
# 克隆
git clone https://github.com/budongshu/duptwo.git

# 从源码构建
make build-all

# Docker 镜像构建
make build-docker

# 开发模式（前端 Vite :4000 + 后端 Go :18421）
make run-dev
```

---

## 📁 项目结构

```
DataRegistry/
├── backend/              # Go 后端 (package: datauptwo)
│   ├── app/api/v1/      # API 处理器
│   ├── app/service/      # 业务逻辑
│   ├── app/repo/        # 数据访问层
│   ├── app/model/       # GORM 模型
│   └── app/dto/         # 请求/响应 DTO
├── frontend/             # Vue 3 前端
│   └── src/
│       ├── api/         # API 调用
│       ├── views/       # 页面视图
│       ├── composables/ # 组合式函数
│       └── lang/        # 中英双语
└── deploy/              # 部署文件
    ├── docker/         # Docker + Compose
    ├── kubernetes/     # K8s YAML + Helm
    └── sync-test/      # 同步测试用例
```

---

## 🗺️ 截图预览

> 平台界面截图（待补充）

<details>
<summary>📸 点击展开截图</summary>

| 功能 | 预览 |
|-----|------|
| 数据概览仪表盘 | *(截图待补充)* |
| 项目管理 - 看板视图 | *(截图待补充)* |
| 人员负荷矩阵 | *(截图待补充)* |
| Swagger API 文档 | *(截图待补充)* |

> 💡 如何添加截图：将截图文件放入项目根目录或 `docs/` 目录，然后在 PR 中更新本节。

</details>

---

## 📧 技术支持

- 🐛 **问题反馈**: [GitHub Issues](https://github.com/budongshu/duptwo/issues)
- 📘 **文档**: [DEPLOY.md](DEPLOY.md) | 内置 Swagger (`/swagger`)
- 🔧 **API 基础地址**: `http://localhost:18421/api`

---

## 📄 License

MIT © [budongshu](https://github.com/budongshu)

---

---

<!-- ============================================================
  English Documentation
  ============================================================ -->
<h1 id="english">🇺🇸 duptwo Data Registry Platform</h1>

<p align="center">
  <img src="frontend/src/assets/favicon.svg" alt="duptwo" width="64" height="64">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go">
  <img src="https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat-square&logo=vue.js">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=flat-square">
</p>

<p align="center">
  <strong>duptwo</strong> is a full-stack data collection and management platform with flexible API support, multi-format upload (JSON/FormData/Excel), project management, D3.js network graph, RBAC, MFA, LDAP, and bilingual Chinese/English UI.
</p>

---

## ✨ Platform Features

| Feature | Description |
|---------|-------------|
| 📊 **[Dashboard](#-quick-start)** | KPI cards, daily trends, disk status, project distribution, global filter |
| 📁 **Project Management** | Kanban / List / Grid / D3.js network views, project cycle deadline warning (≤30 days) |
| 📤 **Data Collection** | REST API push (`POST /api/upload-records`), Excel import, public upload (no auth) |
| 🔄 **Data Sync** | Agent / Center dual-mode sync, multi-node aggregation and distribution |
| 🏷️ **Disk Labels** | Group by labels, status distribution (completed/failed/mixed/pending) |
| 👥 **Personnel** | Multi-role support, workload matrix (card quadrant / scatter chart views) |
| 🔍 **Field Config** | 7 field types: text, number, select, multi-select, date, timestamp, textarea |
| 👤 **User Management** | RBAC roles, MFA (Google Authenticator), bulk import/export, force logout |
| 🌐 **AD/LDAP Login** | Enterprise Active Directory / LDAP SSO, auto-sync users and departments |
| 📋 **Audit Logs** | Operation logs, login logs, data change traceability |
| 🗂️ **System Settings** | Password policy, login lockout, CORS, session management |
| 🎮 **Easter Egg** | Press ↑↑↓↓←→←→BA anywhere to trigger the developer terminal (Konami Code) |

### Technical Features

| Feature | Description |
|---------|-------------|
| 🌐 Bilingual | Full Chinese/English toggle in top-right corner |
| 📡 RESTful API | Built-in Swagger docs (`/swagger`) |
| ❤️ Health Check | `GET /health` |
| 🗄️ Multi-DB | SQLite / MySQL / PostgreSQL one-click switch, auto-migrate |
| 📦 Multi-Format | JSON / FormData / Excel upload |
| 🔐 Auth | JWT (168h) + MFA + Login lockout |

---

## 🔐 Security

| Feature | Description |
|---------|-------------|
| **Password Policy** | Min length, uppercase/lowercase/digit/special char (configurable) |
| **JWT Token** | 168h validity, refresh token support |
| **MFA** | TOTP (Google Authenticator / 1Password) |
| **Login Lockout** | N failed attempts → account locked (configurable duration) |
| **Session Management** | Max concurrent sessions, admin force logout |
| **CORS** | Whitelist-based origin control |

---

## 🚀 Quick Start

### SQLite (Dev / Personal / Small Team)

```bash
git clone https://github.com/budongshu/duptwo.git
cd duptwo/deploy/docker
docker compose -f docker-compose.sqlite.yml up -d
```

Access **http://localhost:18421** | Login: `admin` / `admin123`

### MySQL (Production)

```bash
cd duptwo/deploy/docker
docker compose -f docker-compose.mysql.yml up -d
```

Access **http://localhost:80** | Login: `admin` / `admin123`

---

## 📖 Detailed Deployment Guides

> Click a title to jump to the full deployment guide 👇

| Deployment | Description | Link |
|-----------|-------------|------|
| 🦪 **SQLite** | Dev, personal use | [→ DEPLOY.md - SQLite](DEPLOY.md#4-sqlite-开发轻量部署-sqlite-dev-setup) |
| 🗄️ **MySQL** | Multi-user, production | [→ DEPLOY.md - MySQL](DEPLOY.md#5-mysql-生产部署-mysql-production) |
| 🐳 **Docker** | Custom image, Compose | [→ DEPLOY.md - Docker](DEPLOY.md#6-docker-部署-docker-deployment) |
| ☸️ **Kubernetes** | Helm Chart, YAML | [→ DEPLOY.md - K8s](DEPLOY.md#7-kubernetes-部署-kubernetes-deployment) |
| 🌐 **Nginx** | Reverse proxy, HTTPS, Let's Encrypt | [→ DEPLOY.md - Nginx](DEPLOY.md#8-nginx-配置-nginx-configuration) |

---

## 🛠️ Development

```bash
# Clone
git clone https://github.com/budongshu/duptwo.git

# Build from source
make build-all

# Docker image
make build-docker

# Dev mode (frontend Vite :4000 + backend Go :18421)
make run-dev
```

---

## 📁 Project Structure

```
DataRegistry/
├── backend/              # Go backend (package: datauptwo)
│   ├── app/api/v1/      # API handlers
│   ├── app/service/      # Business logic
│   ├── app/repo/        # Data access layer
│   ├── app/model/       # GORM models
│   └── app/dto/         # Request/Response DTOs
├── frontend/             # Vue 3 frontend
│   └── src/
│       ├── api/         # API calls
│       ├── views/       # Page views
│       ├── composables/ # Composables
│       └── lang/        # i18n (Chinese/English)
└── deploy/              # Deployment files
    ├── docker/         # Docker + Compose
    ├── kubernetes/     # K8s YAML + Helm
    └── sync-test/      # Sync test cases
```

---

## 🗺️ Screenshots

> Platform screenshots (to be added)

<details>
<summary>📸 Click to expand screenshots</summary>

| Feature | Preview |
|---------|---------|
| Dashboard | *(screenshot to add)* |
| Project - Kanban View | *(screenshot to add)* |
| Personnel Workload Matrix | *(screenshot to add)* |
| Swagger API Docs | *(screenshot to add)* |

> 💡 How to add screenshots: Drop image files into the root or `docs/` directory and update this section in your PR.

</details>

---

## 📧 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/budongshu/duptwo/issues)
- 📘 **Docs**: [DEPLOY.md](DEPLOY.md) | Built-in Swagger (`/swagger`)
- 🔧 **API Base**: `http://localhost:18421/api`

---

## 📄 License

MIT © [budongshu](https://github.com/budongshu)