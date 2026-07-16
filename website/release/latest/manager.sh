#!/bin/sh
# Beacon 一键安装脚本
#
# 用法：
#   curl -fsSL https://help.beacon.amuluze.com/release/latest/manager.sh | sh
#   curl -fsSL https://help.beacon.amuluze.com/release/latest/manager.sh | sh -s -- BEACON_HTTP_PORT=1443
# 或下载后执行：
#   sh manager.sh [KEY=VAL ...]
#
# 所有部署步骤（依赖检查 / 目录生成 / 拉 compose / 生成 .env / 拉镜像 / 启动容器）
# 均在本脚本内完成。通过环境变量或 KEY=VAL 参数覆盖默认值，支持：
#   BEACON_BASE_URL / INSTALL_DIR / BEACON_IMAGE
#   BEACON_HTTP_PORT / BEACON_CONTROL_PORT / BEACON_PUBLIC_BASE_URL
#   BEACON_AGENT_INSTALL_TOKEN / BEACON_AUTH_SIGNING_KEY / BEACON_CONTROL_JOIN_TOKEN
set -eu

DEFAULT_BASE_URL="https://help.beacon.amuluze.com"
DEFAULT_INSTALL_DIR="/data/beacon"
DEFAULT_IMAGE="registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:latest"
DEFAULT_HTTP_PORT="1443"
DEFAULT_CONTROL_PORT="17000"

log()  { printf '[beacon] %s\n' "$*"; }
warn() { printf '[beacon] %s\n' "$*" >&2; }
die()  { printf '[beacon] error: %s\n' "$*" >&2; exit 1; }

# 命令行 KEY=VAL 参数覆盖环境变量
for arg in "$@"; do
    case "$arg" in
        *=*) export "$arg" ;;
        *)   die "无法识别的参数：$arg（格式应为 KEY=VAL）" ;;
    esac
done

# root 校验
if [ "$(id -u)" -ne 0 ]; then
    die "本脚本需要 root 权限运行，请使用 root 或 sudo。"
fi

# 依赖检查
need_cmd() {
    command -v "$1" >/dev/null 2>&1 || die "缺少必要命令：$1"
}
need_cmd curl
need_cmd docker

# docker compose：plugin 优先，回退独立二进制
if docker compose version >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
    DOCKER_COMPOSE="docker-compose"
else
    die "未检测到 docker compose 插件或 docker-compose，请先安装。"
fi

# 随机密钥：openssl 优先，回退 cksum
random_secret() {
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -hex 24
    else
        date +%s | cksum | awk '{ print $1 }'
    fi
}

# 交互读取：TTY 时提示，否则用默认值
prompt() {
    label="$1"; default_value="$2"
    if [ -t 0 ]; then
        printf '%s [%s]: ' "$label" "$default_value" >&2
        read -r value || value=""
        if [ -n "$value" ]; then
            printf '%s' "$value"
            return
        fi
    fi
    printf '%s' "$default_value"
}

# 官网基址：默认官网域名，可用 BEACON_BASE_URL 覆盖
BASE_URL="${BEACON_BASE_URL:-$DEFAULT_BASE_URL}"
BASE_URL="${BASE_URL%/}"

# 参数收集（环境变量优先 → 交互 → 默认）
INSTALL_DIR="${INSTALL_DIR:-$(prompt "安装目录" "$DEFAULT_INSTALL_DIR")}"
BEACON_IMAGE="${BEACON_IMAGE:-$(prompt "Beacon 镜像" "$DEFAULT_IMAGE")}"
BEACON_HTTP_PORT="${BEACON_HTTP_PORT:-$(prompt "Web 控制台宿主端口" "$DEFAULT_HTTP_PORT")}"
BEACON_CONTROL_PORT="${BEACON_CONTROL_PORT:-$(prompt "Agent 控制端口" "$DEFAULT_CONTROL_PORT")}"
BEACON_PUBLIC_BASE_URL="${BEACON_PUBLIC_BASE_URL:-$(prompt "对外访问地址" "http://127.0.0.1:$BEACON_HTTP_PORT")}"
BEACON_AGENT_INSTALL_TOKEN="${BEACON_AGENT_INSTALL_TOKEN:-$(random_secret)}"
BEACON_AUTH_SIGNING_KEY="${BEACON_AUTH_SIGNING_KEY:-$(random_secret)}"
BEACON_CONTROL_JOIN_TOKEN="${BEACON_CONTROL_JOIN_TOKEN:-$(random_secret)}"

log "安装目录：$INSTALL_DIR"
log "镜像：$BEACON_IMAGE"
log "Web 控制台：$BEACON_PUBLIC_BASE_URL"

mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

log "拉取 compose.yaml ..."
curl -fsSL "$BASE_URL/release/latest/compose.yaml" -o compose.yaml

log "生成 .env ..."
cat > .env <<EOF
BEACON_IMAGE=$BEACON_IMAGE
BEACON_CONTAINER_NAME=beacon
BEACON_HTTP_PORT=$BEACON_HTTP_PORT
BEACON_CONTROL_PORT=$BEACON_CONTROL_PORT
BEACON_DATA_DIR=./data
BEACON_LOG_DIR=./logs
BEACON_DB_NAME=/app/data/beacon
BEACON_PUBLIC_BASE_URL=$BEACON_PUBLIC_BASE_URL
BEACON_AGENT_INSTALL_TOKEN=$BEACON_AGENT_INSTALL_TOKEN
BEACON_AUTH_SIGNING_KEY=$BEACON_AUTH_SIGNING_KEY
BEACON_CONTROL_JOIN_TOKEN=$BEACON_CONTROL_JOIN_TOKEN
EOF

mkdir -p data logs

# TTY 下询问是否在启动前编辑 .env
if [ -t 0 ]; then
    printf '是否在启动前编辑 .env？[y/N]: ' >&2
    read -r edit_env || edit_env=""
    case "$edit_env" in
        y|Y|yes|YES)
            "${EDITOR:-vi}" .env
            ;;
    esac
fi

log "拉取镜像（可能需要一些时间）..."
$DOCKER_COMPOSE pull

log "启动 beacon 容器 ..."
$DOCKER_COMPOSE up -d

cat <<EOF

[beacon] 部署完成。
  安装目录        : $INSTALL_DIR
  Web 控制台      : $BEACON_PUBLIC_BASE_URL
  Agent 控制端口  : $BEACON_CONTROL_PORT
  初始账号        : admin / admin123
  ⚠️ 安全警告     : 默认密码已对外公开、极不安全，请立即登录后在「设置」中修改，否则实例存在被未授权访问的风险
  密钥存放位置    : $INSTALL_DIR/.env

常用运维命令：
  cd $INSTALL_DIR
  $DOCKER_COMPOSE ps          # 查看状态
  $DOCKER_COMPOSE logs -f     # 查看日志
  $DOCKER_COMPOSE down        # 停止服务
  $DOCKER_COMPOSE up -d       # 启动 / 更新
EOF
