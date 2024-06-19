#!/bin/bash
#####################################################################
# Author     : wangjialong
# Date       : 2024/03/15 16:04:56
# Description:
#####################################################################

set -e

# 将 dist 文件中
sshpass -f 'dev' scp -r ./dist/* root@101.42.246.113:/data/web/dist/
# sshpass -f 'dev' ssh root@101.42.246.113 << remotessh
# docker restart amprobe
# remotessh
