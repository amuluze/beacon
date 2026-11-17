#!/bin/sh
# Beacon 容器启动入口：校验关键密钥后交棒给 CMD（默认 supervisord）。
set -eu

# 显式覆盖为 shell/调试命令时不做生产密钥校验；仅默认 supervisord 启动路径校验。
if [ "$#" -gt 0 ] && [ "$1" != "/usr/bin/supervisord" ]; then
    exec "$@"
fi

fail() {
    echo "ERROR: $*" >&2
    exit 1
}

# BusyBox ash 不支持 ${#VAR:-}，先归一化默认值再做长度校验。
BEACON_AUTH_SIGNING_KEY="${BEACON_AUTH_SIGNING_KEY:-}"
BEACON_AGENT_INSTALL_TOKEN="${BEACON_AGENT_INSTALL_TOKEN:-}"
BEACON_CONTROL_JOIN_TOKEN="${BEACON_CONTROL_JOIN_TOKEN:-}"
BEACON_DATA_ENCRYPTION_KEY="${BEACON_DATA_ENCRYPTION_KEY:-}"

case "$BEACON_AUTH_SIGNING_KEY" in
    beacon|change-me|"") fail "BEACON_AUTH_SIGNING_KEY must be a strong value" ;;
esac
case "$BEACON_AGENT_INSTALL_TOKEN" in
    beacon|change-me|"") fail "BEACON_AGENT_INSTALL_TOKEN must be a strong value" ;;
esac

[ "${#BEACON_AUTH_SIGNING_KEY}" -ge 48 ] || fail "BEACON_AUTH_SIGNING_KEY must have at least 48 characters"
[ "${#BEACON_AGENT_INSTALL_TOKEN}" -ge 48 ] || fail "BEACON_AGENT_INSTALL_TOKEN must have at least 48 characters"
[ "${#BEACON_CONTROL_JOIN_TOKEN}" -ge 48 ] || fail "BEACON_CONTROL_JOIN_TOKEN must have at least 48 characters"

case "$BEACON_DATA_ENCRYPTION_KEY" in
    *[!0-9A-Fa-f]*|"") fail "BEACON_DATA_ENCRYPTION_KEY must be a 64-character hexadecimal key" ;;
esac
[ "${#BEACON_DATA_ENCRYPTION_KEY}" -eq 64 ] || fail "BEACON_DATA_ENCRYPTION_KEY must be a 64-character hexadecimal key"

if [ "$#" -eq 0 ]; then
    set -- /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
fi
exec "$@"
