#!/bin/bash
# 启动 Center（接收 Agent 推送）

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 杀掉旧的 center 进程
pkill -f "datauptwo_center" 2>/dev/null

# 创建数据目录
mkdir -p ./data_center

# 启动 center
echo "[Center] 启动中..."
./datauptwo_center --config ./conf/app.center.yaml > ./logs/center.log 2>&1 &
echo "[Center] 进程 PID=$!，日志: ./logs/center.log"

sleep 2

# 检查是否启动成功
if curl -s http://127.0.0.1:18421/health > /dev/null 2>&1; then
    echo "[Center] 启动成功，监听 http://127.0.0.1:18421"
else
    echo "[Center] 启动失败，查看日志：tail -f ./logs/center.log"
    exit 1
fi
