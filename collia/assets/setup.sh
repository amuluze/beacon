#!/bin/bash
#################################
# Author     : Amu
# Date       : 2024/7/21 11:48
# Description:
#################################

set -ex


# parse parameters
params="$(getopt -o pi --name "$0" -- "$@")"
eval set -- "$params"
while true; do
    case "$1" in
        -p)
            shift
            ;;
        -i)
            shift
            ;;
        --)
            shift
            break
            ;;
        *)
            echo "Unknown Option: $1" >&2
            exit 1
            ;;
    esac
done

# directories
mkdir -p /etc/collia
chmod 755 /etc/collia

# binary
install -m 755 -b collia /usr/sbin/collia

# config
if [ ! -f /etc/collia/config.yml ]; then
    cp config.yml /etc/collia/config.yml
    chmod 644 /etc/collia/config.yml
else
    echo "collia already exists"
fi
