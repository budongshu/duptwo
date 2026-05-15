# duptwo 数据登记平台 - 部署指南

> **版本**: v1.0.0

## 快速开始

### 二进制部署（推荐）

```bash
# 1. 下载 release 包
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64

# 2. 创建目录
mkdir -p data logs

# 3. 修改配置（必读！）
vim conf/app.yaml

# 4. 启动
./duptwo --config conf/app.yaml
```

### Docker 部署

```bash
cd deploy/docker

# SQLite 版（轻量，开发/小规模）
docker compose -f docker-compose.sqlite.yml up -d

# MySQL 版（生产环境）
docker compose -f docker-compose.mysql.yml up -d
```

### 启动后访问

- **Web UI**: http://localhost:18421
- **API 文档**: http://localhost:18421/swagger
- **健康检查**: http://localhost:18421/health
- **默认账号**: `admin` / `admin123`

---

## 目录结构

```
DataRegistry/
├── deploy/           # 部署相关
│   ├── deploy.sh     # 一键部署脚本
│   ├── docker/       # Docker 配置
│   ├── kubernetes/   # K8s 配置
│   └── clients/      # API 客户端示例
├── scripts/          # 构建脚本
│   └── build-release.sh  # 打包二进制
└── release/          # 发布产物
    └── duptwo-v1.0.0-linux-amd64.tar.gz
```

---

## 二进制部署详解

### 下载地址

```bash
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
```

### 配置说明

修改 `conf/app.yaml`：

```yaml
base:
  mode: prod
  port: 18421
  serve_web: true

database:
  type: sqlite           # sqlite / mysql / postgres
  path: ./data/registry.db

jwt:
  secret: CHANGE_THIS   # ⚠️ 必须修改为长随机字符串
  expire_hours: 168

session:
  secret: CHANGE_THIS   # ⚠️ 必须修改
```

### 命令行

```bash
# 重置 admin 密码
./duptwo reset-admin your_new_password

# 查看帮助
./duptwo --help

# 调试模式
./duptwo --log-level=debug
```

### 更新版本

```bash
# 1. 备份
cp -r data data.backup

# 2. 停止
pkill duptwo

# 3. 替换
mv duptwo duptwo.old
cp new-duptwo duptwo

# 4. 重启
./duptwo --config conf/app.yaml
```

---

## Docker 部署详解

### SQLite 轻量版

```bash
cd deploy/docker

docker compose -f docker-compose.sqlite.yml up -d
docker compose -f docker-compose.sqlite.yml logs -f
```

备份：`docker cp duptwo-app:/app/data/registry.db ./`

### MySQL 生产版

```bash
cd deploy/docker
docker compose -f docker-compose.mysql.yml up -d
```

备份：`docker exec duptwo-mysql mysqldump -u root -p duptwo > backup.sql`

---

## 常见问题

**端口被占用**
```bash
lsof -i :18421
```

**重置系统**
```bash
docker compose -f docker-compose.sqlite.yml down -v
docker compose -f docker-compose.sqlite.yml up -d
```

---

## 技术支持

- GitHub: https://github.com/budongshu/duptwo
- Issues: https://github.com/budongshu/duptwo/issues
