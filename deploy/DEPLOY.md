# duptwo 数据登记平台 - 部署指南

> **版本**: v1.0.0
> **最后更新**: 2026-04-01

## 快速开始

### Docker 一键部署（推荐）

```bash
# 进入部署目录
cd deploy/docker

# SQLite 版（轻量，推荐开发/小规模使用）
docker compose -f docker-compose.sqlite.yml up -d

# MySQL 版（推荐生产环境）
docker compose -f docker-compose.mysql.yml up -d

# PostgreSQL 版（高并发生产环境）
docker compose -f docker-compose.postgres.yml up -d
```

启动后访问：
- **Web UI**: http://localhost:80
- **API 直连**: http://localhost:18421
- **健康检查**: http://localhost:18421/health
- **Swagger API 文档**: http://localhost:18421/swagger
- **默认账号**: `admin` / `admin123`

### 二进制部署

```bash
# 1. 下载 release 包
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo

# 2. 修改配置
vim conf/app.yaml

# 3. 启动
./duptwo --config conf/app.yaml
```

---

## 系统架构

```
                           ┌─────────────────┐
                           │   用户浏览器      │
                           └────────┬────────┘
                                    │ HTTP :80 / HTTPS :443
                           ┌────────▼────────┐
                           │     Nginx       │  ← 反向代理 + SSL（可选）
                           │                 │
                           └────────┬────────┘
                                    │
                           ┌────────▼────────┐
                           │   Go API        │
                           │   :18421        │
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
- 支持 SQLite、MySQL、PostgreSQL 三种数据库
- Nginx 用于 HTTPS 终结和负载均衡（可选）

---

## 环境要求

### 硬件要求

| 场景 | CPU | 内存 | 磁盘 |
|------|-----|------|------|
| 开发/测试 | 2核 | 4GB | 20GB |
| 小规模生产（SQLite） | 2核 | 4GB | 50GB |
| 中等规模生产（MySQL） | 4核 | 8GB | 100GB |
| 大规模生产（MySQL） | 8核+ | 16GB+ | 200GB+ |

### 软件要求

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Docker | 20.10 | 24.0+ |
| Docker Compose | 2.0 | 2.20+ |
| MySQL（可选） | 8.0 | 8.0 LTS |
| PostgreSQL（可选） | 13 | 16 LTS |
| Go（开发/二进制构建） | 1.21 | 1.23 |
| Node.js（前端构建） | 18.0 | 20.0 LTS |

---

## 部署方式

### Docker 部署

#### SQLite 轻量版（推荐）

适用于：开发、小规模使用、个人/团队内部系统

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

**数据备份**：
```bash
# 备份 SQLite 数据库
docker cp duptwo-app:/app/data/registry.db ./backup/

# 或直接复制
docker cp duptwo-app:/app/data/registry.db ./registry.db.$(date +%Y%m%d)
```

#### MySQL 生产版

适用于：生产环境、多用户、高并发

```bash
cd deploy/docker

# 启动（首次启动会自动创建数据库）
docker compose -f docker-compose.mysql.yml up -d

# 查看 MySQL 日志
docker compose -f docker-compose.mysql.yml logs -f mysql

# 查看应用日志
docker compose -f docker-compose.mysql.yml logs -f app
```

**MySQL 数据备份**：
```bash
# 备份数据库
docker exec duptwo-mysql mysqldump -u root -pduptwo_root_pass duptwo > backup_$(date +%Y%m%d).sql

# 恢复数据库
docker exec -i duptwo-mysql mysql -u root -pduptwo_root_pass duptwo < backup_20260401.sql
```

#### PostgreSQL 生产版

适用于：生产环境、高并发、强一致性需求

```bash
cd deploy/docker

# 启动
docker compose -f docker-compose.postgres.yml up -d

# 查看日志
docker compose -f docker-compose.postgres.yml logs -f

# 停止
docker compose -f docker-compose.postgres.yml down
```

**PostgreSQL 数据备份**：
```bash
# 备份
docker exec duptwo-postgres pg_dump -U duptwo duptwo > backup_$(date +%Y%m%d).sql

# 恢复
docker exec -i duptwo-postgres psql -U duptwo duptwo < backup_20260401.sql
```

#### 二进制部署

##### Linux amd64

```bash
# 1. 下载 release 包
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64

# 2. 修改配置
vim conf/app.yaml

# 3. 创建数据目录
mkdir -p data logs

# 4. 启动
./duptwo --config conf/app.yaml
```

##### 初始化管理员密码

```bash
# 重置 admin 密码
./duptwo reset-admin your_new_password

# 查看帮助
./duptwo --help
```

---

## 配置说明

### 配置文件位置

- SQLite 生产配置: `deploy/docker/conf/app.sqlite.yaml`
- MySQL 生产配置: `deploy/docker/conf/app.mysql.yaml`
- PostgreSQL 生产配置: `deploy/docker/conf/app.postgres.yaml`
- 开发配置: `backend/conf/app.yaml`

### SQLite 配置示例

```yaml
base:
  mode: prod
  port: 18421
  serve_web: true

database:
  type: sqlite
  path: ./data/registry.db

jwt:
  secret: CHANGE_THIS_TO_A_LONG_RANDOM_SECRET  # ⚠️ 必须修改
  expire_hours: 168  # 7天

session:
  secret: CHANGE_THIS_TO_ANOTHER_LONG_SECRET
```

### MySQL 配置示例

```yaml
database:
  type: mysql
  host: mysql
  port: 3306
  user: duptwo
  pass: duptwo_secure_pass  # ⚠️ 修改密码
  name: duptwo
```

### 环境变量覆盖

Docker 部署时可通过环境变量覆盖配置：

```bash
# SQLite 模式
docker run -e DATABASE_TYPE=sqlite -e DATABASE_PATH=/data/registry.db ...

# MySQL 模式
docker run -e DATABASE_TYPE=mysql -e DB_HOST=mysql -e DB_PORT=3306 ...
```

---

## 数据库支持

| 数据库 | 适用场景 | 数据存储 | 备份方式 |
|--------|---------|---------|---------|
| **SQLite** | 开发、小规模 (< 100用户) | 单文件 `./data/registry.db` | `cp` 或 `tar` 打包 |
| **MySQL** | 生产环境、多用户 | MySQL 数据库 | `mysqldump` |
| **PostgreSQL** | 生产环境、多用户 | PostgreSQL 数据库 | `pg_dump` |

### 从 SQLite 迁移到 MySQL

1. 导出 SQLite 数据
2. 导入 MySQL
3. 修改配置 `database.type: mysql`

---

## 生产环境检查清单

### 安全配置 ⚠️

- [ ] JWT Secret 已修改为强随机字符串（32+ 字符）
- [ ] Admin 默认密码已修改（`./duptwo reset-admin <新密码>`）
- [ ] MySQL 密码已修改（MySQL 版）
- [ ] CORS 白名单已配置
- [ ] 生产环境设置 `mode: prod`

### 网络配置

- [ ] 配置了 HTTPS（Let's Encrypt 或商业证书）
- [ ] 防火墙已开放必要端口（80, 443）
- [ ] 数据库端口（3306）仅内网访问

### 数据安全

- [ ] 配置了定时备份策略
- [ ] 备份存储与生产分离

---

## 常见问题

### Q: 启动报 `port already in use`

```bash
# 查看占用端口的进程
lsof -i :18421

# 杀掉占用进程或修改配置中的端口
```

### Q: MySQL 连接失败

```bash
# 检查容器网络
docker network ls
docker compose -f docker-compose.mysql.yml exec mysql mysql -u duptwo -pduptwo_secure_pass

# 检查连接字符串
docker compose -f docker-compose.mysql.yml exec app env | grep DB_
```

### Q: 如何完全重置系统？

```bash
# 停止并删除所有容器和数据
docker compose -f docker-compose.sqlite.yml down -v

# 重新启动（会自动初始化）
docker compose -f docker-compose.sqlite.yml up -d
```

### Q: 二进制部署如何更新？

```bash
# 1. 备份数据
cp -r data data.backup

# 2. 停止服务
pkill duptwo

# 3. 替换二进制
mv duptwo duptwo.old
cp new-duptwo duptwo

# 4. 启动
./duptwo --config conf/app.yaml
```

### Q: 查看详细日志？

```bash
# Docker
docker compose -f docker-compose.sqlite.yml logs -f --tail=100 app

# 二进制
./duptwo --log-level=debug
```

---

## 技术支持

- GitHub Issues: https://github.com/budongshu/duptwo/issues
- 文档: https://github.com/budongshu/duptwo#readme
