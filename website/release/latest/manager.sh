#!/bin/sh
# Beacon 一键安装脚本
#
# 用法：
#   bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)"
#   bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)" -- install BEACON_HTTP_PORT=1443
#
# 或下载后执行：
#   sh manager.sh [install|upgrade|uninstall|restart] [KEY=VAL ...]
#
# 支持的操作：install / upgrade / uninstall / restart
# 通过环境变量或 KEY=VAL 参数覆盖默认值，支持：
#   BEACON_BASE_URL / INSTALL_DIR / BEACON_IMAGE
#   BEACON_VERSION / BEACON_ARCH / BEACON_HTTP_PORT / BEACON_CONTROL_PORT / BEACON_PUBLIC_BASE_URL
#   BEACON_AGENT_INSTALL_TOKEN / BEACON_AUTH_SIGNING_KEY / BEACON_CONTROL_JOIN_TOKEN
# BEACON_ARCH 默认由主机架构推断（x86_64→amd64，aarch64→arm64），仅在需要拉取跨架构镜像时手动覆盖。
set -eu
umask 077

DEFAULT_BASE_URL="https://help.beacon.amuluze.com"
DEFAULT_INSTALL_DIR="/data/beacon"
DEFAULT_VERSION="v3.0.4"
DEFAULT_HTTP_PORT="1443"
DEFAULT_CONTROL_PORT="17000"
MIN_DOCKER_VERSION=20
MIN_COMPOSE_VERSION=2
MIN_MEMORY_BYTES=1073741824
MIN_DISK_BYTES=5368709120
HEALTH_TIMEOUT=120

# 颜色定义
if [ -t 2 ]; then
    C_RESET='\033[0m'
    C_RED='\033[31m'
    C_YELLOW='\033[33m'
else
    C_RESET=''
    C_RED=''
    C_YELLOW=''
fi

log()  { printf '[beacon] %s\n' "$*"; }
warn() { printf '%b[beacon]%b %s\n' "$C_YELLOW" "$C_RESET" "$*" >&2; }
err()  { printf '%b[beacon]%b %s\n' "$C_RED" "$C_RESET" "$*" >&2; }
die()  { printf '%b[beacon] error:%b %s\n' "$C_RED" "$C_RESET" "$*" >&2; exit 1; }

# 全局状态
ACTION=""
INSTALL_DIR=""
DOCKER_COMPOSE=""
BACKUP_ENV=""
BACKUP_COMPOSE=""

# 清理临时备份
cleanup() {
    if [ -n "$BACKUP_ENV" ] && [ -f "$BACKUP_ENV" ]; then
        rm -f "$BACKUP_ENV"
    fi
    if [ -n "$BACKUP_COMPOSE" ] && [ -f "$BACKUP_COMPOSE" ]; then
        rm -f "$BACKUP_COMPOSE"
    fi
}
trap cleanup INT TERM EXIT

# 命令依赖检查
need_cmd() {
    command -v "$1" >/dev/null 2>&1 || die "缺少必要命令：$1"
}

# root 校验
check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        die "本脚本需要 root 权限运行，请使用 root 或 sudo。"
    fi
}

# 系统校验：仅 Linux
check_linux() {
    sysname="$(uname -s 2>/dev/null || true)"
    if [ "$sysname" != "Linux" ]; then
        die "Beacon 当前仅支持 Linux，检测到系统：${sysname:-未知}"
    fi
}

# 架构校验：归一化为 amd64/arm64
check_arch() {
    machine="$(uname -m 2>/dev/null || true)"
    case "$machine" in
        x86_64|amd64)
            TARGET_ARCH="amd64"
            ;;
        aarch64|arm64)
            TARGET_ARCH="arm64"
            ;;
        *) die "Beacon 暂时不支持 $machine 架构" ;;
    esac
    # 命令行 / 环境变量优先，否则跟随主机架构
    BEACON_ARCH="${BEACON_ARCH:-$TARGET_ARCH}"
    log "检测架构：$machine → $BEACON_ARCH"
}

# 内存校验（可用内存 >= 1GB）
check_memory() {
    if [ ! -r /proc/meminfo ]; then
        warn "无法读取 /proc/meminfo，跳过内存检查"
        return 0
    fi
    available_kb="$(awk '/^MemAvailable:/ {print $2; exit}' /proc/meminfo)"
    if [ -z "$available_kb" ]; then
        warn "无法获取可用内存，跳过检查"
        return 0
    fi
    available_bytes=$((available_kb * 1024))
    if [ "$available_bytes" -lt "$MIN_MEMORY_BYTES" ]; then
        die "可用内存不足 1GB（当前约 $((available_bytes / 1024 / 1024)) MB），请先扩容"
    fi
}

# 磁盘校验
check_disk() {
    target="$1"
    while [ ! -e "$target" ] && [ "$target" != "/" ]; do
        target="$(dirname "$target")"
    done
    need_cmd df
    free_kb="$(df -k "$target" | awk 'NR==2 {print $4}')"
    if [ -z "$free_kb" ]; then
        warn "无法查询 $target 磁盘容量"
        return 0
    fi
    free_bytes=$((free_kb * 1024))
    if [ "$free_bytes" -lt "$MIN_DISK_BYTES" ]; then
        die "$1 磁盘容量不足，安装 Beacon 至少需要 5GB 可用空间"
    fi
    log "$1 可用空间：$((free_bytes / 1024 / 1024 / 1024)) GB"
}

# 端口占用检查
check_port() {
    port="$1"
    has_ss=false
    has_netstat=false
    command -v ss >/dev/null 2>&1 && has_ss=true
    command -v netstat >/dev/null 2>&1 && has_netstat=true

    if [ "$has_ss" = false ] && [ "$has_netstat" = false ]; then
        warn "未安装 ss/netstat，跳过端口 $port 占用检查"
        return 0
    fi

    if [ "$has_ss" = true ]; then
        if ss -Hln "sport = :$port" 2>/dev/null | grep -q "."; then
            die "端口 $port 已被占用，请更换端口"
        fi
    elif [ "$has_netstat" = true ]; then
        if netstat -tlnp 2>/dev/null | grep -q ":$port "; then
            die "端口 $port 已被占用，请更换端口"
        fi
    fi
}

# Docker 版本检查
check_docker() {
    need_cmd docker
    version_output="$(docker --version 2>/dev/null || true)"
    if [ -z "$version_output" ]; then
        die "无法获取 Docker 版本，请确认 Docker 已正确安装"
    fi
    major_version="$(printf '%s' "$version_output" | sed -n 's/^Docker version \([0-9]\+\).*/\1/p')"
    if [ -z "$major_version" ]; then
        die "无法解析 Docker 版本：$version_output"
    fi
    if [ "$major_version" -lt "$MIN_DOCKER_VERSION" ]; then
        die "Docker 版本过低（需要 >= $MIN_DOCKER_VERSION，当前 $major_version），请先升级"
    fi
    log "Docker 版本：$version_output"
}

# Docker Compose 检查：plugin 优先，回退独立二进制
check_compose() {
    if docker compose version >/dev/null 2>&1; then
        version_output="$(docker compose version 2>/dev/null || true)"
        major_version="$(printf '%s' "$version_output" | sed -n 's/.*v\?\([0-9]\+\).*/\1/p')"
        if [ -n "$major_version" ] && [ "$major_version" -ge "$MIN_COMPOSE_VERSION" ]; then
            DOCKER_COMPOSE="docker compose"
            log "Docker Compose：$version_output"
            return 0
        fi
    fi
    if command -v docker-compose >/dev/null 2>&1; then
        version_output="$(docker-compose version 2>/dev/null || true)"
        major_version="$(printf '%s' "$version_output" | sed -n 's/.*v\?\([0-9]\+\).*/\1/p')"
        if [ -n "$major_version" ] && [ "$major_version" -ge "$MIN_COMPOSE_VERSION" ]; then
            DOCKER_COMPOSE="docker-compose"
            log "Docker Compose：$version_output"
            return 0
        fi
    fi
    die "未检测到 Docker Compose 插件或 docker-compose（需要 >= v$MIN_COMPOSE_VERSION），请先安装"
}

# 综合环境预检查
precheck() {
    check_root
    check_linux
    check_arch
    check_docker
    check_compose
    check_memory
}

# 随机密钥：仅接受密码学安全随机源
random_secret() {
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -hex 32
        return
    fi
    if [ -r /dev/urandom ] && command -v od >/dev/null 2>&1 && command -v tr >/dev/null 2>&1; then
        od -An -N32 -tx1 /dev/urandom | tr -d ' \n'
        return
    fi
    die "无法获得密码学安全随机源（需要 openssl 或 /dev/urandom + od）"
}

validate_secret() {
    name="$1"; value="$2"
    if [ "${#value}" -lt 48 ]; then
        die "$name 长度不足，请提供至少 48 个字符的随机值"
    fi
}

verify_sha256() {
    expected="$1"; file="$2"
    if command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$file" | awk '{ print $1 }')"
    elif command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$file" | awk '{ print $1 }')"
    else
        die "缺少 sha256sum 或 shasum，无法校验发布物"
    fi
    [ "$actual" = "$expected" ] || die "$file SHA-256 校验失败"
}

# 带重试的下载
download() {
    url="$1"; output="$2"
    need_cmd curl
    if ! curl -fsSL --retry 3 --retry-delay 2 --connect-timeout 10 --max-time 120 "$url" -o "$output"; then
        die "下载失败：$url"
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

# TTY 确认
confirm() {
    question="$1"; default="$2"
    if [ -t 0 ]; then
        printf '%s [%s]: ' "$question" "$default" >&2
        read -r answer || answer=""
        case "$answer" in
            y|Y|yes|YES) return 0 ;;
            n|N|no|NO) return 1 ;;
            *) [ "$default" = "y" ] && return 0 || return 1 ;;
        esac
    fi
    [ "$default" = "y" ]
}

# 从 KEY=VAL 文件读取配置到变量
read_env() {
    env_file="$1"
    if [ ! -f "$env_file" ]; then
        return 0
    fi
    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in
            ''|\#*) continue ;;
        esac
        key="${line%%=*}"
        value="${line#*=}"
        case "$key" in
            BEACON_IMAGE)            BEACON_IMAGE="$value" ;;
            BEACON_ARCH)             BEACON_ARCH="$value" ;;
            BEACON_VERSION)          BEACON_VERSION="$value" ;;
            BEACON_HTTP_PORT)        BEACON_HTTP_PORT="$value" ;;
            BEACON_CONTROL_PORT)     BEACON_CONTROL_PORT="$value" ;;
            BEACON_PUBLIC_BASE_URL)  BEACON_PUBLIC_BASE_URL="$value" ;;
            BEACON_AGENT_INSTALL_TOKEN) BEACON_AGENT_INSTALL_TOKEN="$value" ;;
            BEACON_AUTH_SIGNING_KEY)    BEACON_AUTH_SIGNING_KEY="$value" ;;
            BEACON_CONTROL_JOIN_TOKEN)  BEACON_CONTROL_JOIN_TOKEN="$value" ;;
        esac
    done < "$env_file"
}

# 版本号解析：v3.0.4 -> 3 0 4
parse_version() {
    ver="$1"
    ver="${ver#v}"
    printf '%s' "$ver" | awk -F. '{print $1, $2, $3}'
}

# 比较版本：$1 >= $2 返回 0
version_gte() {
    old="$1"; new="$2"
    set -- $(parse_version "$old")
    o1="${1:-0}"; o2="${2:-0}"; o3="${3:-0}"
    set -- $(parse_version "$new")
    n1="${1:-0}"; n2="${2:-0}"; n3="${3:-0}"
    if [ "$n1" -gt "$o1" ]; then return 0; fi
    if [ "$n1" -lt "$o1" ]; then return 1; fi
    if [ "$n2" -gt "$o2" ]; then return 0; fi
    if [ "$n2" -lt "$o2" ]; then return 1; fi
    if [ "$n3" -ge "$o3" ]; then return 0; fi
    return 1
}

# 获取目标版本
resolve_version() {
    current="${1:-}"
    if [ -n "$BEACON_VERSION" ]; then
        if [ -n "$current" ] && [ "$current" != "$BEACON_VERSION" ]; then
            if ! version_gte "$current" "$BEACON_VERSION"; then
                die "Beacon 不支持从 $current 降级到 $BEACON_VERSION"
            fi
        fi
        printf '%s' "$BEACON_VERSION"
        return
    fi
    BEACON_VERSION="$DEFAULT_VERSION"
    printf '%s' "$BEACON_VERSION"
}

# 生成 .env 文件
generate_env() {
    env_file="$1"

    # 保存命令行/环境变量显式传入的值（高优先级）
    cli_image="${BEACON_IMAGE:-}"
    cli_arch="${BEACON_ARCH:-}"
    cli_version="${BEACON_VERSION:-}"
    cli_http_port="${BEACON_HTTP_PORT:-}"
    cli_control_port="${BEACON_CONTROL_PORT:-}"
    cli_public_url="${BEACON_PUBLIC_BASE_URL:-}"
    cli_agent_token="${BEACON_AGENT_INSTALL_TOKEN:-}"
    cli_auth_key="${BEACON_AUTH_SIGNING_KEY:-}"
    cli_join_token="${BEACON_CONTROL_JOIN_TOKEN:-}"

    # 读取已有配置（低优先级）
    read_env "$env_file"

    # 恢复命令行/环境变量传入的值
    [ -n "$cli_image" ]       && BEACON_IMAGE="$cli_image"
    [ -n "$cli_arch" ]        && BEACON_ARCH="$cli_arch"
    [ -n "$cli_version" ]     && BEACON_VERSION="$cli_version"
    [ -n "$cli_http_port" ]   && BEACON_HTTP_PORT="$cli_http_port"
    [ -n "$cli_control_port" ] && BEACON_CONTROL_PORT="$cli_control_port"
    [ -n "$cli_public_url" ]  && BEACON_PUBLIC_BASE_URL="$cli_public_url"
    [ -n "$cli_agent_token" ] && BEACON_AGENT_INSTALL_TOKEN="$cli_agent_token"
    [ -n "$cli_auth_key" ]    && BEACON_AUTH_SIGNING_KEY="$cli_auth_key"
    [ -n "$cli_join_token" ]  && BEACON_CONTROL_JOIN_TOKEN="$cli_join_token"

    # 参数收集（环境变量/命令行优先 → 已有配置 → 交互 → 默认）
    BASE_URL="${BEACON_BASE_URL:-$DEFAULT_BASE_URL}"
    BASE_URL="${BASE_URL%/}"

    INSTALL_DIR="${INSTALL_DIR:-$(prompt "安装目录" "$DEFAULT_INSTALL_DIR")}"
    BEACON_VERSION="${BEACON_VERSION:-$DEFAULT_VERSION}"
    # 若运行未经过 precheck（例如直接生成 .env），兜底用 uname -m 推断一次
    if [ -z "${TARGET_ARCH:-}" ]; then
        machine="$(uname -m 2>/dev/null || true)"
        case "$machine" in
            x86_64|amd64) TARGET_ARCH="amd64" ;;
            aarch64|arm64) TARGET_ARCH="arm64" ;;
            *) TARGET_ARCH="amd64" ;;
        esac
    fi
    BEACON_ARCH="${BEACON_ARCH:-$TARGET_ARCH}"
    case "$BEACON_ARCH" in
        amd64|arm64) ;;
        *) die "BEACON_ARCH 非法：$BEACON_ARCH（应为 amd64 或 arm64）" ;;
    esac
    # amd64 与 arm64 走不同的镜像仓库：
    #   amd64 -> amuluze/beacon
    #   arm64 -> amuluze/beacon-arm
    case "$BEACON_ARCH" in
        amd64) DEFAULT_IMAGE="registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:${BEACON_VERSION}" ;;
        arm64) DEFAULT_IMAGE="registry.cn-hangzhou.aliyuncs.com/amuluze/beacon-arm:${BEACON_VERSION}" ;;
    esac
    BEACON_IMAGE="${BEACON_IMAGE:-$(prompt "Beacon 镜像" "$DEFAULT_IMAGE")}"
    BEACON_HTTP_PORT="${BEACON_HTTP_PORT:-$(prompt "Web 控制台宿主端口" "$DEFAULT_HTTP_PORT")}"
    BEACON_CONTROL_PORT="${BEACON_CONTROL_PORT:-$(prompt "Agent 控制端口" "$DEFAULT_CONTROL_PORT")}"
    BEACON_PUBLIC_BASE_URL="${BEACON_PUBLIC_BASE_URL:-$(prompt "对外访问地址" "http://127.0.0.1:$BEACON_HTTP_PORT")}"
    BEACON_AGENT_INSTALL_TOKEN="${BEACON_AGENT_INSTALL_TOKEN:-$(random_secret)}"
    BEACON_AUTH_SIGNING_KEY="${BEACON_AUTH_SIGNING_KEY:-$(random_secret)}"
    BEACON_CONTROL_JOIN_TOKEN="${BEACON_CONTROL_JOIN_TOKEN:-$(random_secret)}"

    # 端口可用性校验
    check_port "$BEACON_HTTP_PORT"
    check_port "$BEACON_CONTROL_PORT"

    validate_secret "BEACON_AGENT_INSTALL_TOKEN" "$BEACON_AGENT_INSTALL_TOKEN"
    validate_secret "BEACON_AUTH_SIGNING_KEY" "$BEACON_AUTH_SIGNING_KEY"
    validate_secret "BEACON_CONTROL_JOIN_TOKEN" "$BEACON_CONTROL_JOIN_TOKEN"

    cat > "$env_file" <<EOF
BEACON_IMAGE=$BEACON_IMAGE
BEACON_VERSION=$BEACON_VERSION
BEACON_ARCH=$BEACON_ARCH
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
    chmod 600 "$env_file"
}

# 下载并校验 compose.yaml
fetch_compose() {
    base_url="$1"; target_dir="$2"
    compose_file="$target_dir/compose.yaml"
    sums_file="$target_dir/SHA256SUMS"

    log "拉取 compose.yaml ..."
    download "$base_url/release/latest/compose.yaml" "$compose_file"
    download "$base_url/release/latest/SHA256SUMS" "$sums_file"

    compose_sha256="$(awk '$2 == "compose.yaml" { print $1 }' "$sums_file")"
    [ -n "$compose_sha256" ] || die "SHA256SUMS 中缺少 compose.yaml"
    verify_sha256 "$compose_sha256" "$compose_file"
    chmod 600 "$compose_file" "$sums_file"
}

# 镜像拉取，带重试
pull_images() {
    target_dir="$1"
    log "拉取镜像（可能需要一些时间）..."
    if (cd "$target_dir" && $DOCKER_COMPOSE pull); then
        return 0
    fi
    warn "首次拉取镜像失败，5 秒后重试 ..."
    sleep 5
    (cd "$target_dir" && $DOCKER_COMPOSE pull) || die "拉取镜像失败"
}

# 启动容器
compose_up() {
    target_dir="$1"
    log "启动 beacon 容器 ..."
    (cd "$target_dir" && $DOCKER_COMPOSE up -d --remove-orphans) || die "启动容器失败"
}

# 等待容器健康
wait_healthy() {
    target_dir="$1"
    container_name="${BEACON_CONTAINER_NAME:-beacon}"
    log "等待 $container_name 健康检查通过（最多 ${HEALTH_TIMEOUT} 秒）..."
    start_at="$(date +%s)"
    while true; do
        if [ "$(($(date +%s) - start_at))" -gt "$HEALTH_TIMEOUT" ]; then
            err "$container_name 在 ${HEALTH_TIMEOUT} 秒内未通过健康检查"
            if docker logs --tail 50 "$container_name" >/dev/null 2>&1; then
                err "最近 50 行日志："
                docker logs --tail 50 "$container_name" | sed 's/^/  /' >&2
            fi
            return 1
        fi
        status="$(docker inspect --format='{{.State.Health.Status}}' "$container_name" 2>/dev/null || true)"
        if [ "$status" = "healthy" ]; then
            log "$container_name 已通过健康检查"
            return 0
        fi
        sleep 3
    done
}

# 安装操作
install() {
    ACTION="install"
    log "即将安装 Beacon ..."

    precheck

    # 安装目录收集
    INSTALL_DIR="${INSTALL_DIR:-$(prompt "安装目录" "$DEFAULT_INSTALL_DIR")}"
    if [ -e "$INSTALL_DIR" ]; then
        die "安装目录 $INSTALL_DIR 已存在，请选择新目录"
    fi
    case "$INSTALL_DIR" in
        /*) ;;
        *) die "$INSTALL_DIR 不是合法的绝对路径" ;;
    esac
    check_disk "$INSTALL_DIR"

    mkdir -p "$INSTALL_DIR"
    cd "$INSTALL_DIR"

    BASE_URL="${BEACON_BASE_URL:-$DEFAULT_BASE_URL}"
    BASE_URL="${BASE_URL%/}"

    # 生成 .env
    generate_env "$INSTALL_DIR/.env"

    # 重新读取生成后的关键变量用于展示
    read_env "$INSTALL_DIR/.env"

    # 下载 compose
    fetch_compose "$BASE_URL" "$INSTALL_DIR"

    mkdir -p data logs

    # TTY 下询问是否在启动前编辑 .env
    if [ -t 0 ] && confirm "是否在启动前编辑 .env？" "n"; then
        "${EDITOR:-vi}" .env
        read_env "$INSTALL_DIR/.env"
    fi

    pull_images "$INSTALL_DIR"
    compose_up "$INSTALL_DIR"
    wait_healthy "$INSTALL_DIR" || die "服务启动失败，请检查日志"

    print_finish
}

# 升级操作
upgrade() {
    ACTION="upgrade"
    log "即将升级 Beacon ..."

    precheck

    INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    if [ ! -d "$INSTALL_DIR" ]; then
        die "未找到安装目录 $INSTALL_DIR，无法升级"
    fi
    cd "$INSTALL_DIR"

    if [ ! -f .env ]; then
        die "未找到 $INSTALL_DIR/.env，无法升级"
    fi
    if [ ! -f compose.yaml ]; then
        die "未找到 $INSTALL_DIR/compose.yaml，无法升级"
    fi

    # 备份当前配置
    BACKUP_ENV="$INSTALL_DIR/.env.bak.$(date +%Y%m%d%H%M%S)"
    BACKUP_COMPOSE="$INSTALL_DIR/compose.yaml.bak.$(date +%Y%m%d%H%M%S)"
    cp -f .env "$BACKUP_ENV"
    cp -f compose.yaml "$BACKUP_COMPOSE"

    # 读取当前版本
    read_env "$INSTALL_DIR/.env"
    current_version="$BEACON_VERSION"

    BASE_URL="${BEACON_BASE_URL:-$DEFAULT_BASE_URL}"
    BASE_URL="${BASE_URL%/}"

    # 生成新版 .env（保留已有配置并允许覆盖版本等）
    BEACON_VERSION="${BEACON_VERSION:-}"
    BEACON_VERSION="$(resolve_version "$current_version")"
    generate_env "$INSTALL_DIR/.env"
    read_env "$INSTALL_DIR/.env"

    # 下载新版 compose
    fetch_compose "$BASE_URL" "$INSTALL_DIR"

    # 停止旧容器
    log "停止旧容器 ..."
    (cd "$INSTALL_DIR" && $DOCKER_COMPOSE down) || warn "停止旧容器失败，将继续尝试启动"

    pull_images "$INSTALL_DIR"
    compose_up "$INSTALL_DIR"
    if wait_healthy "$INSTALL_DIR"; then
        rm -f "$BACKUP_ENV" "$BACKUP_COMPOSE"
        BACKUP_ENV=""
        BACKUP_COMPOSE=""
    else
        err "升级后服务未通过健康检查，尝试回滚 ..."
        (cd "$INSTALL_DIR" && $DOCKER_COMPOSE down) || true
        cp -f "$BACKUP_ENV" .env
        cp -f "$BACKUP_COMPOSE" compose.yaml
        (cd "$INSTALL_DIR" && $DOCKER_COMPOSE up -d) || die "回滚失败，请手动检查 $INSTALL_DIR"
        die "已回滚到升级前状态"
    fi

    print_finish
}

# 重启操作
restart() {
    ACTION="restart"
    log "即将重启 Beacon ..."

    precheck

    INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    if [ ! -d "$INSTALL_DIR" ]; then
        die "未找到安装目录 $INSTALL_DIR"
    fi
    cd "$INSTALL_DIR"

    if [ ! -f compose.yaml ]; then
        die "未找到 $INSTALL_DIR/compose.yaml"
    fi

    log "停止容器 ..."
    (cd "$INSTALL_DIR" && $DOCKER_COMPOSE down) || die "停止容器失败"

    log "启动容器 ..."
    (cd "$INSTALL_DIR" && $DOCKER_COMPOSE up -d) || die "启动容器失败"

    wait_healthy "$INSTALL_DIR" || die "服务启动失败"
    print_finish
}

# 卸载操作
uninstall() {
    ACTION="uninstall"
    log "即将卸载 Beacon ..."

    check_root

    INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
    if [ ! -d "$INSTALL_DIR" ]; then
        die "未找到安装目录 $INSTALL_DIR"
    fi

    if ! confirm "确认卸载 Beacon？该操作会删除 $INSTALL_DIR 下所有数据" "n"; then
        log "已取消卸载"
        return 0
    fi

    if [ -f "$INSTALL_DIR/compose.yaml" ]; then
        (cd "$INSTALL_DIR" && $DOCKER_COMPOSE down --volumes --remove-orphans) || warn "停止容器失败"
    fi

    rm -rf "$INSTALL_DIR"
    log "Beacon 已卸载"
}

# 完成提示
print_finish() {
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
}

# 打印用法
usage() {
    cat <<EOF
Beacon 一键安装脚本

用法：
  sh manager.sh [install|upgrade|uninstall|restart] [KEY=VAL ...]

操作：
  install    安装 Beacon（默认）
  upgrade    升级 Beacon
  uninstall  卸载 Beacon
  restart    重启 Beacon

可覆盖变量：
  BEACON_BASE_URL
  INSTALL_DIR
  BEACON_IMAGE
  BEACON_VERSION
  BEACON_ARCH
  BEACON_HTTP_PORT
  BEACON_CONTROL_PORT
  BEACON_PUBLIC_BASE_URL
  BEACON_AGENT_INSTALL_TOKEN
  BEACON_AUTH_SIGNING_KEY
  BEACON_CONTROL_JOIN_TOKEN
EOF
}

# 解析位置参数和 KEY=VAL 参数
parse_args() {
    while [ $# -gt 0 ]; do
        case "$1" in
            install|upgrade|uninstall|restart)
                if [ -z "$ACTION" ]; then
                    ACTION="$1"
                else
                    die "只能指定一个操作"
                fi
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *=*)
                export "$1"
                ;;
            *)
                die "无法识别的参数：$1（格式应为 install/upgrade/uninstall/restart 或 KEY=VAL）"
                ;;
        esac
        shift
    done

    if [ -z "$ACTION" ]; then
        ACTION="install"
    fi
}

# 入口
main() {
    parse_args "$@"

    case "$ACTION" in
        install) install ;;
        upgrade) upgrade ;;
        uninstall) uninstall ;;
        restart) restart ;;
        *) die "未知操作：$ACTION" ;;
    esac
}

main "$@"
