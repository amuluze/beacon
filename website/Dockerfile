############################
# Stage 1: 构建前端（Nuxt）
############################
FROM node:20.18.3-bullseye-slim AS web-builder
WORKDIR /build/web
RUN corepack enable
# 先装依赖以利用 Docker 层缓存
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
# 再拷源码触发构建
COPY web/ ./
RUN pnpm run build

############################
# Stage 2: 构建 Go server
############################
FROM golang:1.25-alpine AS server-builder
ARG TARGETOS=linux
ARG TARGETARCH=amd64
WORKDIR /build/server
# 先下载依赖以利用层缓存
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
# glebarez/sqlite 为纯 Go，可关闭 CGO 产出静态二进制
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

############################
# Stage 3: 运行时（Node + Go + supervisor）
############################
FROM node:20.18.3-bullseye-slim AS runtime
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends supervisor \
    && rm -rf /var/lib/apt/lists/*

# 前端产物（nitro 自包含）
COPY --from=web-builder /build/web/.output/ /app/.output/
# Go server 二进制
COPY --from=server-builder /out/server /app/server
# 配置与 supervisor 编排
COPY server/configs/ /app/configs/
COPY supervisor/ /etc/supervisor/conf.d

EXPOSE 3000

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
