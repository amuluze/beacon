package service

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultAmprobeImage = "registry.cn-hangzhou.aliyuncs.com/amuluze/amprobe:latest"

func (a *Router) WebsiteInstallScript(ctx *fiber.Ctx) error {
	baseURL := requestBaseURL(ctx)
	ctx.Type("sh")
	ctx.Set("Content-Disposition", `attachment; filename="install.sh"`)
	return ctx.SendString(buildWebsiteInstallScript(baseURL))
}

func (a *Router) WebsiteCompose(ctx *fiber.Ctx) error {
	ctx.Type("yaml")
	ctx.Set("Content-Disposition", `attachment; filename="compose.yaml"`)
	return ctx.SendString(websiteComposeYAML())
}

func buildWebsiteInstallScript(baseURL string) string {
	baseURL = strings.TrimRight(baseURL, "/")

	return fmt.Sprintf(`#!/bin/sh
set -eu

BASE_URL=%s
DEFAULT_INSTALL_DIR="/data/amprobe"
DEFAULT_IMAGE=%s

prompt() {
  label="$1"
  default_value="$2"
  if [ -t 0 ]; then
    printf "%%s [%%s]: " "$label" "$default_value" >&2
    read -r value || value=""
    if [ -n "$value" ]; then
      printf "%%s" "$value"
      return
    fi
  fi
  printf "%%s" "$default_value"
}

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 24
  else
    date +%%s | cksum | awk '{print $1}'
  fi
}

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing command: $1" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd docker

if docker compose version >/dev/null 2>&1; then
  DOCKER_COMPOSE="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DOCKER_COMPOSE="docker-compose"
else
  echo "missing docker compose plugin or docker-compose" >&2
  exit 1
fi

INSTALL_DIR="${INSTALL_DIR:-$(prompt "Install directory" "$DEFAULT_INSTALL_DIR")}"
AMPROBE_IMAGE="${AMPROBE_IMAGE:-$(prompt "Amprobe image" "$DEFAULT_IMAGE")}"
AMPROBE_HTTP_PORT="${AMPROBE_HTTP_PORT:-$(prompt "Web console host port" "1443")}"
AMPROBE_CONTROL_PORT="${AMPROBE_CONTROL_PORT:-$(prompt "Agent control host port" "17000")}"
AMPROBE_PUBLIC_BASE_URL="${AMPROBE_PUBLIC_BASE_URL:-$(prompt "Public base URL" "http://127.0.0.1:$AMPROBE_HTTP_PORT")}"
AMPROBE_AGENT_INSTALL_TOKEN="${AMPROBE_AGENT_INSTALL_TOKEN:-$(random_secret)}"
AMPROBE_AUTH_SIGNING_KEY="${AMPROBE_AUTH_SIGNING_KEY:-$(random_secret)}"

mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

curl -fsSL "$BASE_URL/download/compose.yaml" -o compose.yaml

cat > .env <<EOF
AMPROBE_IMAGE=$AMPROBE_IMAGE
AMPROBE_CONTAINER_NAME=amprobe
AMPROBE_HTTP_PORT=$AMPROBE_HTTP_PORT
AMPROBE_CONTROL_PORT=$AMPROBE_CONTROL_PORT
AMPROBE_DATA_DIR=./data
AMPROBE_LOG_DIR=./logs
AMPROBE_DB_NAME=/app/data/probe
AMPROBE_PUBLIC_BASE_URL=$AMPROBE_PUBLIC_BASE_URL
AMPROBE_AGENT_INSTALL_TOKEN=$AMPROBE_AGENT_INSTALL_TOKEN
AMPROBE_AUTH_SIGNING_KEY=$AMPROBE_AUTH_SIGNING_KEY
EOF

mkdir -p data logs

if [ -t 0 ]; then
  printf "Edit .env before start? [y/N]: " >&2
  read -r edit_env || edit_env=""
  case "$edit_env" in
    y|Y|yes|YES)
      "${EDITOR:-vi}" .env
      ;;
  esac
fi

$DOCKER_COMPOSE pull
$DOCKER_COMPOSE up -d

echo "Amprobe started."
echo "Install directory: $INSTALL_DIR"
echo "Web console: $AMPROBE_PUBLIC_BASE_URL"
echo "Agent install token is stored in $INSTALL_DIR/.env"
`, shellQuote(baseURL), shellQuote(defaultAmprobeImage))
}

func websiteComposeYAML() string {
	return `services:
  amprobe:
    image: ${AMPROBE_IMAGE:-registry.cn-hangzhou.aliyuncs.com/amuluze/amprobe:latest}
    container_name: ${AMPROBE_CONTAINER_NAME:-amprobe}
    restart: unless-stopped
    ports:
      - "${AMPROBE_HTTP_PORT:-1443}:80"
      - "${AMPROBE_CONTROL_PORT:-17000}:${AMPROBE_CONTROL_PORT:-17000}"
    volumes:
      - ${AMPROBE_DATA_DIR:-./data}:/app/data
      - ${AMPROBE_LOG_DIR:-./logs}:/app/logs
    environment:
      AMPROBE_DB_NAME: ${AMPROBE_DB_NAME:-/app/data/probe}
      AMPROBE_PUBLIC_BASE_URL: ${AMPROBE_PUBLIC_BASE_URL:-}
      AMPROBE_AGENT_INSTALL_TOKEN: ${AMPROBE_AGENT_INSTALL_TOKEN:-change-me}
      AMPROBE_AUTH_SIGNING_KEY: ${AMPROBE_AUTH_SIGNING_KEY:-amprobe}
      AMPROBE_CONTROL_PORT: ${AMPROBE_CONTROL_PORT:-17000}
    command:
      - /bin/sh
      - -c
      - |
        set -eu
        sed -i "s|DBName = \".*\"|DBName = \"$${AMPROBE_DB_NAME:-/app/data/probe}\"|" /app/configs/config.toml
        sed -i "s|SigningKey = \".*\"|SigningKey = \"$${AMPROBE_AUTH_SIGNING_KEY:-amprobe}\"|" /app/configs/config.toml
        sed -i "s|Address = \".*\"|Address = \":$${AMPROBE_CONTROL_PORT:-17000}\"|" /app/configs/config.toml
        sed -i "s|Token = \".*\"|Token = \"$${AMPROBE_AGENT_INSTALL_TOKEN:-change-me}\"|" /app/configs/config.toml
        sed -i "s|PublicBaseURL = \".*\"|PublicBaseURL = \"$${AMPROBE_PUBLIC_BASE_URL:-}\"|" /app/configs/config.toml
        sed -i "s|ControlPort = .*|ControlPort = $${AMPROBE_CONTROL_PORT:-17000}|" /app/configs/config.toml
        exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
`
}

func requestBaseURL(ctx *fiber.Ctx) string {
	scheme := ctx.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = ctx.Protocol()
	}
	host := ctx.Get("Host")
	if host == "" {
		host = ctx.Hostname()
	}
	return strings.TrimRight(scheme+"://"+host, "/")
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
