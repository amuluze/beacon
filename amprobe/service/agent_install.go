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

const colliaBinaryName = "collia"

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

	baseURL := a.config.AgentInstall.PublicBaseURL
	if baseURL == "" {
		baseURL = requestBaseURL(ctx)
	}

	ctx.Type("sh")
	return ctx.SendString(buildAgentInstallScript(baseURL, node))
}

func (a *Router) AgentInstallPackage(ctx *fiber.Ctx) error {
	if err := a.verifyAgentInstallToken(ctx); err != nil {
		return err
	}

	arch := ctx.Query("arch", "amd64")
	if arch != "amd64" && arch != "arm64" {
		return fiber.NewError(http.StatusBadRequest, "unsupported arch")
	}

	packagePath, err := safeJoin(a.agentInstallPackageDir(), arch, colliaBinaryName)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	if _, err := os.Stat(packagePath); err != nil {
		return fiber.NewError(http.StatusNotFound, "collia binary not found")
	}
	return ctx.Download(packagePath, colliaBinaryName)
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
	reportURL := ""
	controlServer := "127.0.0.1:17000"
	controlPort := a.config.AgentInstall.ControlPort
	if controlPort == 0 {
		controlPort = 17000
	}
	if a.config.AgentInstall.PublicBaseURL != "" {
		reportURL = a.config.AgentInstall.PublicBaseURL + "/api/v1/host/report"
		controlServer = extractHost(a.config.AgentInstall.PublicBaseURL) + ":" + strconv.Itoa(controlPort)
	}

	return fmt.Sprintf(`control:
  server: %s
  agent_id: %s
  join_token: "%s"
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
  report:
    url: "%s"
    token: "%s"
    agent_id: %s
db:
  dbtype: sqlite
  dbname: /data/amprobe/resources/collia/storage/collia
variables:
  image_tag: latest
  host_prefix: /data/amprobe
  container_prefix: /
  node: %s
`, controlServer, node, controlJoinToken(a.config), reportURL, a.config.AgentInstall.Token, node, node)
}

func (a *Router) agentInstallPackageDir() string {
	if a.config.AgentInstall.PackageDir != "" {
		return a.config.AgentInstall.PackageDir
	}
	return defaultAgentInstallPackageDir
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

func buildAgentInstallScript(baseURL string, node string) string {
	return fmt.Sprintf(`#!/bin/sh
set -eu

BASE_URL=%s
NODE=%s
TOKEN=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --token=*) TOKEN="${1#*=}" ;;
    --token) shift; TOKEN="${1:-}" ;;
    *) echo "unknown argument: $1" >&2; exit 1 ;;
  esac
  shift
done

if [ -z "$TOKEN" ]; then
  echo "missing --token" >&2
  exit 1
fi

case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

WORKDIR="$(mktemp -d /tmp/collia-install.XXXXXX)"
cleanup() {
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

download() {
  curl -kfsSL -H "X-Install-Token: $TOKEN" "$1" -o "$2"
}

mkdir -p /etc/collia /data/amprobe/resources/collia/storage /data/amprobe/logs/collia /usr/sbin

download "$BASE_URL/api/v1/host/install/package?arch=$ARCH" "$WORKDIR/collia"
download "$BASE_URL/api/v1/host/install/config?node=$NODE" "$WORKDIR/config.yml"

install -m 0755 "$WORKDIR/collia" /usr/sbin/collia
install -m 0644 "$WORKDIR/config.yml" /etc/collia/config.yml

collia install || true
collia stop || true
collia start

echo "collia installed and started, reverse tunnel -> $BASE_URL"
`, shellQuote(baseURL), shellQuote(node))
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

// extractHost extracts the host portion from a URL (e.g., "http://example.com:8000" -> "example.com").
func extractHost(rawURL string) string {
	if rawURL == "" {
		return "127.0.0.1"
	}
	// Strip protocol prefix
	clean := rawURL
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(clean, prefix) {
			clean = strings.TrimPrefix(clean, prefix)
			break
		}
	}
	// Strip port and path
	if idx := strings.Index(clean, ":"); idx > 0 {
		clean = clean[:idx]
	} else if idx := strings.Index(clean, "/"); idx > 0 {
		clean = clean[:idx]
	}
	return clean
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
