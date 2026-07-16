############################
# Stage 1: 构建前端（Nuxt）
############################
FROM node:22.22.2-bookworm-slim AS web-builder
WORKDIR /build/web
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
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
FROM node:22.22.2-bookworm-slim AS runtime
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends supervisor \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 app \
    && useradd --system --uid 10001 --gid 10001 --no-create-home --home-dir /app app

# 前端产物（nitro 自包含）
COPY --from=web-builder /build/web/.output/ /app/.output/
# Go server 二进制
COPY --from=server-builder /out/server /app/server
# 配置与 supervisor 编排
COPY server/configs/ /app/configs/
COPY supervisor/ /etc/supervisor/conf.d
# 发布物静态资源（/release/* 对外暴露 manager.sh 等）
COPY release/ /app/release/

# 运行期可写目录（应用日志、sqlite 数据），归属降权用户 app
RUN mkdir -p /app/logs /app/data \
    && chown -R app:app /app

EXPOSE 3000

# supervisord 作为 PID1 保留 root 以管理 socket 与日志句柄；
# 实际业务进程（node / server）通过 program 的 user=app 降权运行（见 supervisor/*.conf）。
# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
