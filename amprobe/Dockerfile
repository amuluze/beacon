FROM golang:1.21 AS builder

WORKDIR /src
COPY go.work ./
COPY common ./common
COPY amprobe ./amprobe
COPY collia ./collia
WORKDIR /src/amprobe
RUN go build -o /out/amprobe ./cmd/amprobe

FROM ubuntu:22.04
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx supervisor ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -M -s /sbin/nologin nginx && \
    mkdir -p /app/logs /app/downloads

COPY --from=builder /out/amprobe /app/amprobe
COPY amprobe/configs /app/configs
COPY amprobe/nginx /etc/nginx
COPY amprobe/supervisor /etc/supervisor/conf.d

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
