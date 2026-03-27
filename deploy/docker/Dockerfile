# ==================== 第一阶段：构建 Builder（debian-slim 支持 CGO） ====================
FROM debian:bookworm-slim AS builder

# 安装构建依赖（gcc 用于 CGO SQLite）
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        libc6-dev \
        make \
        curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# 复制源码
COPY . .

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --parseInternal --parseDependency ./docs

# 编译：Linux amd64，启用 CGO（SQLite 需要）
ARG TARGETARCH=amd64
RUN GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=1 \
    CC=gcc \
    go build -ldflags="-s -w" \
    -o duptwo \ 
    ./cmd/server/main.go

# ==================== 第二阶段：运行镜像（debian-slim） ====================
FROM debian:bookworm-slim AS runner

# 安全加固：创建非 root 用户
RUN groupadd --gid 1000 appgroup && \
    useradd --uid 1000 --gid appgroup --shell /bin/false --create-home appuser

WORKDIR /app

# 安装运行时最小依赖（SQLite + tini）
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libsqlite3-0 \
        tini \
    && rm -rf /var/lib/apt/lists/*

# 从 Builder 复制二进制（动态链接 debian libc）
COPY --from=builder /build/duptwo  /app/duptwo

# 复制 Web 前端
COPY --from=builder /build/cmd/server/web /app/web

# 复制配置文件模板
COPY --from=builder /build/conf /app/conf

# 创建数据目录
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appgroup /app

# 使用 tini 作为 init（PID 1）
ENTRYPOINT ["/usr/bin/tini", "--"]

# 切换到非 root 用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 默认启动命令
CMD ["./duptwo", "--config", "conf/app.yaml"]
