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
