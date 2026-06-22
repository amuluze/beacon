package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const defaultAgentInstallPackageDir = "/app/downloads/collia"

var safeInstallNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func (a *Router) AgentInstallScript(ctx *fiber.Ctx) error {
	if !a.config.AgentInstall.Enable {
		return fiber.ErrNotFound
	}

	node := ctx.Query("node")
	if node == "" {
		return fiber.NewError(http.StatusBadRequest, "missing node")
	}
	if !isSafeInstallName(node) {
		return fiber.NewError(http.StatusBadRequest, "invalid node")
	}

	osType := ctx.Query("os_type", "linux")
	if !isSafeInstallName(osType) {
		return fiber.NewError(http.StatusBadRequest, "invalid os_type")
	}

	baseURL := a.config.AgentInstall.PublicBaseURL
	if baseURL == "" {
		baseURL = requestBaseURL(ctx)
	}

	ctx.Type("sh")
	return ctx.SendString(buildAgentInstallScript(baseURL, node, osType, a.agentInstallRPCPort(), a.config.AgentInstall.TLSEnable))
}

func (a *Router) AgentInstallPackage(ctx *fiber.Ctx) error {
	if err := a.verifyAgentInstallToken(ctx); err != nil {
		return err
	}

	osType := ctx.Query("os_type", "linux")
	arch := ctx.Query("arch", "amd64")
	if !isSafeInstallName(osType) || !isSafeInstallName(arch) {
		return fiber.NewError(http.StatusBadRequest, "invalid package selector")
	}

	packagePath, err := safeJoin(a.agentInstallPackageDir(), osType, arch, "collia.install")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	if _, err := os.Stat(packagePath); err != nil {
		return fiber.NewError(http.StatusNotFound, "collia install package not found")
	}
	return ctx.Download(packagePath, "collia.install")
}

func (a *Router) AgentInstallConfig(ctx *fiber.Ctx) error {
	if err := a.verifyAgentInstallToken(ctx); err != nil {
		return err
	}

	node := ctx.Query("node")
	if node == "" || !isSafeInstallName(node) {
		return fiber.NewError(http.StatusBadRequest, "invalid node")
	}

	ctx.Type("yaml")
	return ctx.SendString(a.buildColliaConfig(node))
}

func (a *Router) AgentInstallCerts(ctx *fiber.Ctx) error {
	if err := a.verifyAgentInstallToken(ctx); err != nil {
		return err
	}

	node := ctx.Query("node")
	if node == "" || !isSafeInstallName(node) {
		return fiber.NewError(http.StatusBadRequest, "invalid node")
	}

	certsPath, err := safeJoin(a.agentInstallPackageDir(), "certs", node+".tar.gz")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	if _, err := os.Stat(certsPath); err != nil {
		return fiber.NewError(http.StatusNotFound, "collia cert package not found")
	}
	return ctx.Download(certsPath, node+"-certs.tar.gz")
}

func (a *Router) verifyAgentInstallToken(ctx *fiber.Ctx) error {
	if !a.config.AgentInstall.Enable {
		return fiber.ErrNotFound
	}
	if a.config.AgentInstall.Token == "" {
		return fiber.NewError(http.StatusForbidden, "agent install token is not configured")
	}
	token := ctx.Get("X-Install-Token")
	if token == "" {
		token = strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")
	}
	if token == "" {
		token = ctx.Query("token")
	}
	if token != a.config.AgentInstall.Token {
		return fiber.NewError(http.StatusUnauthorized, "invalid install token")
	}
	return nil
}

func (a *Router) buildColliaConfig(node string) string {
	certDir := a.config.AgentInstall.CertDir
	if certDir == "" {
		certDir = "/etc/collia/certs"
	}

	return fmt.Sprintf(`rpc:
  network: tcp
  address: 0.0.0.0:%d
  tls:
    enable: %t
    cert_dir: %s
    client_names:
      - amprobe
log:
  output: /data/amprobe/logs/collia/collia.log
  level: info
  rotation: 1
  max_age: 7
task:
  interval: 30
  max_age: 1
  disk:
    devices:
      - vda2
  ethernet:
    names:
      - eth0
db:
  dbtype: sqlite
  dbname: /data/amprobe/resources/collia/storage/collia
variables:
  image_tag: latest
  host_prefix: /data/amprobe
  container_prefix: /
  node: %s
`, a.agentInstallRPCPort(), a.config.AgentInstall.TLSEnable, certDir, node)
}

func (a *Router) agentInstallPackageDir() string {
	if a.config.AgentInstall.PackageDir != "" {
		return a.config.AgentInstall.PackageDir
	}
	return defaultAgentInstallPackageDir
}

func (a *Router) agentInstallRPCPort() int {
	if a.config.AgentInstall.RPCPort > 0 {
		return a.config.AgentInstall.RPCPort
	}
	return 18080
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

func buildAgentInstallScript(baseURL string, node string, osType string, rpcPort int, tlsEnabled bool) string {
	return fmt.Sprintf(`#!/bin/sh
set -eu

BASE_URL=%s
NODE=%s
OS_TYPE=%s
RPC_PORT=%s
TLS_ENABLED=%s
TOKEN=""
ARCH=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --token=*) TOKEN="${1#*=}" ;;
    --token) shift; TOKEN="${1:-}" ;;
    --arch=*) ARCH="${1#*=}" ;;
    --arch) shift; ARCH="${1:-}" ;;
    *) echo "unknown argument: $1" >&2; exit 1 ;;
  esac
  shift
done

if [ -z "$TOKEN" ]; then
  echo "missing --token" >&2
  exit 1
fi

if [ -z "$ARCH" ]; then
  case "$(uname -m)" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
  esac
fi

WORKDIR="$(mktemp -d /tmp/collia-install.XXXXXX)"
cleanup() {
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

download() {
  curl -kfsSL -H "X-Install-Token: $TOKEN" "$1" -o "$2"
}

mkdir -p /etc/collia /data/amprobe/resources/collia/storage /data/amprobe/logs/collia

download "$BASE_URL/api/v1/host/install/package?node=$NODE&os_type=$OS_TYPE&arch=$ARCH" "$WORKDIR/collia.install"
download "$BASE_URL/api/v1/host/install/config?node=$NODE" "$WORKDIR/config.yml"

chmod +x "$WORKDIR/collia.install"
"$WORKDIR/collia.install"
install -m 0644 "$WORKDIR/config.yml" /etc/collia/config.yml

if [ "$TLS_ENABLED" = "true" ]; then
  mkdir -p /etc/collia/certs
  download "$BASE_URL/api/v1/host/install/certs?node=$NODE" "$WORKDIR/certs.tar.gz"
  tar -xzf "$WORKDIR/certs.tar.gz" -C /etc/collia/certs
fi

collia install || true
collia stop || true
collia start

echo "collia installed and started on tcp://0.0.0.0:$RPC_PORT"
`, shellQuote(baseURL), shellQuote(node), shellQuote(osType), shellQuote(strconv.Itoa(rpcPort)), shellQuote(strconv.FormatBool(tlsEnabled)))
}

func isSafeInstallName(s string) bool {
	return safeInstallNamePattern.MatchString(s)
}

func safeJoin(base string, elems ...string) (string, error) {
	if base == "" {
		return "", fmt.Errorf("empty base dir")
	}
	cleanBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	parts := append([]string{cleanBase}, elems...)
	joined := filepath.Clean(filepath.Join(parts...))
	rel, err := filepath.Rel(cleanBase, joined)
	if err != nil {
		return "", err
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", fmt.Errorf("path escapes base dir")
	}
	return joined, nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
