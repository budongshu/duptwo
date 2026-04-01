#!/bin/bash
# ============================================================
# duptwo v1.0.0 Release 打包脚本
# 支持：Linux amd64 (Docker + Binary)
# ============================================================

set -e

VERSION=${1:-v1.0.0}
BUILD_DIR="release"
PROJECT_NAME="duptwo"
ARCH="amd64"
PLATFORM="linux"

echo "=========================================="
echo "  duptwo ${VERSION} Release Builder"
echo "=========================================="

# 清理旧文件
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 1. 构建二进制包 (Linux amd64)
echo ""
echo "[1/4] 构建二进制包..."
cd backend

# 构建 Linux 版本
GOOS=linux GOARCH=${ARCH} CGO_ENABLED=1 CC=gcc \
    go build -ldflags="-s -w -X main.version=${VERSION}" \
    -o ../${BUILD_DIR}/${PROJECT_NAME} ./cmd/server/main.go

# 复制前端
cp -r cmd/server/web ../${BUILD_DIR}/web

# 复制默认配置
mkdir -p ../${BUILD_DIR}/conf
cp conf/app.yaml ../${BUILD_DIR}/conf/app.yaml 2>/dev/null || true

cd ..

# 2. 打包二进制 release
echo ""
echo "[2/4] 打包二进制 release..."
cd ${BUILD_DIR}
tar -czvf ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}.tar.gz \
    ${PROJECT_NAME} \
    web/ \
    conf/ \
    README.md \
    2>/dev/null || true

# 创建精简包（只有二进制和配置）
mkdir -p ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal
cp ${PROJECT_NAME} ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal/
cp conf/app.yaml ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal/
cd ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal
tar -czvf ../${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal.tar.gz .
cd ..
rm -rf ${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-minimal
cd ..

# 3. 构建 Docker 镜像
echo ""
echo "[3/4] 构建 Docker 镜像..."

# 构建 SQLite 版本镜像
docker build \
    -t ${PROJECT_NAME}:${VERSION}-sqlite \
    -t ${PROJECT_NAME}:latest-sqlite \
    -f deploy/docker/Dockerfile .

# 标记为 MySQL 版本（使用相同的镜像，配置不同）
docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:${VERSION}
docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:${VERSION}-mysql
docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:latest

# 4. 导出 Docker 镜像为 tgz
echo ""
echo "[4/4] 导出 Docker 镜像..."

# 导出完整镜像
docker save ${PROJECT_NAME}:${VERSION} \
    -o ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar

# 压缩
gzip ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar

# 创建 Docker 部署包
mkdir -p ${BUILD_DIR}/docker-deploy
cp ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar.gz ${BUILD_DIR}/docker-deploy/

# 复制 docker-compose 文件
mkdir -p ${BUILD_DIR}/docker-deploy/docker
cp deploy/docker/docker-compose.yml ${BUILD_DIR}/docker-deploy/docker/
cp deploy/docker/docker-compose.sqlite.yml ${BUILD_DIR}/docker-deploy/docker/
cp deploy/docker/docker-compose.mysql.yml ${BUILD_DIR}/docker-deploy/docker/
cp -r deploy/docker/nginx ${BUILD_DIR}/docker-deploy/docker/
cp -r deploy/docker/conf ${BUILD_DIR}/docker-deploy/docker/

# 复制 README
cat > ${BUILD_DIR}/docker-deploy/README.md << 'README_EOF'
# duptwo v1.0.0 Docker 部署包

## 快速开始

### 1. 加载 Docker 镜像
```bash
tar -xzf duptwo-v1.0.0-docker.tar.gz
docker load -i duptwo-v1.0.0-docker.tar
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

- **Web UI**: http://localhost:80
- **API 直连**: http://localhost:18421
- **健康检查**: http://localhost:18421/health
- **Swagger**: http://localhost:18421/swagger

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

## 数据备份

### SQLite 备份
```bash
docker cp duptwo-app:/app/data/registry.db ./backup/
```

### MySQL 备份
```bash
docker exec duptwo-mysql mysqldump -u root -pduptwo_root_pass duptwo > backup.sql
```

## 配置说明

配置文件位于 `docker/conf/` 目录：
- `app.sqlite.yaml` - SQLite 配置
- `app.mysql.yaml` - MySQL 配置

修改配置后重新启动：
```bash
docker compose -f docker-compose.sqlite.yml restart
```
README_EOF

# 打包 Docker 部署包
cd ${BUILD_DIR}
tar -czvf ${PROJECT_NAME}-${VERSION}-docker.tar.gz docker-deploy/
rm -rf docker-deploy/
cd ..

# 显示结果
echo ""
echo "=========================================="
echo "  Release 构建完成！"
echo "=========================================="
echo ""
echo "输出文件:"
ls -lh ${BUILD_DIR}/
echo ""
echo "推荐部署方式:"
echo ""
echo "1. Docker 部署（推荐）:"
echo "   tar -xzf ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar.gz"
echo "   cd docker && docker compose -f docker-compose.sqlite.yml up -d"
echo ""
echo "2. 二进制部署:"
echo "   tar -xzf ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}.tar.gz"
echo "   cd ${PROJECT_NAME} && ./${PROJECT_NAME} --config conf/app.yaml"
echo ""
