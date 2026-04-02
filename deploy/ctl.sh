#!/bin/bash
# ============================================================
# duptwo 管理脚本
# 支持 systemctl 和直接后台运行
# ============================================================

set -e

APP_NAME="duptwo"
BINARY="./${APP_NAME}"
CONFIG="./conf/app.yaml"
PID_FILE="./data/${APP_NAME}.pid"
LOG_FILE="./logs/${APP_NAME}.log"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

usage() {
    echo -e "${BLUE}duptwo 管理脚本${NC}"
    echo ""
    echo "用法: $0 <命令> [选项]"
    echo ""
    echo "命令:"
    echo "  start       启动服务"
    echo "  stop        停止服务"
    echo "  restart     重启服务"
    echo "  status      查看状态"
    echo "  log         查看日志"
    echo "  install     安装 systemd 服务（需 root）"
    echo "  uninstall   卸载 systemd 服务（需 root）"
    echo ""
    echo "选项:"
    echo "  --no-daemon   直接前台运行（用于调试）"
}

# 创建必要目录
init_dirs() {
    mkdir -p ./data ./logs
}

# 检查进程是否运行
is_running() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            return 0
        fi
    fi
    # 也检查进程名
    if pgrep -x "$APP_NAME" > /dev/null 2>&1; then
        return 0
    fi
    return 1
}

# 获取进程 PID
get_pid() {
    if [ -f "$PID_FILE" ]; then
        cat "$PID_FILE"
    else
        pgrep -x "$APP_NAME" 2>/dev/null || echo ""
    fi
}

# 启动服务
do_start() {
    init_dirs

    if is_running; then
        PID=$(get_pid)
        echo -e "${YELLOW}duptwo 已在运行 (PID: $PID)${NC}"
        return 1
    fi

    # 检查二进制文件
    if [ ! -f "$BINARY" ]; then
        echo -e "${RED}错误: 未找到 $BINARY${NC}"
        echo "请在正确目录执行此脚本"
        exit 1
    fi

    # 检查配置文件
    if [ ! -f "$CONFIG" ]; then
        echo -e "${RED}错误: 未找到配置文件 $CONFIG${NC}"
        exit 1
    fi

    echo -e "${GREEN}启动 duptwo...${NC}"

    if [ "$1" == "--no-daemon" ]; then
        echo -e "${YELLOW}前台运行模式（Ctrl+C 退出）${NC}"
        $BINARY --config "$CONFIG"
    else
        # 后台运行
        nohup $BINARY --config "$CONFIG" >> "$LOG_FILE" 2>&1 &
        PID=$!
        echo $PID > "$PID_FILE"
        sleep 1

        if is_running; then
            echo -e "${GREEN}duptwo 已启动 (PID: $PID)${NC}"
            echo "日志: $LOG_FILE"
        else
            echo -e "${RED}启动失败，请查看日志: $LOG_FILE${NC}"
            cat "$LOG_FILE"
            exit 1
        fi
    fi
}

# 停止服务
do_stop() {
    if ! is_running; then
        echo -e "${YELLOW}duptwo 未运行${NC}"
        rm -f "$PID_FILE"
        return
    fi

    PID=$(get_pid)
    echo -e "${YELLOW}停止 duptwo (PID: $PID)...${NC}"

    # 优雅停止（发送 SIGTERM）
    kill "$PID" 2>/dev/null || true

    # 等待最多 10 秒
    for i in {1..10}; do
        if ! kill -0 "$PID" 2>/dev/null; then
            break
        fi
        sleep 1
    done

    # 强制杀死
    if kill -0 "$PID" 2>/dev/null; then
        echo "强制停止..."
        kill -9 "$PID" 2>/dev/null || true
    fi

    rm -f "$PID_FILE"
    echo -e "${GREEN}duptwo 已停止${NC}"
}

# 重启服务
do_restart() {
    echo -e "${YELLOW}重启 duptwo...${NC}"
    do_stop
    sleep 1
    do_start
}

# 查看状态
do_status() {
    if is_running; then
        PID=$(get_pid)
        echo -e "${GREEN}● duptwo 运行中 (PID: $PID)${NC}"
    else
        echo -e "${RED}○ duptwo 未运行${NC}"
    fi
}

# 查看日志
do_log() {
    if [ ! -f "$LOG_FILE" ]; then
        echo -e "${YELLOW}日志文件不存在${NC}"
        return
    fi

    if [ "$1" == "-f" ] || [ "$1" == "--follow" ]; then
        echo -e "${BLUE}跟踪日志 (Ctrl+C 退出)${NC}"
        tail -f "$LOG_FILE"
    else
        tail -50 "$LOG_FILE"
    fi
}

# 安装 systemd 服务
do_install() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}错误: 安装服务需要 root 权限${NC}"
        echo "请使用: sudo $0 install"
        exit 1
    fi

    # 获取当前目录的绝对路径
    CURRENT_DIR=$(cd "$(dirname "$0")" && pwd)

    SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

    echo -e "${YELLOW}安装 systemd 服务...${NC}"

    # 生成 service 文件
    sed -e "s|{{USER}}|$(whoami)|g" \
        -e "s|{{WORKDIR}}|${CURRENT_DIR}|g" \
        -e "s|{{BINDIR}}|${CURRENT_DIR}|g" \
        -e "s|{{CONFDIR}}|${CURRENT_DIR}/conf|g" \
        -e "s|{{DATADIR}}|${CURRENT_DIR}/data|g" \
        -e "s|{{LOGDIR}}|${CURRENT_DIR}/logs|g" \
        "${CURRENT_DIR}/deploy/duptwo.service" > /tmp/${APP_NAME}.service

    mv /tmp/${APP_NAME}.service "$SERVICE_FILE"
    chmod 644 "$SERVICE_FILE"

    systemctl daemon-reload
    systemctl enable "${APP_NAME}.service"

    init_dirs

    echo -e "${GREEN}安装完成！${NC}"
    echo ""
    echo "使用以下命令管理服务:"
    echo "  sudo systemctl start ${APP_NAME}    # 启动"
    echo "  sudo systemctl stop ${APP_NAME}     # 停止"
    echo "  sudo systemctl restart ${APP_NAME}  # 重启"
    echo "  sudo systemctl status ${APP_NAME}   # 状态"
    echo "  sudo journalctl -u ${APP_NAME} -f    # 日志"
}

# 卸载 systemd 服务
do_uninstall() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}错误: 卸载服务需要 root 权限${NC}"
        echo "请使用: sudo $0 uninstall"
        exit 1
    fi

    echo -e "${YELLOW}卸载 systemd 服务...${NC}"

    systemctl stop "${APP_NAME}.service" 2>/dev/null || true
    systemctl disable "${APP_NAME}.service" 2>/dev/null || true
    rm -f "/etc/systemd/system/${APP_NAME}.service"
    systemctl daemon-reload

    echo -e "${GREEN}卸载完成${NC}"
}

# -------------------- 主程序 --------------------
COMMAND=${1:-usage}
shift || true

case $COMMAND in
    start)
        do_start "$@"
        ;;
    stop)
        do_stop
        ;;
    restart)
        do_restart
        ;;
    status)
        do_status
        ;;
    log)
        do_log "$@"
        ;;
    install)
        do_install
        ;;
    uninstall)
        do_uninstall
        ;;
    help|--help|-h)
        usage
        ;;
    *)
        usage
        exit 1
        ;;
esac
