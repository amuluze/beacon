#!/bin/bash
#################################
# Author     : Amu
# Date       : 2024/7/21 11:48
# Description:
#################################

set -ex

# directories
mkdir -p /etc/amvector
chmod 755 /etc/amvector

# binary
install -m 755 -b amvector /usr/sbin/amvector

# config
if [ ! -f /etc/amvector/config.yml ]; then
    cp config.yml /etc/amvector/config.yml
    chmod 644 /etc/amvector/config.yml
else
    echo "amvector already exists"
fi
