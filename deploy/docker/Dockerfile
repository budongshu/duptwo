# ==================== 第一阶段：构建 Builder ====================
FROM golang:1.24-bookworm AS builder

WORKDIR /build

# 设置 Go 环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1
ENV GO111MODULE=on
# 设置 Go 工具链为自动
ENV GOTOOLCHAIN=auto

COPY backend/ . 


# 复制 go.mod 和 go.sum 文件（利用 Docker 缓存，避免重复下载依赖）
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

# 复制项目源码
COPY . .

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --parseInternal --parseDependency -g ./cmd/server/main.go -o ./docs || true

# 构建参数
ARG TARGETARCH=amd64
ARG VERSION=v1.0.0
ARG BUILD_TIME

# 编译二进制文件
RUN GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=1 \
    go build -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=$(date +%Y-%m-%d_%H:%M:%S)" \
    -o duptwo \
    ./cmd/server/main.go

# ==================== 第二阶段：运行镜像 ====================
FROM debian:bookworm-slim AS runner

# 替换为阿里云源（加速软件包下载）
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources

# 安装运行时依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libsqlite3-0 \
        tini \
        curl \
        net-tools \
        iputils-ping \
        iproute2 \
        telnet \
        procps \
        vim-tiny \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# 创建非 root 用户
RUN groupadd --gid 1000 appgroup && \
    useradd --uid 1000 --gid appgroup --shell /bin/bash --create-home appuser

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /build/duptwo /app/duptwo

# 复制 web 静态文件（如果存在）
COPY --from=builder /build/cmd/server/web /app/web  

# 复制 Swagger 文档（如果存在）
COPY --from=builder /build/docs /app/docs 

# 创建必要的目录
RUN mkdir -p /app/data /app/logs /app/conf /app/tmp && \
    chown -R 1000:1000 /app && \
    chmod +x /app/duptwo

# 切换为非 root 用户
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:18421/health || exit 1

# 使用 tini 作为 init 进程，正确处理信号
ENTRYPOINT ["/usr/bin/tini", "--"]

# 启动应用
CMD ["./duptwo", "--config", "conf/app.yaml"]
