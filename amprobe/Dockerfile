FROM ubuntu:22.04
WORKDIR /app

COPY ./amprobe /app/

RUN apt update && \
    apt install -y nginx supervisor && \
    useradd -M -s /sbin/nologin nginx && \
    touch /app/probe.db && \
    chmod +x /app/amprobe

COPY ./configs /app/configs
COPY ./dist /usr/share/nginx/html/app
COPY ./nginx /etc/nginx
COPY ./supervisor /etc/supervisor/conf.d

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
