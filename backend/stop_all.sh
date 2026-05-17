#!/bin/bash
# 停止所有服务

echo "[Stop] 正在停止 Center 和 Agent..."
pkill -f "datauptwo_center" 2>/dev/null && echo "[Stop] Center 已停止" || echo "[Stop] Center 未运行"
pkill -f "datauptwo_agent" 2>/dev/null && echo "[Stop] Agent 已停止" || echo "[Stop] Agent 未运行"
echo "[Stop] 完成"
