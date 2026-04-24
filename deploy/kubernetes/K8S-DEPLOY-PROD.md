# DataRegistry K8s 生产环境部署文档

> **版本**: v1.1.0 (MySQL 生产版)
> **Namespace**: infra
> **数据库**: MySQL 云 PaaS (阿里云 RDS)

---

## 目录

- [架构概览](#架构概览)
- [前提条件](#前提条件)
- [快速部署](#快速部署)
- [详细配置](#详细配置)
- [部署清单](#部署清单)
- [验证部署](#验证部署)
- [常见问题](#常见问题)

---

## 架构概览

```
                          ┌─────────────────────────────────────────┐
                          │         Ingress (nginx)                  │
                          │        duptwo.your-domain.com           │
                          └─────────────────────────────────────────┘
                                            │
                                            ▼
                          ┌─────────────────────────────────────────┐
                          │         Service (ClusterIP)              │
                          │           port: 18421                    │
                          └─────────────────────────────────────────┘
                                            │
                              ┌─────────────┴─────────────┐
                              ▼                           ▼
                    ┌─────────────────┐       ┌─────────────────┐
                    │   Pod (duptwo-1)│       │   Pod (duptwo-2)│
                    │  image: duptwo  │       │  image: duptwo  │
                    │  port: 18421   │       │  port: 18421    │
                    └─────────────────┘       └─────────────────┘
                              │                           │
                              └─────────────┬─────────────┘
                                            ▼
                          ┌─────────────────────────────────────────┐
                          │         MySQL 云 PaaS (RDS)             │
                          │    rm-xxxxx.mysql.rds.aliyuncs.com      │
                          └─────────────────────────────────────────┘
```

---

## 前提条件

1. **Kubernetes 集群** (v1.22+)
2. **kubectl** 已配置集群访问
3. **Helm 3** 已安装
4. **Ingress Controller** 已部署 (nginx-ingress)
5. **Docker 镜像** 已推送到镜像仓库
6. **MySQL 云 PaaS** 已创建 (阿里云 RDS)

---

## 快速部署

### 1. 构建并推送镜像

```bash
cd DataRegistry

# 构建镜像
docker build -t duptwo:v1.0.0 -f deploy/docker/Dockerfile .

# 推送到阿里云 ACR
docker tag duptwo:v1.0.0 registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0
docker push registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo:v1.0.0
```

### 2. 创建 infra namespace

```bash
kubectl create namespace infra
```

### 3. 配置并安装

```bash
cd DataRegistry/deploy/kubernetes/helm-duptwo

# 生成安全密钥
SESSION_SECRET=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)

# 安装 Helm Chart
helm install duptwo . \
  --namespace infra \
  --create-namespace \
  --set image.repository=registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo \
  --set image.tag=v1.0.0 \
  --set ingress.host=duptwo.your-domain.com \
  --set secret.sessionSecret="$SESSION_SECRET" \
  --set secret.jwtSecret="$JWT_SECRET" \
  --set mysql.externalHost=rm-xxxxx.mysql.rds.aliyuncs.com \
  --set mysql.password=your_mysql_password
```

### 4. 验证部署

```bash
# 查看 Pod 状态
kubectl -n infra get pods

# 查看服务
kubectl -n infra get svc,ingress

# 查看日志
kubectl -n infra logs -l app.kubernetes.io/name=duptwo -f
```

---

## 详细配置

### values.yaml 配置项

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `replicaCount` | 2 | Pod 副本数 |
| `image.repository` | duptwo | 镜像仓库地址 |
| `image.tag` | v1.0.0 | 镜像版本 |
| `ingress.host` | duptwo.your-domain.com | 访问域名 |
| `mysql.externalHost` | "" | MySQL 云 PaaS 地址 |
| `mysql.database` | duptwo | 数据库名 |
| `mysql.username` | duptwo | 数据库用户名 |
| `mysql.password` | "" | 数据库密码（必填）|
| `secret.sessionSecret` | "" | Session 密钥（必填）|
| `secret.jwtSecret` | "" | JWT 密钥（必填）|
| `resources.requests.cpu` | 100m | 最小 CPU |
| `resources.requests.memory` | 256Mi | 最小内存 |
| `resources.limits.cpu` | 1000m | 最大 CPU |
| `resources.limits.memory` | 1Gi | 最大内存 |

### 使用外部配置文件

创建 `production.yaml`:

```yaml
replicaCount: 2

image:
  repository: registry.cn-hangzhou.aliyuncs.com/你的命名空间/duptwo
  tag: v1.0.0
  pullPolicy: Always

ingress:
  enabled: true
  className: nginx
  host: duptwo.your-domain.com
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
    nginx.ingress.kubernetes.io/proxy-body-size: "500m"

mysql:
  externalHost: rm-xxxxx.mysql.rds.aliyuncs.com
  database: duptwo
  username: duptwo
  password: "your_mysql_password"

secret:
  sessionSecret: "your_session_secret_32_chars"
  jwtSecret: "your_jwt_secret_32_chars"

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 100m
    memory: 256Mi
```

安装:
```bash
helm install duptwo . -f production.yaml --namespace infra --create-namespace
```

---

## 部署清单

### 资源列表

| 资源类型 | 名称 | 说明 |
|---------|------|------|
| Deployment | duptwo | 应用部署 (2 副本) |
| Service | duptwo-duptwo-svc | ClusterIP 服务 |
| ConfigMap | duptwo-config | 应用配置 |
| Secret | duptwo-secret | 敏感信息 (密钥、数据库密码) |
| Ingress | duptwo-ingress | HTTP 入口 |

### 部署顺序

1. Secret (数据库密码、密钥)
2. ConfigMap (应用配置)
3. Deployment (应用)
4. Service (服务)
5. Ingress (入口)

---

## 验证部署

### 健康检查

```bash
# 获取 Ingress IP
kubectl -n infra get ingress

# 健康检查
curl http://duptwo.your-domain.com/health

# 预期响应: {"status":"ok"}
```

### 访问服务

- **Web UI**: http://duptwo.your-domain.com
- **API 文档**: http://duptwo.your-domain.com/swagger
- **健康检查**: http://duptwo.your-domain.com/health
- **默认账号**: `admin` / `admin123`

### 日志查看

```bash
# 查看所有 Pod 日志
kubectl -n infra logs -l app.kubernetes.io/name=duptwo --tail=100

# 实时跟踪
kubectl -n infra logs -l app.kubernetes.io/name=duptwo -f

# 查看特定 Pod
kubectl -n infra logs duptwo-xxxxx-xxxx
```

---

## 运维操作

### 升级

```bash
# 更新镜像版本
helm upgrade duptwo . \
  --namespace infra \
  --set image.tag=v1.1.0

# 使用自定义配置
helm upgrade duptwo . -f production.yaml --namespace infra
```

### 回滚

```bash
# 查看发布历史
helm history duptwo -n infra

# 回滚到上一个版本
helm rollback duptwo -n infra

# 回滚到指定版本
helm rollback duptwo 1 -n infra
```

### 卸载

```bash
helm uninstall duptwo -n infra
kubectl delete namespace infra
```

---

## 常见问题

### Pod 无法启动

```bash
# 查看 Pod 详情
kubectl -n infra describe pod -l app.kubernetes.io/name=duptwo

# 查看日志
kubectl -n infra logs -l app.kubernetes.io/name=duptwo --previous
```

### 数据库连接失败

1. 确认 MySQL 云 PaaS 地址正确
2. 确认数据库账号密码正确
3. 确认云 PaaS 已开放连接白名单
4. 测试连接: `telnet rm-xxxxx.mysql.rds.aliyuncs.com 3306`

### Ingress 无法访问

1. 确认域名 DNS 解析正确
2. 确认 Ingress Controller 运行正常
3. 检查 Ingress 注解配置

### 资源不足

```bash
# 扩展副本数
helm upgrade duptwo . --namespace infra --set replicaCount=3

# 调整资源限制
helm upgrade duptwo . --namespace infra --set resources.limits.memory=2Gi
```

---

## 技术支持

- GitHub: https://github.com/budongshu/duptwo
- Issues: https://github.com/budongshu/duptwo/issues
