#!/bin/bash
#####################################################################
# Author     : Amu
# Date       : 2024/07/21 16:04:07
# Description:
#####################################################################

set -ex


# 检查是否已安装 collia，若已安装则跳过
[ -r "/etc/collia/config.yml" ] && {
    echo "Collia already installed!"
    exit 0
}

# 安装 collia
./collia.install  || {
    echo "Installation failed!"
    exit 1
}

mkdir -p /etc/collia/versions
cp -v VERSION.TXT /etc/collia/versions
