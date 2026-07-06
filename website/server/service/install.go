package service

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultBeaconImage = "registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:latest"

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
DEFAULT_INSTALL_DIR="/data/beacon"
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
BEACON_IMAGE="${BEACON_IMAGE:-$(prompt "Beacon image" "$DEFAULT_IMAGE")}"
BEACON_HTTP_PORT="${BEACON_HTTP_PORT:-$(prompt "Web console host port" "1443")}"
BEACON_CONTROL_PORT="${BEACON_CONTROL_PORT:-$(prompt "Agent control host port" "17000")}"
BEACON_PUBLIC_BASE_URL="${BEACON_PUBLIC_BASE_URL:-$(prompt "Public base URL" "http://127.0.0.1:$BEACON_HTTP_PORT")}"
BEACON_AGENT_INSTALL_TOKEN="${BEACON_AGENT_INSTALL_TOKEN:-$(random_secret)}"
BEACON_AUTH_SIGNING_KEY="${BEACON_AUTH_SIGNING_KEY:-$(random_secret)}"

mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

curl -fsSL "$BASE_URL/download/compose.yaml" -o compose.yaml

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

echo "Beacon started."
echo "Install directory: $INSTALL_DIR"
echo "Web console: $BEACON_PUBLIC_BASE_URL"
echo "Agent install token is stored in $INSTALL_DIR/.env"
`, shellQuote(baseURL), shellQuote(defaultBeaconImage))
}

func websiteComposeYAML() string {
	return `services:
  beacon:
    image: ${BEACON_IMAGE:-registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:latest}
    container_name: ${BEACON_CONTAINER_NAME:-beacon}
    restart: unless-stopped
    ports:
      - "${BEACON_HTTP_PORT:-1443}:80"
      - "${BEACON_CONTROL_PORT:-17000}:${BEACON_CONTROL_PORT:-17000}"
    volumes:
      - ${BEACON_DATA_DIR:-./data}:/app/data
      - ${BEACON_LOG_DIR:-./logs}:/app/logs
    environment:
      BEACON_DB_NAME: ${BEACON_DB_NAME:-/app/data/beacon}
      BEACON_PUBLIC_BASE_URL: ${BEACON_PUBLIC_BASE_URL:-}
      BEACON_AGENT_INSTALL_TOKEN: ${BEACON_AGENT_INSTALL_TOKEN:-change-me}
      BEACON_AUTH_SIGNING_KEY: ${BEACON_AUTH_SIGNING_KEY:-beacon}
      BEACON_CONTROL_PORT: ${BEACON_CONTROL_PORT:-17000}
    command:
      - /bin/sh
      - -c
      - |
        set -eu
        sed -i "s|DBName = \".*\"|DBName = \"$${BEACON_DB_NAME:-/app/data/beacon}\"|" /app/configs/config.toml
        sed -i "s|SigningKey = \".*\"|SigningKey = \"$${BEACON_AUTH_SIGNING_KEY:-beacon}\"|" /app/configs/config.toml
        sed -i "s|Address = \".*\"|Address = \":$${BEACON_CONTROL_PORT:-17000}\"|" /app/configs/config.toml
        sed -i "s|Token = \".*\"|Token = \"$${BEACON_AGENT_INSTALL_TOKEN:-change-me}\"|" /app/configs/config.toml
        sed -i "s|PublicBaseURL = \".*\"|PublicBaseURL = \"$${BEACON_PUBLIC_BASE_URL:-}\"|" /app/configs/config.toml
        sed -i "s|ControlPort = .*|ControlPort = $${BEACON_CONTROL_PORT:-17000}|" /app/configs/config.toml
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
