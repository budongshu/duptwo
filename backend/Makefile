# ============================================================
# datauptwo 构建与部署 Makefile
# ============================================================

.PHONY: help build build-linux build-docker build-all clean
.PHONY: run run-dev test fmt lint
.PHONY: docker-build docker-up docker-down docker-restart
.PHONY: k8s-apply k8s-delete k8s-logs
.PHONY: swagger-gen

# ========== 变量 ==========
APP_NAME    := datauptwo
VERSION     ?= v1.0.0
REGISTRY    ?= registry.example.com
IMAGE_NAME  := $(REGISTRY)/$(APP_NAME):$(VERSION)
BUILD_TIME  := $$(date +"%Y-%m-%d %H:%M:%S")
GO_LDFLAGS  := -s -w -X datauptwo/global.VERSION=$(VERSION) -X datauptwo/global.BUILD_TIME=$(BUILD_TIME)

# 前端目录
FRONTEND_DIR := ../frontend

# Docker 上下文
DOCKER_DIR := deploy/docker
DOCKER_COMPOSE := $(DOCKER_DIR)/docker-compose.yml

# K8s 目录
K8S_DIR := deploy/kubernetes

# ========== 帮助信息 ==========
help:
	@echo "DataRegistry 构建与部署命令"
	@echo ""
	@echo "=== 本地开发 ==="
	@echo "  make run-dev          启动开发模式（前端 Vite + 后端 Go）"
	@echo "  make test             运行测试"
	@echo "  make fmt              格式化代码"
	@echo "  make lint             代码检查"
	@echo ""
	@echo "=== 构建 ==="
	@echo "  make build            构建本地可执行文件"
	@echo "  make build-linux      构建 Linux amd64 可执行文件"
	@echo "  make build-docker     构建 Docker 镜像"
	@echo "  make build-all        构建前端 + 后端 + Docker 镜像"
	@echo "  make swagger-gen     生成 Swagger 文档"
	@echo ""
	@echo "=== Docker 部署 ==="
	@echo "  make docker-up        启动 Docker 服务（MySQL 版）"
	@echo "  make docker-up-sqlite 启动 Docker 服务（SQLite 版，轻量）"
	@echo "  make docker-down      停止并删除 Docker 容器"
	@echo "  make docker-restart   重启 Docker 服务"
	@echo "  make docker-logs      查看 Docker 日志"
	@echo ""
	@echo "=== Kubernetes 部署 ==="
	@echo "  make k8s-apply        部署到 Kubernetes"
	@echo "  make k8s-delete       从 Kubernetes 删除"
	@echo "  make k8s-logs         查看 K8s Pod 日志"
	@echo ""
	@echo "=== 清理 ==="
	@echo "  make clean            删除构建产物"
	@echo ""

# ========== 代码质量 ==========
fmt:
	cd $(FRONTEND_DIR) && npm run lint:fix 2>/dev/null || true
	go fmt ./...

lint:
	go vet ./...
	cd $(FRONTEND_DIR) && npm run lint 2>/dev/null || true

test:
	go test -v ./...

# ========== Swagger ==========
swagger-gen:
	@echo "生成 Swagger 文档..."
	swag init --parseInternal --parseDependency ./docs

# ========== 前端构建 ==========
frontend-build:
	@echo "构建前端..."
	cd $(FRONTEND_DIR) && npm run build

# ========== Go 构建 ==========
build: swagger-gen
	@echo "构建 Go 应用 (本地平台，CGO 启用)..."
	go build -ldflags="$(GO_LDFLAGS)" -o $(APP_NAME) ./cmd/server/main.go

build-linux: swagger-gen
	@echo "交叉编译 Linux amd64（使用 Docker 多阶段构建更稳定）..."
	@echo "推荐使用: make build-docker"
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc go build -ldflags="$(GO_LDFLAGS)" -o $(APP_NAME)-linux-amd64 ./cmd/server/main.go

# ========== Docker 镜像构建 ==========
build-docker: frontend-build swagger-gen
	@echo "构建 Docker 镜像: $(IMAGE_NAME)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		-t $(IMAGE_NAME) \
		-t $(APP_NAME):latest \
		-f $(DOCKER_DIR)/Dockerfile \
		..
	docker tag $(APP_NAME):latest $(IMAGE_NAME)
	@echo "构建完成！"
	@echo "  镜像名: $(IMAGE_NAME)"
	@echo "  推送命令: docker push $(IMAGE_NAME)"

# ========== 全量构建 ==========
build-all: frontend-build build-linux
	@echo "全量构建完成"

# ========== 清理 ==========
clean:
	rm -f $(APP_NAME) $(APP_NAME)-linux-amd64
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules/.vite
	find . -name "*.test" -delete 2>/dev/null || true

# ========== Docker 部署 ==========
docker-up:
	@echo "启动 Docker 服务 (MySQL 版)..."
	cd $(DOCKER_DIR) && docker compose up -d
	@echo "服务已启动:"
	@echo "  Web UI:  http://localhost:8080"
	@echo "  API:     http://localhost:8080/api"
	@echo "  Swagger: http://localhost:8080/swagger"

docker-up-sqlite:
	@echo "启动 Docker 服务 (SQLite 版)..."
	cd $(DOCKER_DIR) && docker compose -f docker-compose.sqlite.yml up -d
	@echo "服务已启动: http://localhost:8080"

docker-down:
	cd $(DOCKER_DIR) && docker compose down
	cd $(DOCKER_DIR) && docker compose -f docker-compose.sqlite.yml down

docker-restart: docker-down docker-up

docker-logs:
	cd $(DOCKER_DIR) && docker compose logs -f

# ========== Kubernetes ==========
k8s-apply:
	@echo "部署到 Kubernetes..."
	kubectl apply -f $(K8S_DIR)/namespace.yaml
	kubectl apply -f $(K8S_DIR)/manifests.yaml
	@echo "部署完成，等待 Pod 就绪..."
	kubectl wait --for=condition=ready pod -l app=datauptwo -n dataregistry --timeout=120s
	@echo "查看状态: kubectl get pods -n dataregistry"

k8s-delete:
	@echo "从 Kubernetes 删除..."
	kubectl delete -f $(K8S_DIR)/manifests.yaml || true
	kubectl delete namespace dataregistry || true

k8s-logs:
	kubectl logs -l app=datauptwo -n dataregistry -f

# ========== 开发模式 ==========
run-dev:
	@echo "启动开发模式..."
	@echo "前端: http://localhost:4000"
	@echo "后端: http://localhost:18421"
	cd $(FRONTEND_DIR) && npm run dev &
	go run cmd/server/main.go
