// Package service
// Date: 2026/07/15
// Author: Amu
// Description: unit tests for install script / compose generation
package service

import (
	"strings"
	"testing"
)

func TestShellQuote(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "''"},
		{"plain", "https://example.com", "'https://example.com'"},
		// 单引号必须被转义，避免拼接出的脚本被注入切断
		{"single quote", "a'b", "'a'\"'\"'b'"},
		{"trailing slash", "https://example.com/", "'https://example.com/'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shellQuote(tt.in); got != tt.want {
				t.Errorf("shellQuote(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBuildWebsiteInstallScript_TrimBaseURL(t *testing.T) {
	// 末尾斜杠需被裁剪，避免脚本里出现 //download/install.sh
	withSlash := buildWebsiteInstallScript("https://help.beacon.amuluze.com/")
	withoutSlash := buildWebsiteInstallScript("https://help.beacon.amuluze.com")

	if strings.Contains(withSlash, "com//download") {
		t.Errorf("baseURL 末尾斜杠未被裁剪，出现双斜杠")
	}
	if withSlash != withoutSlash {
		t.Errorf("带/不带末尾斜杠应产出相同脚本")
	}
}

func TestBuildWebsiteInstallScript_ContainsExpectedPieces(t *testing.T) {
	script := buildWebsiteInstallScript("https://help.beacon.amuluze.com")

	mustContain := []string{
		"BASE_URL='https://help.beacon.amuluze.com'",
		"/download/compose.yaml",
		"need_cmd curl",
		"need_cmd docker",
		// docker compose 命令以变量形式调用，同时校验探测逻辑
		"docker compose version",
		"$DOCKER_COMPOSE up -d",
	}
	for _, want := range mustContain {
		if !strings.Contains(script, want) {
			t.Errorf("install 脚本缺少片段: %q", want)
		}
	}
}

func TestWebsiteComposeYAML(t *testing.T) {
	yaml := websiteComposeYAML()

	mustContain := []string{
		"services:",
		"beacon:",
		"${BEACON_HTTP_PORT:-1443}:80",
		"${BEACON_CONTROL_PORT:-17000}",
		"restart: unless-stopped",
	}
	for _, want := range mustContain {
		if !strings.Contains(yaml, want) {
			t.Errorf("compose yaml 缺少片段: %q", want)
		}
	}
}
