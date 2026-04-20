#!/bin/bash
# ============================================================
# duptwo 数据同步系统 - 一键部署脚本
# 同时部署 Center 和 Agent 节点进行测试
#
# 用法:
#   ./deploy-sync.sh              # 交互式选择
#   ./deploy-sync.sh center        # 仅部署 Center
#   ./deploy-sync.sh agent         # 仅部署 Agent
#   ./deploy-sync.sh all           # 同时部署 Center + Agent
#   ./deploy-sync.sh status        # 查看状态
#   ./deploy-sync.sh stop          # 停止全部
#   ./deploy-sync.sh clean         # 清理全部
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 项目根目录（脚本所在目录的父目录）
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 部署目录
CENTER_DIR="${SCRIPT_DIR}/sync-test/center"
AGENT_DIR="${SCRIPT_DIR}/sync-test/agent"

# 端口配置
CENTER_APP_PORT=18421
CENTER_WEB_PORT=18080
AGENT_APP_PORT=18422
AGENT_WEB_PORT=18081

# Center 地址（Agent 用来连接）
CENTER_URL="http://localhost:${CENTER_APP_PORT}"

# 日志函数
log_info()   { echo -e "${BLUE}[INFO]${NC}  $*"; }
log_ok()     { echo -e "${GREEN}[OK]${NC}   $*"; }
log_warn()   { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $*"; }

# 打印标题
print_title() {
    echo ""
    echo -e "${CYAN}====================================================${NC}"
    echo -e "${CYAN}  $*${NC}"
    echo -e "${CYAN}====================================================${NC}"
}

# 检查 Docker 是否可用
check_docker() {
    if command -v docker &> /dev/null; then
        log_ok "Docker 已安装: $(docker --version 2>/dev/null | cut -d' ' -f3)"
    elif command -v nerdctl &> /dev/null; then
        log_ok "nerdctl 已安装"
        alias docker=nerdctl 2>/dev/null || true
    else
        log_error "未找到 Docker 或 nerdctl"
        exit 1
    fi

    # 检查 docker compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    elif docker-compose --version &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        log_error "未找到 docker compose"
        exit 1
    fi
    log_ok "使用命令: $COMPOSE_CMD"
}

# 构建 Docker 镜像
build_image() {
    print_title "构建 Docker 镜像"
    cd "$PROJECT_ROOT"

    log_info "开始构建镜像（可能需要几分钟）..."
    if docker build -t duptwo:sync-test -f deploy/docker/Dockerfile .; then
        log_ok "镜像构建成功: duptwo:sync-test"
    else
        log_error "镜像构建失败"
        exit 1
    fi
}

# 创建目录结构
create_dirs() {
    log_info "创建部署目录..."
    mkdir -p "${CENTER_DIR}/conf" "${CENTER_DIR}/data" "${CENTER_DIR}/logs"
    mkdir -p "${AGENT_DIR}/conf" "${AGENT_DIR}/data" "${AGENT_DIR}/logs"
    log_ok "目录创建完成"
}

# 生成 Center 配置
generate_center_config() {
    cat > "${CENTER_DIR}/conf/app.yaml" << 'EOF'
# Center 站点配置 - 数据同步中心
base:
  mode: prod
  port: 18421
  serve_web: true
  web_root: /app/web

database:
  type: sqlite
  path: /app/data/center.db

log:
  level: info
  time_zone: Asia/Shanghai
  log_name: center
  max_backup: 10
  max_size: 50

session:
  timeout: 7200
  secret: center-session-secret-change-in-production

cors:
  allow_origins:
    - "http://localhost"
    - "http://localhost:18080"
    - "http://localhost:18421"
  allow_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allow_headers:
    - Origin
    - Content-Type
    - Authorization
    - X-API-Key
  max_age: 86400

jwt:
  secret: center-jwt-secret-change-in-production

# ========== 数据同步配置 ==========
sync:
  enabled: true
  mode: "center"
  interval: "1m"
  retry_count: 3
  retry_interval: "30s"
  proxy:
    enabled: false
EOF
    log_ok "Center 配置已生成"
}

# 生成 Agent 配置
generate_agent_config() {
    cat > "${AGENT_DIR}/conf/app.yaml" << 'EOF'
# Agent 站点配置 - 数据采集端
base:
  mode: prod
  port: 18422
  serve_web: true
  web_root: /app/web

database:
  type: sqlite
  path: /app/data/agent.db

log:
  level: info
  time_zone: Asia/Shanghai
  log_name: agent
  max_backup: 10
  max_size: 50

session:
  timeout: 7200
  secret: agent-session-secret-change-in-production

cors:
  allow_origins:
    - "http://localhost"
    - "http://localhost:18081"
    - "http://localhost:18422"
  allow_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allow_headers:
    - Origin
    - Content-Type
    - Authorization
    - X-API-Key
  max_age: 86400

jwt:
  secret: agent-jwt-secret-change-in-production

# ========== 数据同步配置 ==========
sync:
  enabled: true
  mode: "agent"
  center_url: "REPLACE_CENTER_URL"
  station_id: "agent-local-01"
  station_name: "本地测试Agent"
  api_key: "REPLACE_API_KEY"
  interval: "30s"
  retry_count: 3
  retry_interval: "10s"
  proxy:
    enabled: false
EOF
    log_ok "Agent 配置已生成"
}

# 生成 Center Docker Compose
generate_center_docker() {
    cat > "${CENTER_DIR}/docker-compose.yml" << EOF
# Center 站点 - Docker Compose
version: '3.8'

services:
  center:
    image: duptwo:sync-test
    container_name: duptwo-sync-center
    restart: unless-stopped
    ports:
      - "${CENTER_APP_PORT}:18421"
      - "${CENTER_WEB_PORT}:8080"
    volumes:
      - ./conf/app.yaml:/app/conf/app.yaml:ro
      - ./data:/app/data
      - ./logs:/app/logs
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:18421/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - sync-network

networks:
  sync-network:
    driver: bridge
EOF
    log_ok "Center Docker Compose 已生成"
}

# 生成 Agent Docker Compose
generate_agent_docker() {
    cat > "${AGENT_DIR}/docker-compose.yml" << EOF
# Agent 站点 - Docker Compose
version: '3.8'

services:
  agent:
    image: duptwo:sync-test
    container_name: duptwo-sync-agent
    restart: unless-stopped
    ports:
      - "${AGENT_APP_PORT}:18422"
      - "${AGENT_WEB_PORT}:8080"
    volumes:
      - ./conf/app.yaml:/app/conf/app.yaml:ro
      - ./data:/app/data
      - ./logs:/app/logs
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:18422/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - sync-network

networks:
  sync-network:
    driver: bridge
EOF
    log_ok "Agent Docker Compose 已生成"
}

# 部署 Center
deploy_center() {
    print_title "部署 Center 站点"

    create_dirs
    generate_center_config
    generate_center_docker

    cd "${CENTER_DIR}"
    log_info "启动 Center 容器..."
    $COMPOSE_CMD up -d

    log_info "等待 Center 启动..."
    sleep 5

    # 检查健康状态
    for i in {1..10}; do
        if curl -s "http://localhost:${CENTER_APP_PORT}/health" | grep -q "ok"; then
            log_ok "Center 启动成功!"
            break
        fi
        log_info "等待中... ($i/10)"
        sleep 2
    done

    echo ""
    log_ok "=========================================="
    log_ok "  Center 部署完成!"
    log_ok "=========================================="
    echo "  应用端口:   http://localhost:${CENTER_APP_PORT}"
    echo "  Web 界面:   http://localhost:${CENTER_WEB_PORT}"
    echo "  健康检查:   http://localhost:${CENTER_APP_PORT}/health"
    echo "  数据目录:   ${CENTER_DIR}/data"
    echo "  日志目录:   ${CENTER_DIR}/logs"
    echo ""
    log_info "请访问 Web 界面，在【系统设置 → 数据同步】中添加 Agent 站点"
    log_info "添加后会生成 API Key，需要复制给 Agent 配置使用"
    echo ""
}

# 部署 Agent
deploy_agent() {
    print_title "部署 Agent 站点"

    create_dirs
    generate_agent_config
    generate_agent_docker

    cd "${AGENT_DIR}"
    log_info "启动 Agent 容器..."
    $COMPOSE_CMD up -d

    log_info "等待 Agent 启动..."
    sleep 5

    # 检查健康状态
    for i in {1..10}; do
        if curl -s "http://localhost:${AGENT_APP_PORT}/health" | grep -q "ok"; then
            log_ok "Agent 启动成功!"
            break
        fi
        log_info "等待中... ($i/10)"
        sleep 2
    done

    echo ""
    log_ok "=========================================="
    log_ok "  Agent 部署完成!"
    log_ok "=========================================="
    echo "  应用端口:   http://localhost:${AGENT_APP_PORT}"
    echo "  Web 界面:   http://localhost:${AGENT_WEB_PORT}"
    echo "  健康检查:   http://localhost:${AGENT_APP_PORT}/health"
    echo "  数据目录:   ${AGENT_DIR}/data"
    echo "  日志目录:   ${AGENT_DIR}/logs"
    echo ""
    log_info "Agent 已配置连接到: ${CENTER_URL}"
    log_info "注意: 需要在 Center 端添加站点并获取 API Key 后，才能完成同步"
    echo ""
}

# 查看状态
show_status() {
    print_title "数据同步系统状态"

    echo -e "\n${BLUE}[Center 站点]${NC}"
    if docker ps --filter "name=duptwo-sync-center" --format "{{.Status}}" 2>/dev/null | grep -q "Up"; then
        echo -e "  ${GREEN}● 运行中${NC}"
        echo "  应用: http://localhost:${CENTER_APP_PORT}"
        echo "  Web:  http://localhost:${CENTER_WEB_PORT}"

        # 获取 API Key
        if [ -f "${CENTER_DIR}/data/center.db" ]; then
            echo "  数据: ${CENTER_DIR}/data/center.db"
        fi

        # 检查同步配置
        echo "  同步模式: center (等待 Agent 连接)"
    else
        echo -e "  ${RED}○ 未运行${NC}"
        echo "  部署: ./deploy-sync.sh center"
    fi

    echo -e "\n${BLUE}[Agent 站点]${NC}"
    if docker ps --filter "name=duptwo-sync-agent" --format "{{.Status}}" 2>/dev/null | grep -q "Up"; then
        echo -e "  ${GREEN}● 运行中${NC}"
        echo "  应用: http://localhost:${AGENT_APP_PORT}"
        echo "  Web:  http://localhost:${AGENT_WEB_PORT}"
        echo "  数据: ${AGENT_DIR}/data/agent.db"
        echo "  同步模式: agent"
        echo "  目标: ${CENTER_URL}"

        # 检查 Agent 配置
        if grep -q "REPLACE_API_KEY" "${AGENT_DIR}/conf/app.yaml" 2>/dev/null; then
            echo -e "  ${YELLOW}⚠ API Key 未配置${NC} - 需要从 Center 获取"
        fi
    else
        echo -e "  ${RED}○ 未运行${NC}"
        echo "  部署: ./deploy-sync.sh agent"
    fi

    echo ""
}

# 停止服务
stop_services() {
    print_title "停止服务"

    if [ -d "${CENTER_DIR}" ]; then
        log_info "停止 Center..."
        cd "${CENTER_DIR}" && $COMPOSE_CMD down 2>/dev/null || true
        log_ok "Center 已停止"
    fi

    if [ -d "${AGENT_DIR}" ]; then
        log_info "停止 Agent..."
        cd "${AGENT_DIR}" && $COMPOSE_CMD down 2>/dev/null || true
        log_ok "Agent 已停止"
    fi
}

# 清理
clean_all() {
    print_title "清理全部"

    stop_services

    log_warn "即将删除所有数据..."
    read -p "确认删除? (yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        rm -rf "${SCRIPT_DIR}/sync-test"
        log_ok "已清理全部数据和配置"
    else
        log_info "取消清理"
    fi
}

# 配置 Agent API Key（交互式）
config_agent_apikey() {
    print_title "配置 Agent API Key"

    echo ""
    echo "请到 Center 管理界面添加站点并获取 API Key"
    echo ""
    read -p "请输入 API Key: " api_key
    read -p "请输入 Center URL [${CENTER_URL}]: " center_url

    center_url=${center_url:-${CENTER_URL}}

    if [ -z "$api_key" ]; then
        log_error "API Key 不能为空"
        return 1
    fi

    # 更新 Agent 配置
    sed -i "s|REPLACE_API_KEY|${api_key}|g" "${AGENT_DIR}/conf/app.yaml"
    sed -i "s|REPLACE_CENTER_URL|${center_url}|g" "${AGENT_DIR}/conf/app.yaml"

    log_ok "Agent 配置已更新"

    # 重启 Agent
    log_info "重启 Agent..."
    cd "${AGENT_DIR}" && $COMPOSE_CMD restart
    log_ok "Agent 已重启"
}

# 显示使用说明
show_usage() {
    print_title "数据同步测试环境部署"

    echo "
${GREEN}功能说明:${NC}
  本脚本用于部署 Center + Agent 同步测试环境

${GREEN}架构说明:${NC}
  ┌─────────────┐         ┌─────────────┐
  │  Agent 站点  │ ──push──>│  Center 站点 │
  │  :18422     │          │  :18421     │
  └─────────────┘          └─────────────┘

${GREEN}部署步骤:${NC}
  1. 部署 Center:   ./deploy-sync.sh center
  2. 访问 Center Web，在【系统设置 → 数据同步】添加站点
  3. 获取 API Key
  4. 配置 Agent:   ./deploy-sync.sh config-agent
  5. 查看状态:     ./deploy-sync.sh status

${GREEN}完整部署:${NC}
  ./deploy-sync.sh all     # 同时部署 Center + Agent

${GREEN}管理命令:${NC}
  ./deploy-sync.sh status   # 查看状态
  ./deploy-sync.sh stop     # 停止服务
  ./deploy-sync.sh clean    # 清理全部

${GREEN}端口分配:${NC}
  Center App:  ${CENTER_APP_PORT}  Center Web:  ${CENTER_WEB_PORT}
  Agent  App:  ${AGENT_APP_PORT}   Agent  Web:  ${AGENT_WEB_PORT}
"
}

# 主程序
main() {
    COMMAND=${1:-usage}

    check_docker

    case $COMMAND in
        center)
            deploy_center
            ;;
        agent)
            deploy_agent
            ;;
        all)
            build_image
            deploy_center
            deploy_agent
            show_status
            ;;
        status)
            show_status
            ;;
        stop)
            stop_services
            ;;
        clean)
            clean_all
            ;;
        config-agent)
            config_agent_apikey
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            show_usage
            ;;
    esac
}

main "$@"