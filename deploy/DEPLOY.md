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

如果二进制无法运行（如 glibc 版本问题），使用 Docker/nerdctl：

```bash
cd DataRegistry
docker compose -f deploy/docker/docker-compose.sqlite.yml up -d
# 或 containerd 环境：
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml up -d
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
  # web_root: ./web      # 前端静态文件目录（可选，留空自动推导）
  # 配置示例：
  #   ./web               - 相对于项目根目录（release 包默认）
  #   ./cmd/server/web    - 相对于项目根目录（开发时使用）
  #   /etc/duptwo/web     - 绝对路径

database:
  type: sqlite           # sqlite / mysql / postgres
  path: ./data/registry.db

jwt:
  secret: CHANGE_THIS   # ⚠️ 必须修改为长随机字符串
  expire_hours: 168

session:
  secret: CHANGE_THIS   # ⚠️ 必须修改
```

### 目录结构与路径关系

```
项目根目录/                     ← web_root 的参照目录
├── duptwo                   # 可执行文件
├── web/                     # 前端静态文件（release 包默认）
│   ├── index.html
│   └── assets/
├── cmd/server/web/          # 前端静态文件（开发源码目录）
│   ├── index.html
│   └── assets/
└── conf/
    └── app.yaml             # 配置文件（web_root 相对于 conf/ 的父目录）

# web_root 为空时自动推导：
#   1. 先找 ./web（release 包）
#   2. 再找 ./cmd/server/web（开发目录）
#   3. 均未找到则不加载前端（API 仍可用）
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
cd DataRegistry
docker compose -f deploy/docker/docker-compose.sqlite.yml up -d
docker compose -f deploy/docker/docker-compose.sqlite.yml logs -f
```

备份：`docker cp duptwo-app:/app/data/registry.db ./`

### MySQL 生产版

```bash
cd DataRegistry
docker compose -f deploy/docker/docker-compose.mysql.yml up -d
```

备份：`docker exec duptwo-mysql mysqldump -u root -p duptwo > backup.sql`

### Containerd 环境部署（nerdctl）

服务器使用 containerd（无 Docker daemon）时，用 nerdctl 代替：

```bash
# 1. 安装 nerdctl（阿里云镜像下载）
wget https://mirrors.aliyun.com/docker-toolbox/linux/tools/nerdctl-1.7.6-linux-amd64.tar.gz
tar -xzf nerdctl-1.7.6-linux-amd64.tar.gz -C /usr/local/bin/
nerdctl --version

# 2. 在有 Docker 的机器上构建镜像并导出（推荐在本地/CI机器上执行）
docker build -t duptwo:v1.0.0 -f deploy/docker/Dockerfile .
docker save duptwo:v1.0.0 -o duptwo-v1.0.0-docker.tar

# 3. 上传 tar 包到服务器，导入镜像（阿里云镜像仓库加速）
# 方式A：直接导入 tar 包
nerdctl load -i duptwo-v1.0.0-docker.tar

# 方式B：推送到阿里云 ACR（推荐生产环境）
# 先在 https://cr.aliyun.com 创建命名空间
nerdctl tag duptwo:v1.0.0 registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0
nerdctl push registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0

# 4. 启动（从项目根目录）
cd DataRegistry
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml up -d

# 5. 查看状态
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml ps
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml logs -f
```

---

## Kubernetes 部署（阿里云 ACK）

### 前提条件

1. 阿里云 ACK 集群（Kubernetes 1.22+）
2. kubectl 已配置集群访问
3. 镜像已推送到阿里云 ACR

### 方式一：原生 Manifest 部署（SQLite 版本）

```bash
cd DataRegistry

# 1. 修改镜像地址
# 编辑 deploy/kubernetes/01-manifests.yaml
# 将 image: duptwo:v1.0.0 替换为你的阿里云镜像
# 例如: registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0

# 2. 生成安全密钥（替换默认值）
cat > /tmp/secrets.env << EOF
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)
EOF

# 3. 替换 manifest 中的占位符
sed -i 's/REPLACE_WITH_SESSION_SECRET_32CHARS/'"$(openssl rand -base64 32)"'/g' deploy/kubernetes/01-manifests.yaml
sed -i 's/REPLACE_WITH_JWT_SECRET_32CHARS/'"$(openssl rand -base64 32)"'/g' deploy/kubernetes/01-manifests.yaml

# 4. 替换域名为你的实际域名
sed -i 's/duptwo.your-domain.com/你的实际域名/g' deploy/kubernetes/01-manifests.yaml

# 5. 部署
kubectl apply -f deploy/kubernetes/

# 6. 查看状态
kubectl -n duptwo get pods,svc,ingress

# 7. 查看日志
kubectl -n duptwo logs -l app=duptwo -f
```

### 方式二：Helm 部署（推荐）

```bash
cd DataRegistry

# 1. 构建并推送镜像到阿里云 ACR
docker build -t duptwo:v1.0.0 -f deploy/docker/Dockerfile .
docker tag duptwo:v1.0.0 registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0
docker push registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0

# 2. 安装 Helm Chart
cd deploy/kubernetes/helm-duptwo

# 生成密钥
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)

# 安装
helm install duptwo . \
  --namespace duptwo \
  --create-namespace \
  --set image.repository=registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo \
  --set image.tag=v1.0.0 \
  --set ingress.host=duptwo.你的域名.com \
  --set secret.sessionSecret="$SESSION_SECRET" \
  --set secret.jwtSecret="$JWT_SECRET"

# 3. 升级（如已安装）
helm upgrade duptwo . \
  --namespace duptwo \
  --set image.tag=v1.1.0

# 4. 卸载
helm uninstall duptwo -n duptwo
```

### 配置说明

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `replicaCount` | 2 | Pod 副本数 |
| `persistence.size` | 10Gi | SQLite 数据卷大小 |
| `persistence.storageClass` | "" | 存储类（空=默认，可设为 nfs-csi） |
| `resources.requests.cpu` | 100m | 最小 CPU |
| `resources.requests.memory` | 256Mi | 最小内存 |
| `resources.limits.cpu` | 1000m | 最大 CPU |
| `resources.limits.memory` | 1Gi | 最大内存 |

### K8s 架构说明

```
┌─────────────────────────────────────────────────────────────┐
│                      Ingress (nginx)                        │
│                   duptwo.your-domain.com                    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service (ClusterIP)                     │
│                      port: 18421                            │
└─────────────────────────────────────────────────────────────┘
                            │
              ┌─────────────┴─────────────┐
              ▼                           ▼
┌─────────────────────┐       ┌─────────────────────┐
│   Pod (duptwo-1)    │       │   Pod (duptwo-2)   │
│  image: duptwo:v1   │       │  image: duptwo:v1  │
│  port: 18421        │       │  port: 18421       │
└─────────────────────┘       └─────────────────────┘
              │                           │
              └─────────────┬─────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│               PVC (duptwo-data-pvc) - 10Gi                 │
│                    /app/data (SQLite)                      │
└─────────────────────────────────────────────────────────────┘
```

### 数据备份

```bash
# 使用 kubectl cp 备份数据库
kubectl -n duptwo exec deploy/duptwo -- tar -czf /tmp/backup.tar.gz /app/data
kubectl -n duptwo cp deploy/duptwo:/tmp/backup.tar.gz ./duptwo-backup-$(date +%Y%m%d).tar.gz
```

### 常见问题

**Pod 无法启动**
```bash
kubectl -n duptwo describe pod -l app=duptwo
kubectl -n duptwo logs -l app=duptwo --previous
```

**PVC 绑定失败**
- 检查集群是否有足够的存储
- 确认 StorageClass 存在：`kubectl get storageclass`

**Ingress 无法访问**
- 检查域名 DNS 解析
- 确认 Ingress Controller 运行正常

---

## 数据同步系统（Data Sync）

### 系统架构

duptwo 数据同步系统采用 **Center + Agent** 架构，实现多节点数据汇聚：

```
                          ┌──────────────────────────────────────┐
                          │           Center 节点                  │
                          │         (数据汇聚中心)                  │
                          │                                      │
                          │  ┌────────────────────────────────┐   │
                          │  │  /api/v1/sync/stations        │   │
                          │  │  /api/v1/sync/history         │   │
                          │  │  /api/v1/sync/status          │   │
                          │  │  /api/v1/sync/register        │   │
                          │  │  /api/v1/sync/upload          │   │
                          │  └────────────────────────────────┘   │
                          │                                      │
                          │  ┌─────────────┐ ┌────────────────┐  │
                          │  │ Station Mgmt │ │ History Store  │  │
                          │  └─────────────┘ └────────────────┘  │
                          └──────────────────────────────────────┘
                                    ▲          ▲          ▲
                                    │          │          │
                          ┌─────────┴┐   ┌─────┴────┐ ┌──┴───────┐
                          │ Agent 1  │   │ Agent 2  │ │ Agent N  │
                          │ 站点A    │   │ 站点B    │ │ 站点...  │
                          └──────────┘   └──────────┘ └──────────┘
```

- **Center 节点**：数据汇聚中心，负责站点管理、API Key 分配、记录存储
- **Agent 节点**：各数据采集站点，独立运行，定期推送数据到 Center

### 配置说明

#### Center 模式配置

作为数据汇聚中心启用：

```yaml
sync:
  enabled: true
  mode: "center"    # 作为中心节点运行
  center_url: ""    # Center 模式下留空
  api_key: ""       # 自动生成
  station_id: ""
  station_name: ""
  interval: "5m"
  retry_count: 3
  retry_interval: "1m"
  proxy:
    enabled: false
    url: ""
    username: ""
    password: ""
```

#### Agent 模式配置

作为数据采集站启用：

```yaml
sync:
  enabled: true
  mode: "agent"     # 作为采集站运行
  center_url: "http://center-node:8080"  # Center 节点地址
  api_key: ""       # Agent 注册后由 Center 下发
  station_id: ""    # Agent 注册后由 Center 下发
  station_name: "北京站点"
  interval: "5m"    # 同步间隔
  retry_count: 3
  retry_interval: "1m"
  proxy:
    enabled: false
    url: "http://proxy:8080"  # 如需代理
    username: ""              # 代理认证用户名
    password: ""              # 代理认证密码
```

### API Key 说明

| 场景 | 说明 |
|------|------|
| **Agent 注册** | 调用 `POST /api/v1/sync/register`，Center 返回 `api_key` 和 `station_id` |
| **Agent 上传** | 请求头携带 `X-API-Key: <your_api_key>`，调用 `POST /api/v1/sync/upload` |
| **管理操作** | 所有 `/api/v1/sync/stations` 和 `/api/v1/sync/history` 接口需 JWT 认证 |

### 代理支持

Agent 模式下支持 HTTP/HTTPS 代理，适用于内网环境：

```yaml
sync:
  proxy:
    enabled: true
    url: "http://proxy.example.com:8080"
    username: "proxy_user"
    password: "proxy_pass"
```

### 相关 API

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/v1/sync/stations` | JWT | 站点列表 |
| POST | `/api/v1/sync/stations` | JWT | 创建站点 |
| PUT | `/api/v1/sync/stations/{id}` | JWT | 更新站点 |
| DELETE | `/api/v1/sync/stations/{id}` | JWT | 删除站点 |
| GET | `/api/v1/sync/history` | JWT | 同步历史 |
| GET | `/api/v1/sync/status` | JWT | 同步状态 |
| GET | `/api/v1/sync/history/{id}/details` | JWT | 同步详情 |
| POST | `/api/v1/sync/register` | None | Agent 注册 |
| POST | `/api/v1/sync/upload` | API Key | Agent 上传记录 |

---

## 常见问题

**glibc 版本不兼容**
```
./duptwo: /lib64/libc.so.6: version `GLIBC_2.33' not found
```

解决方式：
1. 使用 Docker 部署：`nerdctl compose -f deploy/docker/docker-compose.sqlite.yml up -d`
2. 在新系统上重新编译：`go build -ldflags="-s -w" -o duptwo ./cmd/server/main.go`

**端口被占用**
```bash
lsof -i :18421
```

**重置系统**
```bash
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml down -v
nerdctl compose -f deploy/docker/docker-compose.sqlite.yml up -d
```

---

## 技术支持

- GitHub: https://github.com/budongshu/duptwo
- Issues: https://github.com/budongshu/duptwo/issues
