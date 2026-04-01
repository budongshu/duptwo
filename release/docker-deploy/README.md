# duptwo v1.0.0 Docker 部署包

## 快速开始

### 1. 加载 Docker 镜像

```bash
# 先构建镜像（如果已有镜像则跳过）
cd docker
docker build -t duptwo:v1.0.0 -f Dockerfile ../..

# 导出/导入镜像
docker save duptwo:v1.0.0 -o duptwo-v1.0.0-docker.tar
# 在目标机器上: docker load -i duptwo-v1.0.0-docker.tar
```

### 2. SQLite 轻量版部署

```bash
cd docker
docker compose -f docker-compose.sqlite.yml up -d
```

### 3. MySQL 生产版部署

```bash
cd docker
docker compose -f docker-compose.mysql.yml up -d
```

## 访问地址

| 服务 | 地址 |
|------|------|
| Web UI | http://localhost:80 |
| API 直连 | http://localhost:18421 |
| 健康检查 | http://localhost:18421/health |
| Swagger | http://localhost:18421/swagger |

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

## 配置说明

配置文件位于 `docker/conf/` 目录：
- `app.sqlite.yaml` - SQLite 配置（轻量版）
- `app.mysql.yaml` - MySQL 配置（生产版）

## 数据备份

### SQLite
```bash
docker cp duptwo-app:/app/data/registry.db ./backup/
```

### MySQL
```bash
docker exec duptwo-mysql mysqldump -u root -pduptwo_root_pass duptwo > backup.sql
```

## 停止服务

```bash
# SQLite
docker compose -f docker-compose.sqlite.yml down

# MySQL（保留数据）
docker compose -f docker-compose.mysql.yml down

# MySQL（删除数据）
docker compose -f docker-compose.mysql.yml down -v
```
