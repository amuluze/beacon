ARG NODE_IMAGE=public.ecr.aws/docker/library/node:22.22.2-bookworm-slim
ARG RUNTIME_IMAGE=public.ecr.aws/docker/library/alpine:3.22
ARG APK_REPOSITORY=

FROM --platform=$BUILDPLATFORM ${NODE_IMAGE} AS beacon-web-builder
WORKDIR /src/beacon/web
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
COPY beacon/web/package.json beacon/web/pnpm-lock.yaml beacon/web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY beacon/web/ ./
RUN pnpm run build

FROM --platform=$BUILDPLATFORM public.ecr.aws/docker/library/golang:1.25 AS beacon-builder
ARG TARGETOS=linux
ARG TARGETARCH
ARG BEACON_VERSION=dev
ENV GOPROXY=https://goproxy.cn|https://goproxy.io|direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /src
COPY go.work ./
COPY common ./common
COPY beacon ./beacon
COPY collia ./collia
COPY --from=beacon-web-builder /src/beacon/web/dist ./beacon/web/dist
WORKDIR /src/beacon
# 依赖下载单独成层，使 amd64/arm64 Beacon 构建复用同一模块缓存。
# 代理回退到 VCS 后可能留下无法继续 fetch 的半成品缓存；失败重试前清理，
# 让下一轮重新从 GOPROXY 开始，而不是永久卡在 origin。
RUN for attempt in 1 2 3 4 5; do \
        if go mod download; then break; fi; \
        if [ "$attempt" = "5" ]; then exit 1; fi; \
        go clean -modcache; \
        sleep 2; \
    done
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.Version=$BEACON_VERSION" -o /out/beacon ./cmd/beacon

FROM --platform=$BUILDPLATFORM public.ecr.aws/docker/library/golang:1.25 AS collia-builder
ARG COLLIA_VERSION=dev
ENV GOPROXY=https://goproxy.cn|https://goproxy.io|direct
ENV GOSUMDB=sum.golang.google.cn

# collia 要下发到各种架构的 agent 机器，与 beacon 自身架构无关，
# 因此这里固定同时构建 amd64 和 arm64 两个裸二进制。
WORKDIR /src
COPY go.work ./
COPY common ./common
COPY beacon ./beacon
COPY collia ./collia
WORKDIR /src/collia
RUN for attempt in 1 2 3 4 5; do \
        if go mod download; then break; fi; \
        if [ "$attempt" = "5" ]; then exit 1; fi; \
        go clean -modcache; \
        sleep 2; \
    done
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w -X main.version=$COLLIA_VERSION" -o /out/collia-amd64 ./cmd && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -trimpath -ldflags="-s -w -X main.version=$COLLIA_VERSION" -o /out/collia-arm64 ./cmd

FROM ${RUNTIME_IMAGE}
ARG APK_REPOSITORY
ENV BEACON_APP_ENV=production
WORKDIR /app

RUN if [ -n "$APK_REPOSITORY" ]; then \
        sed -i "s|https://dl-cdn.alpinelinux.org/alpine|$APK_REPOSITORY|g" /etc/apk/repositories; \
    fi && \
    apk add --no-cache nginx supervisor ca-certificates && \
    mkdir -p /app/logs /app/downloads/collia/amd64 /app/downloads/collia/arm64 \
        /etc/supervisor/conf.d

COPY --from=beacon-builder /out/beacon /app/beacon
COPY --from=collia-builder /out/collia-amd64 /app/downloads/collia/amd64/collia
COPY --from=collia-builder /out/collia-arm64 /app/downloads/collia/arm64/collia
COPY beacon/configs /app/configs
COPY beacon/nginx /etc/nginx
COPY beacon/supervisor/server.conf beacon/supervisor/nginx.conf /etc/supervisor/conf.d/
COPY beacon/supervisor/supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 80 17000

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
