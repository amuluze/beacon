#!/bin/bash
#####################################################################
# Author     : Amu
# Date       : 2024/07/21 16:04:07
# Description:
#####################################################################

set -ex


# 检查是否已安装 amvector，若已安装则跳过
[ -r "/etc/amvector/config.yml" ] && {
    echo "Amvector already installed!"
    exit 0
}

# 安装 amvector
./amvector.install  || {
    echo "Installation failed!"
    exit 1
}

mkdir -p /etc/amvector/versions
cp -v VERSION.TXT /etc/amvector/versions
