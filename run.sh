#!/bin/bash
#####################################################################
# Author     : Amuluze
# Date       : 2024/05/24 10:47:05
# Description: 一键部署脚本
#              make -p /data
#              unzip amprobe.installer.zip -d /data
#              cd /data/amvector
#              ./amvector install --conf=/data/amvector/configs/config.toml
#              ./amvector start
#               cd /data
#              docker load -i amprobe.tar.gz
#              docker run -itd --name amprobe -p 80:80 \
#               -v /data/amprobe/configs:/app/configs \
#               -v /data/amprobe/nginx/nginx.conf:/etc/nginx/nginx.conf \
#               -v /data/amvector/amvector.socket:/app/amvector.socket \
#               -v /var/run/docker.sock:/var/run/docker.sock \
#               amuluze/amprobe:v1.3.4
#####################################################################
