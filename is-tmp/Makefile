.PHONY: help build build_frontend build_backend build_all clean run test fmt lint

help: ## 显示帮助
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[0;32m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## 构建后端
build_backend:
	@echo "Building backend..."
	cd backend && go build -ldflags "-s -w" -o bin/data-registry cmd/server/main.go

## 构建前端
build_frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

## 完整构建
build_all: build_backend build_frontend
	@echo "Build completed!"

## 开发模式
dev:
	@echo "Starting development server..."
	cd frontend && npm run dev

## 后端运行
run:
	@echo "Running backend..."
	cd backend && go run cmd/server/main.go

## 清理
clean:
	@echo "Cleaning..."
	rm -rf backend/bin
	rm -rf frontend/dist

## 测试
test:
	@echo "Running tests..."
	cd backend && go test -v ./...

## 代码格式化
fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...
	cd frontend && npx prettier --write "src/**/*.{ts,vue,scss}"

## 代码检查
lint:
	@echo "Linting code..."
	cd backend && go vet ./...
	cd frontend && npm run lint

## 交叉编译 - Linux
build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o backend/bin/data-registry-linux-amd64 backend/cmd/server/main.go

build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o backend/bin/data-registry-linux-arm64 backend/cmd/server/main.go

## Docker 构建
docker_build:
	@echo "Building Docker image..."
	docker build -t data-registry:latest .

## Docker 运行
docker_run:
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 --name data-registry data-registry:latest

## Docker 停止
docker_stop:
	docker stop data-registry || true
	docker rm data-registry || true
