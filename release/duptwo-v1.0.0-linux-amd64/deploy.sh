#!/bin/bash
# ============================================================
# duptwo 数据登记平台 一键部署脚本
# 产品名称: duptwo / 数据登记平台
# 配置统一在 deploy/config.yaml 中管理
# ============================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 产品信息
PRODUCT_NAME="duptwo"
PRODUCT_CN="数据登记平台"
IMAGE_NAME="${IMAGE_NAME:-duptwo}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
NAMESPACE="${NAMESPACE:-duptwo}"

# ========== 颜色定义 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "${RED}[ERR]${NC}  $*"; }
log_step()  { echo -e "${BLUE}[STEP]${NC}  $*" ; }
log_ok()    { echo -e "${CYAN}[ OK ]${NC}  $*"; }

# ========== 依赖检查 ==========
check_deps() {
    log_step "检查依赖..."
    local missing=""
    for cmd in docker python3; do
        if ! command -v $cmd &> /dev/null; then
            missing="$missing $cmd"
        fi
    done
    # 检查 docker compose
    if ! docker compose version &> /dev/null && ! docker compose --version &> /dev/null; then
        missing="$missing docker compose"
    fi
    if [ -n "$missing" ]; then
        log_error "缺少依赖:$missing"
        log_info "安装方法:"
        log_info "  Ubuntu/Debian: apt update && apt install -y docker.io docker compose python3"
        log_info "  CentOS/RHEL:   yum install -y docker docker compose python3"
        log_info "  macOS:         brew install docker docker compose python3"
        exit 1
    fi
    local docker_ver=$(docker --version 2>/dev/null | grep -oP '\d+\.\d+\.\d+' | head -1)
    log_info "Docker: $docker_ver"
    log_ok "依赖检查通过"
}

# ========== 生成配置 ==========
gen_config() {
    log_step "从 config.yaml 生成各环境配置..."
    if [ -f "config_gen.py" ]; then
        python3 config_gen.py --all 2>&1 | while read line; do
            echo -e "         $line"
        done
    fi
    log_ok "配置生成完成"
}

# ========== Docker Compose 封装 ==========
docker_compose() {
    cd "$SCRIPT_DIR/docker"
    if docker compose version &> /dev/null; then
        docker compose "$@"
    else
        docker compose "$@"
    fi
}

# ========== 构建 Docker 镜像 ==========
build_docker() {
    log_step "构建 Docker 镜像: ${IMAGE_NAME}:${IMAGE_TAG}"

    # 生成配置
    gen_config

    # 构建前端
    if [ -d "../frontend" ]; then
        log_info "构建前端..."
        (cd ../frontend && npm run build 2>&1 | tail -5)
    fi

    # 生成 Swagger 文档
    log_info "生成 Swagger 文档..."
    (cd ../backend && \
        (command -v swag &>/dev/null || go install github.com/swaggo/swag/cmd/swag@latest) && \
        swag init --parseInternal --parseDependency ./docs 2>&1 | tail -3 || true)

    # 构建镜像
    log_info "构建镜像: ${IMAGE_NAME}:${IMAGE_TAG}"
    docker build \
        -t "${IMAGE_NAME}:${IMAGE_TAG}" \
        -f docker/Dockerfile \
        .. 2>&1 | tail -5

    local size=$(docker images "${IMAGE_NAME}:${IMAGE_TAG}" --format '{{.Size}}' 2>/dev/null)
    log_ok "镜像构建完成！大小: $size"
}

# ========== Docker 启动（MySQL 版）==========
docker_up() {
    log_step "启动 Docker 服务（MySQL 版）..."
    cd "$SCRIPT_DIR/docker"

    docker_compose -f docker-compose.mysql.yml up -d --remove-orphans
    sleep 5

    log_info ""
    log_ok "=========================================="
    log_ok "  $PRODUCT_CN ($PRODUCT_NAME) 部署成功！"
    log_ok "=========================================="
    log_info "  Web UI:    http://localhost:80"
    log_info "  API:       http://localhost:80/api"
    log_info "  Swagger:   http://localhost:80/swagger"
    log_info "  健康检查:  http://localhost:80/health"
    log_info ""
    log_info "默认账号: admin / admin123"
    log_info ""
    docker_compose -f docker-compose.mysql.yml ps
}

# ========== Docker 启动（PostgreSQL 版）==========
docker_up_postgres() {
    log_step "启动 Docker 服务（PostgreSQL 版）..."
    cd "$SCRIPT_DIR/docker"

    log_info "使用 PostgreSQL 生产版"
    docker_compose -f docker-compose.postgres.yml up -d --remove-orphans
    sleep 5

    log_ok "=========================================="
    log_ok "  $PRODUCT_CN ($PRODUCT_NAME) PostgreSQL 版部署成功！"
    log_ok "=========================================="
    log_info "  Web UI:    http://localhost:80"
    log_info "  API:       http://localhost:80/api"
    log_info "  Swagger:   http://localhost:80/swagger"
    log_info "  健康检查:  http://localhost:80/health"
    log_info ""
    log_info "默认账号: admin / admin123"
    log_info ""
    docker_compose -f docker-compose.postgres.yml ps
}

# ========== Docker 启动（SQLite 版）==========
docker_up_sqlite() {
    log_step "启动 Docker 服务（SQLite 轻量版）..."
    cd "$SCRIPT_DIR/docker"

    log_info "使用 SQLite 轻量版，无需 MySQL"
    docker_compose -f docker-compose.sqlite.yml up -d --remove-orphans
    sleep 5

    log_ok "=========================================="
    log_ok "  $PRODUCT_CN ($PRODUCT_NAME) SQLite 版部署成功！"
    log_ok "=========================================="
    log_info "  Web UI:    http://localhost:8080"
    log_info "  Swagger:   http://localhost:8080/swagger"
    log_info "  健康检查:  http://localhost:8080/health"
    log_info ""
    log_info "默认账号: admin / admin123"
    log_info ""
    docker_compose -f docker-compose.sqlite.yml ps
}

# ========== Docker 停止 ==========
docker_down() {
    log_step "停止 Docker 服务..."
    cd "$SCRIPT_DIR/docker"
    docker_compose -f docker-compose.mysql.yml down 2>/dev/null || true
    docker_compose -f docker-compose.sqlite.yml down 2>/dev/null || true
    docker_compose -f docker-compose.postgres.yml down 2>/dev/null || true
    log_ok "服务已停止"
}

# ========== Docker 重启 ==========
docker_restart() {
    docker_down
    docker_up
}

# ========== K8s 部署 ==========
k8s_deploy() {
    log_step "部署到 Kubernetes..."

    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi

    gen_config

    log_warn "镜像需先推送到仓库，请修改 manifests.yaml 中的镜像地址："
    log_warn "  your-registry.com/duptwo:${IMAGE_TAG}"
    echo ""

    kubectl apply -f "$SCRIPT_DIR/kubernetes/namespace.yaml"

    # 替换命名空间
    sed "s/namespace: duptwo/namespace: $NAMESPACE/g" \
        "$SCRIPT_DIR/kubernetes/manifests.yaml" | \
        kubectl apply -f - 2>/dev/null || \
        kubectl apply -f "$SCRIPT_DIR/kubernetes/manifests.yaml"

    log_info "等待 Pod 就绪（最多 2 分钟）..."
    kubectl wait --for=condition=ready pod -l app=duptwo -n "$NAMESPACE" --timeout=120s 2>/dev/null || true

    kubectl get pods,svc,ingress -n "$NAMESPACE" 2>/dev/null || \
        kubectl get pods,svc -n "$NAMESPACE"

    log_ok "K8s 部署完成！"
    log_info "查看日志: kubectl -n $NAMESPACE logs -l app=duptwo -f"
}

# ========== K8s 删除 ==========
k8s_delete() {
    log_step "从 Kubernetes 删除..."
    kubectl delete -f "$SCRIPT_DIR/kubernetes/manifests.yaml" 2>/dev/null || true
    kubectl delete namespace "$NAMESPACE" 2>/dev/null || true
    log_ok "已删除"
}

# ========== 进入容器 ==========
docker_shell() {
    cd "$SCRIPT_DIR/docker"
    docker_compose exec app sh
}

# ========== 查看日志 ==========
logs() {
    cd "$SCRIPT_DIR/docker"
    docker_compose logs -f --tail=50
}

# ========== 查看状态 ==========
status() {
    cd "$SCRIPT_DIR/docker"
    docker_compose ps 2>/dev/null || echo "Docker Compose 未运行"

    echo ""
    echo -e "${BLUE}健康检查:${NC}"
    local health=$(curl -s --connect-timeout 3 http://localhost:80/health 2>/dev/null)
    if [ -n "$health" ]; then
        echo "$health" | python3 -m json.tool 2>/dev/null || echo "$health"
        log_ok "服务正常运行"
    else
        log_warn "服务未启动或无法访问"
    fi
}

# ========== 完全清理 ==========
clean() {
    log_warn "清理所有容器和数据（不可恢复）..."
    cd "$SCRIPT_DIR/docker"
    docker_compose -f docker-compose.mysql.yml down -v --remove-orphans 2>/dev/null || true
    docker_compose -f docker-compose.sqlite.yml down -v --remove-orphans 2>/dev/null || true
    docker_compose -f docker-compose.postgres.yml down -v --remove-orphans 2>/dev/null || true
    log_ok "清理完成"
}

# ========== 帮助 ==========
usage() {
    echo -e "${CYAN}duptwo 数据登记平台 - 一键部署脚本${NC}"
    echo ""
    echo "用法: $0 <命令>"
    echo ""
    echo -e "${GREEN}Docker 部署:${NC}"
    echo "  $0 deploy          构建 + 部署（MySQL 版，推荐生产）"
    echo "  $0 deploy:sqlite   构建 + 部署（SQLite 轻量版）"
    echo "  $0 deploy:postgres  构建 + 部署（PostgreSQL 生产版）"
    echo "  $0 build          仅构建镜像"
    echo "  $0 start          仅启动（使用已有镜像，MySQL）"
    echo "  $0 start:sqlite   仅启动 SQLite 版"
    echo "  $0 start:postgres 仅启动 PostgreSQL 版"
    echo "  $0 stop           停止所有容器"
    echo "  $0 restart        重启"
    echo ""
    echo -e "${GREEN}查看和调试:${NC}"
    echo "  $0 logs           实时查看日志"
    echo "  $0 status         查看服务状态"
    echo "  $0 shell          进入容器终端"
    echo ""
    echo -e "${GREEN}Kubernetes:${NC}"
    echo "  $0 k8s:deploy     部署到 K8s"
    echo "  $0 k8s:delete     从 K8s 删除"
    echo ""
    echo -e "${GREEN}配置:${NC}"
    echo "  $0 config         仅生成配置文件"
    echo "  $0 clean          完全清理（慎用）"
    echo ""
    echo -e "${GREEN}环境变量:${NC}"
    echo "  IMAGE_NAME=duptwo IMAGE_TAG=v2.0.0 $0 deploy"
    echo ""
    echo -e "${GREEN}快速开始:${NC}"
    echo "  $0 config          # 先生成配置"
    echo "  $0 deploy          # MySQL 版完整部署"
    echo "  $0 deploy:sqlite   # SQLite 版（无需 MySQL）"
    echo "  $0 deploy:postgres  # PostgreSQL 版（高并发）"
}

case "${1:-help}" in
    config)          check_deps && gen_config ;;
    build)           check_deps && build_docker ;;
    deploy)          check_deps && build_docker && docker_up ;;
    deploy:sqlite)   check_deps && build_docker && docker_up_sqlite ;;
    deploy:postgres)  check_deps && build_docker && docker_up_postgres ;;
    start)           docker_up ;;
    start:sqlite)    docker_up_sqlite ;;
    start:postgres)  docker_up_postgres ;;
    stop)            docker_down ;;
    restart)         docker_restart ;;
    logs)            logs ;;
    status)          status ;;
    shell)           docker_shell ;;
    clean)           clean ;;
    k8s:deploy)      check_deps && gen_config && k8s_deploy ;;
    k8s:delete)      k8s_delete ;;
    help|--help|-h)  usage ;;
    *)               usage ;;
esac
