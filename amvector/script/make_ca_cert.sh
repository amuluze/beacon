#!/bin/bash
#################################
# Author     : Amu
# Date       : 2024/5/16 15:14
# Description: 生产 CA 证书
#################################

ca_base_dir="$1"
ca_key="$ca_base_dir/ca.key"
ca_pem="$ca_base_dir/ca.pem"

# 生成CA私钥
openssl genrsa -out "$ca_key" 4096

# 生成CA证书
openssl req -new -x509 -days 365 -key "$ca_key" -out "$ca_pem" -nodes -subj "/C=CN/ST=Beijing/L=Beijing/O=Amuluze/OU=Amuluze/CN=root"
rm "$ca_key"

# 生成CA私钥和证书
# openssl req -x509 -newkey rsa:4096 -keyout $ca_key -out $ca_crt -days 365 -nodes -subj "/C=CN/ST=Beijing/L=Beijing/O=Amuluze/OU=Amuluze/CN=root"
