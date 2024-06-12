#!/bin/bash
#################################
# Author     : Amu
# Date       : 2024/5/16 15:16
# Description: 生成自签名证书
#################################

ca_base_dir="$1"
output_dir="$2"

ca_key="$ca_base_dir/ca.key"
ca_pem="$ca_base_dir/ca.pem"

cert_key="$output_dir/tls.key"
cert_csr="$output_dir/tls.csr"
cert_crt="$output_dir/tls.crt"

if [[ -f $ca_key && -f $ca_pem ]]; then
    echo "ca.key and ca.pem exist"
    openssl genrsa -out "$cert_key" 4096 || exit 1
    openssl req -new -key "$cert_key" -out "$cert_csr" -config ./san.conf -sha256 -subj "/C=CN/ST=Beijing/L=Beijing/O=Amuluze/OU=Amuluze/CN=amvector" || exit 1
    openssl x509 -req -days 365 -in "$cert_csr" -CA "$ca_pem" -CAkey "$ca_key" -set_serial $RANDOM -out "$cert_crt" -extfile san.conf -extensions v3_req || exit 1
    rm "$cert_csr"
    echo "tls.key, tls.csr and tls.crt created"
else
    echo "ca.key and ca.pem not exist"
    exit 2
fi

exit 0
