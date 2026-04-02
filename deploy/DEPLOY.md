# duptwo 数据登记平台 - 部署指南

> **版本**: v1.0.0

## 快速开始

### 二进制部署（推荐）

```bash
# 1. 下载 release 包
wget https://github.com/budongshu/duptwo/releases/download/v1.0.0/duptwo-v1.0.0-linux-amd64.tar.gz
tar -xzf duptwo-v1.0.0-linux-amd64.tar.gz
cd duptwo-v1.0.0-linux-amd64

# 2. 修改配置
vim conf/app.yaml

# 3. 启动
./ctl.sh start
```

### Docker 部署（通用）

如果二进制无法运行（如 glibc 版本问题），使用 Docker：

```bash
cd deploy/docker
docker compose -f docker-compose.sqlite.yml up -d
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
│   ├── DEPLOY.md     # 本文档
│   ├── ctl.sh        # 服务管理脚本
│   ├── duptwo.service # systemd 服务文件
│   ├── docker/       # Docker 配置
│   ├── kubernetes/   # K8s 配置
│   └── clients/      # API 客户端示例
├── scripts/          # 构建脚本
│   └── build-release.sh  # 打包二进制
└── release/          # 发布产物
    └── duptwo-v1.0.0-linux-amd64.tar.gz
```

---

## 服务管理

### 方式一：使用 ctl.sh 脚本（推荐）

```bash
# 启动服务（后台运行）
./ctl.sh start

# 查看状态
./ctl.sh status

# 查看日志
./ctl.sh log

# 跟踪日志
./ctl.sh log -f

# 停止服务
./ctl.sh stop

# 重启服务
./ctl.sh restart

# 调试模式（前台运行）
./ctl.sh start --no-daemon
```

### 方式二：systemd 服务（生产环境推荐）

```bash
# 安装服务（需要 root）
sudo ./ctl.sh install

# 管理服务
sudo systemctl start duptwo
sudo systemctl stop duptwo
sudo systemctl restart duptwo
sudo systemctl status duptwo

# 查看日志
sudo journalctl -u duptwo -f

# 卸载服务
sudo ./ctl.sh uninstall
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
./ctl.sh stop
cp -r data data.backup

# 2. 替换二进制
cp new-duptwo duptwo

# 3. 重启
./ctl.sh start
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

**glibc 版本不兼容**
```
./duptwo: /lib64/libc.so.6: version `GLIBC_2.33' not found
```

解决方式：
1. 使用 Docker 部署：`docker compose -f deploy/docker/docker-compose.sqlite.yml up -d`
2. 在新系统上重新构建：`./scripts/build-release.sh binary --static`（静态编译版）

**端口被占用**
```bash
lsof -i :18421
```

**重置系统**
```bash
./ctl.sh stop
docker compose -f deploy/docker/docker-compose.sqlite.yml down -v
docker compose -f deploy/docker/docker-compose.sqlite.yml up -d
```

---

## 技术支持

- GitHub: https://github.com/budongshu/duptwo
- Issues: https://github.com/budongshu/duptwo/issues
