# ============================================================
# duptwo Dockerfile
# 多阶段构建：Builder + Runner
# 支持 SQLite 和 MySQL/PostgreSQL
# ============================================================

# ==================== 第一阶段：构建 Builder ====================
FROM debian:bookworm-slim AS builder

# 安装构建依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        libc6-dev \
        make \
        curl \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# 复制源码
COPY . .

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    export PATH=$PATH:$(go env GOPATH)/bin && \
    swag init --parseInternal --parseDependency ./docs || true

# 编译：Linux，启用 CGO（SQLite 需要）
ARG TARGETARCH=amd64
RUN GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=1 \
    CC=gcc \
    go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildTime=$(date +%Y-%m-%d)" \
    -o duptwo \
    ./cmd/server/main.go

# ==================== 第二阶段：运行镜像 ====================
FROM debian:bookworm-slim AS runner

# 安全加固：创建非 root 用户
RUN groupadd --gid 1000 appgroup && \
    useradd --uid 1000 --gid appgroup --shell /bin/false --create-home appuser

WORKDIR /app

# 安装运行时依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libsqlite3-0 \
        tini \
        curl \
    && rm -rf /var/lib/apt/lists/*

# 从 Builder 复制二进制
COPY --from=builder /build/duptwo /app/duptwo

# 复制 Web 前端
COPY --from=builder /build/cmd/server/web /app/web

# 创建目录
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:18421/health || exit 1

# 默认启动命令（SQLite 模式）
CMD ["./duptwo", "--config", "conf/app.yaml", "--port", "18421"]
