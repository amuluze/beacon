FROM node:20.18.3-bullseye-slim

WORKDIR /app

# 创建工作目录
RUN mkdir -p /app && \
    apt update && \
    apt install -y supervisor

# Copy built assets from previous stage
COPY ./web/.output/ /app/.output/
COPY ./web/package.json /app/

COPY ./server/server /app
COPY ./server/configs /app/configs
COPY ./supervisor /etc/supervisor/conf.d


EXPOSE 3000

# ENV NODE_ENV=production \
#     NUXT_PUBLIC_BASE_URL=http://127.0.0.1:8000

# ENTRYPOINT [ "node", "/app/.output/server/index.mjs" ]

# 这里有点坑，不加 -n 服务启动不了
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
