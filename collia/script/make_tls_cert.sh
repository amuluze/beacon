#!/bin/bash
#################################
# Author     : Amu
# Date       : 2024/5/16 15:16
# Description: 生成自签名证书
#################################

ca_base_dir="$1"
output_dir="$2"
common_name="${3:-amprobe/collia}"
dns_names="${4:-$common_name}"
ip_names="${5:-127.0.0.1}"

ca_key="$ca_base_dir/ca.key"
ca_pem="$ca_base_dir/ca.pem"

cert_key="$output_dir/tls.key"
cert_csr="$output_dir/tls.csr"
cert_crt="$output_dir/tls.crt"
san_conf="$output_dir/san.conf"

if [[ -f $ca_key && -f $ca_pem ]]; then
    mkdir -p "$output_dir"
    echo "ca.key and ca.pem exist"
    {
        echo "[ req ]"
        echo "default_bits       = 4096"
        echo "distinguished_name = req_distinguished_name"
        echo "req_extensions     = v3_req"
        echo
        echo "[ req_distinguished_name ]"
        echo "countryName         = CN"
        echo "stateOrProvinceName = Beijing"
        echo "localityName        = Beijing"
        echo "organizationName    = Amuluze"
        echo "commonName          = $common_name"
        echo
        echo "[ v3_req ]"
        echo "subjectAltName = @alt_names"
        echo
        echo "[alt_names]"
        IFS=',' read -ra dns_array <<< "$dns_names"
        dns_index=1
        for dns in "${dns_array[@]}"; do
            [[ -z "$dns" ]] && continue
            echo "DNS.$dns_index = $dns"
            dns_index=$((dns_index + 1))
        done
        IFS=',' read -ra ip_array <<< "$ip_names"
        ip_index=1
        for ip in "${ip_array[@]}"; do
            [[ -z "$ip" ]] && continue
            echo "IP.$ip_index = $ip"
            ip_index=$((ip_index + 1))
        done
    } > "$san_conf"
    openssl genrsa -out "$cert_key" 4096 || exit 1
    openssl req -new -key "$cert_key" -out "$cert_csr" -config "$san_conf" -sha256 -subj "/C=CN/ST=Beijing/L=Beijing/O=Amuluze/OU=Amuluze/CN=$common_name" || exit 1
    openssl x509 -req -days 365 -in "$cert_csr" -CA "$ca_pem" -CAkey "$ca_key" -set_serial $RANDOM -out "$cert_crt" -extfile "$san_conf" -extensions v3_req || exit 1
    cp "$ca_pem" "$output_dir/ca.pem"
    rm "$cert_csr"
    rm "$san_conf"
    echo "ca.pem, tls.key and tls.crt created"
else
    echo "ca.key and ca.pem not exist"
    exit 2
fi

exit 0
