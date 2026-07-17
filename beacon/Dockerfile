FROM --platform=$BUILDPLATFORM golang:1.25 AS beacon-builder
ARG TARGETARCH
ENV GOPROXY=https://goproxy.cn|https://goproxy.io|direct
ENV GOSUMDB=sum.golang.google.cn

WORKDIR /src
COPY go.work ./
COPY common ./common
COPY beacon ./beacon
COPY collia ./collia
WORKDIR /src/beacon
RUN for attempt in 1 2 3 4 5; do \
        go mod download && break; \
        if [ "$attempt" = "5" ]; then exit 1; fi; \
        sleep 2; \
    done && \
    GOARCH=$TARGETARCH go build -o /out/beacon ./cmd/beacon

FROM --platform=$BUILDPLATFORM golang:1.25 AS collia-builder
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
        go mod download && break; \
        if [ "$attempt" = "5" ]; then exit 1; fi; \
        sleep 2; \
    done && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/collia-amd64 ./cmd/collia && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /out/collia-arm64 ./cmd/collia

FROM public.ecr.aws/ubuntu/ubuntu:22.04
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx supervisor ca-certificates wget && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -M -s /sbin/nologin nginx && \
    mkdir -p /app/logs /app/downloads/collia/amd64 /app/downloads/collia/arm64

COPY --from=beacon-builder /out/beacon /app/beacon
COPY --from=collia-builder /out/collia-amd64 /app/downloads/collia/amd64/collia
COPY --from=collia-builder /out/collia-arm64 /app/downloads/collia/arm64/collia
COPY beacon/configs /app/configs
COPY beacon/nginx /etc/nginx
COPY beacon/supervisor /etc/supervisor/conf.d
COPY deploy/downloads/collia/certs /app/downloads/collia/certs

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
