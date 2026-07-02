FROM --platform=$BUILDPLATFORM golang:1.25 AS amprobe-builder
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /src
COPY go.work ./
COPY common ./common
COPY amprobe ./amprobe
COPY collia ./collia
WORKDIR /src/amprobe
RUN GOARCH=$TARGETARCH go build -ldflags "-X main.Version=${VERSION}" -o /out/amprobe ./cmd/amprobe

FROM --platform=$BUILDPLATFORM golang:1.25 AS collia-builder
ARG VERSION=dev

# collia 要下发到各种架构的 agent 机器，与 amprobe 自身架构无关，
# 因此这里固定同时构建 amd64 和 arm64 两个裸二进制。
WORKDIR /src
COPY go.work ./
COPY common ./common
COPY amprobe ./amprobe
COPY collia ./collia
WORKDIR /src/collia
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o /out/collia-amd64 ./cmd/collia && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=${VERSION}" -o /out/collia-arm64 ./cmd/collia

FROM ubuntu:22.04
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx supervisor ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -M -s /sbin/nologin nginx && \
    mkdir -p /app/logs /app/downloads/collia/amd64 /app/downloads/collia/arm64

COPY --from=amprobe-builder /out/amprobe /app/amprobe
COPY --from=collia-builder /out/collia-amd64 /app/downloads/collia/amd64/collia
COPY --from=collia-builder /out/collia-arm64 /app/downloads/collia/arm64/collia
COPY amprobe/configs /app/configs
COPY amprobe/nginx /etc/nginx
COPY amprobe/supervisor /etc/supervisor/conf.d
COPY deploy/downloads/collia/certs /app/downloads/collia/certs

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
