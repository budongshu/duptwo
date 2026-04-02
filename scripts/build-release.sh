#!/bin/bash
# ============================================================
# duptwo Release 打包脚本
# 用法: ./build-release.sh [binary|docker|all] [version] [--static]
#
# 示例:
#   ./build-release.sh binary          # 标准版
#   ./build-release.sh binary --static # 静态版（兼容旧系统）
#   ./build-release.sh all --static     # 全部构建
# ============================================================

set -e

# 解析参数
STATIC=false
TARGET="all"
VERSION="v1.0.0"

for arg in "$@"; do
    case $arg in
        --static)
            STATIC=true
            ;;
        binary|docker|all)
            TARGET=$arg
            ;;
        v*)
            VERSION=$arg
            ;;
    esac
done

BUILD_DIR="release"
PROJECT_NAME="duptwo"
ARCH="amd64"
PLATFORM="linux"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}==========================================${NC}"
echo -e "${GREEN}  duptwo ${VERSION} Release Builder${NC}"
echo -e "${GREEN}==========================================${NC}"
echo "模式: $TARGET"
echo "静态编译: $STATIC"
echo ""

# 清理旧文件
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# -------------------- 二进制包 --------------------
build_binary() {
    echo -e "${YELLOW}[1/2] 构建二进制包...${NC}"

    cd backend

    if [ "$STATIC" = true ]; then
        echo "使用静态编译（兼容旧系统）..."
        # 静态编译：禁用 CGO，不依赖 glibc
        CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
            go build -ldflags="-s -w -X main.version=${VERSION}" \
            -o ../${BUILD_DIR}/${PROJECT_NAME} ./cmd/server/main.go
    else
        echo "使用标准编译..."
        # 标准编译：启用 CGO，支持 SQLite 原生驱动
        GOOS=linux GOARCH=${ARCH} CGO_ENABLED=1 CC=gcc \
            go build -ldflags="-s -w -X main.version=${VERSION}" \
            -o ../${BUILD_DIR}/${PROJECT_NAME} ./cmd/server/main.go
    fi

    mkdir -p ../${BUILD_DIR}/web
    cp -r cmd/server/web ../${BUILD_DIR}/web/

    mkdir -p ../${BUILD_DIR}/conf
    cp conf/app.yaml ../${BUILD_DIR}/conf/app.yaml 2>/dev/null || true

    cd ..

    # 根据是否静态编译调整包名
    if [ "$STATIC" = true ]; then
        PKG_NAME="${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}-static.tar.gz"
    else
        PKG_NAME="${PROJECT_NAME}-${VERSION}-${PLATFORM}-${ARCH}.tar.gz"
    fi

    echo "打包..."
    cd ${BUILD_DIR}
    tar -czvf ${PKG_NAME} \
        ${PROJECT_NAME} web/ conf/

    cd ..

    echo -e "${GREEN}[OK] 二进制包: ${BUILD_DIR}/${PKG_NAME}${NC}"
}

# -------------------- Docker 镜像 --------------------
build_docker() {
    echo -e "${YELLOW}[2/2] 构建 Docker 镜像...${NC}"

    docker build \
        -t ${PROJECT_NAME}:${VERSION}-sqlite \
        -t ${PROJECT_NAME}:latest-sqlite \
        -f deploy/docker/Dockerfile .

    docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:${VERSION}
    docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:${VERSION}-mysql
    docker tag ${PROJECT_NAME}:${VERSION}-sqlite ${PROJECT_NAME}:latest

    echo "导出镜像..."
    docker save ${PROJECT_NAME}:${VERSION} \
        -o ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar
    gzip ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar

    echo "创建部署包..."
    mkdir -p ${BUILD_DIR}/docker-deploy/docker
    cp ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar.gz ${BUILD_DIR}/docker-deploy/
    cp deploy/docker/docker-compose.yml ${BUILD_DIR}/docker-deploy/docker/
    cp deploy/docker/docker-compose.sqlite.yml ${BUILD_DIR}/docker-deploy/docker/
    cp deploy/docker/docker-compose.mysql.yml ${BUILD_DIR}/docker-deploy/docker/
    cp -r deploy/docker/nginx ${BUILD_DIR}/docker-deploy/docker/
    cp -r deploy/docker/conf ${BUILD_DIR}/docker-deploy/docker/

    cd ${BUILD_DIR}
    tar -czvf ${PROJECT_NAME}-${VERSION}-docker.tar.gz docker-deploy/
    rm -rf docker-deploy/
    cd ..

    echo -e "${GREEN}[OK] Docker 包: ${BUILD_DIR}/${PROJECT_NAME}-${VERSION}-docker.tar.gz${NC}"
}

# -------------------- 执行 --------------------
case $TARGET in
    binary)
        build_binary
        ;;
    docker)
        build_docker
        ;;
    all)
        build_binary
        echo ""
        build_docker
        ;;
esac

# 显示结果
echo ""
echo -e "${GREEN}==========================================${NC}"
echo -e "${GREEN}  构建完成！${NC}"
echo -e "${GREEN}==========================================${NC}"
echo ""
ls -lh ${BUILD_DIR}/
