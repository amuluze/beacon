#!/bin/sh
# 为镜像内置的 collia 发布二进制生成 Ed25519 侧车签名。
# 用法：sign_release_packages.sh <out-dir> [sign-binary]
# 私钥路径由 COLLIA_SIGNING_KEY_FILE 指定（默认 /run/secrets/collia_signing_key，
# 由 BuildKit secret 挂载）；私钥缺失或为空时跳过——未签名构建在 Agent 未固定
# 发布公钥时仍按 SHA-256 兼容放行，已固定公钥的 Agent 会拒绝该升级包。
set -eu
OUT_DIR="${1:?用法: $0 <out-dir> [sign-binary]}"
SIGN_BIN="${2:-$OUT_DIR/sign}"
KEY_FILE="${COLLIA_SIGNING_KEY_FILE:-/run/secrets/collia_signing_key}"

if [ ! -s "$KEY_FILE" ]; then
    echo "[sign-collia] 未提供签名私钥（$KEY_FILE），跳过侧车签名（兼容模式）"
    exit 0
fi
if [ ! -x "$SIGN_BIN" ]; then
    echo "[sign-collia] 签名工具不可执行：$SIGN_BIN" >&2
    exit 1
fi

for arch in amd64 arm64; do
    target="$OUT_DIR/collia-$arch"
    if [ ! -f "$target" ]; then
        echo "[sign-collia] 缺少发布二进制：$target" >&2
        exit 1
    fi
    "$SIGN_BIN" -key "$KEY_FILE" "$target"
done
echo "[sign-collia] 已生成 amd64/arm64 侧车签名"
