# ==================== 第一阶段：构建 Builder ====================
FROM golang:1.22-bookworm AS builder

# 替换为阿里云源
RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GOTOOLCHAIN=auto

# 复制整个项目到 /build
COPY . /build/
WORKDIR /build/backend

# 安装编译依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

# 生成 Swagger 文档（假设 main.go 在 cmd/server/main.go）
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --parseInternal --parseDependency -g ./cmd/server/main.go -o ./docs || true

ARG TARGETARCH=amd64
RUN GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=1 \
    go build -ldflags="-s -w -X main.version=v1.0.0 -X main.buildTime=$(date +%Y-%m-%d)" \
    -o duptwo \
    ./cmd/server/main.go

# ==================== 第二阶段：运行镜像 ====================
FROM debian:bookworm-slim AS runner

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list.d/debian.sources

RUN groupadd --gid 1000 appgroup && \
    useradd --uid 1000 --gid appgroup --shell /bin/false --create-home appuser

WORKDIR /app

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
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/backend/duptwo /app/duptwo
COPY --from=builder /build/backend/cmd/server/web /app/web

# 创建目录并设置权限（在切换用户之前）
RUN mkdir -p /app/data /app/logs /app/conf && \
    chown -R 1000:1000  /app

USER appuser

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:18421/health || exit 1

CMD ["./duptwo", "--config", "conf/app.yaml"]
