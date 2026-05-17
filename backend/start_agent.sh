#!/bin/bash
# 启动 Agent（推送到 Center）

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 杀掉旧的 agent 进程
pkill -f "datauptwo_agent" 2>/dev/null

# 创建数据目录
mkdir -p ./data_agent

# 启动 agent
echo "[Agent] 启动中..."
./datauptwo_agent --config ./conf/app.agent.yaml > ./logs/agent.log 2>&1 &
echo "[Agent] 进程 PID=$!，日志: ./logs/agent.log"

sleep 3

# 检查是否启动成功
if curl -s http://127.0.0.1:18422/health > /dev/null 2>&1; then
    echo "[Agent] 启动成功，监听 http://127.0.0.1:18422"
    echo "[Agent] 推送目标: http://localhost:18421"
else
    echo "[Agent] 启动失败，查看日志：tail -f ./logs/agent.log"
    exit 1
fi
