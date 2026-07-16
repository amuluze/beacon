// Package service
// Date: 2026/07/16
// Author: Amu
// Description: 版本检查接口单测
package service

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    int
	}{
		// 语义化版本比较
		{"patch behind", "v3.0.3", "v3.0.4", -1},
		{"minor behind", "v3.0.4", "v3.1.0", -1},
		{"major behind", "v2.9.9", "v3.0.0", -1},
		{"equal", "v3.0.4", "v3.0.4", 0},
		{"ahead", "v3.0.5", "v3.0.4", 1},
		// 缺前缀 v
		{"no prefix current", "3.0.3", "v3.0.4", -1},
		{"no prefix latest", "v3.0.4", "3.0.4", 0},
		// 段数不一致
		{"shorter current", "v3.0", "v3.0.4", -1},
		{"shorter latest", "v3.0.4", "v3.0", 1},
		// 非法格式回退：两端均非法走字符串比较
		{"both invalid", "main", "dev", 1}, // strings.Compare: 'm' > 'd'
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareVersions(tt.current, tt.latest); got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.current, tt.latest, got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	if got := parseVersion("v3.0.4"); got == nil || len(got) != 3 || got[0] != 3 || got[2] != 4 {
		t.Errorf("parseVersion(v3.0.4) = %v, want [3 0 4]", got)
	}
	if got := parseVersion("abc"); got != nil {
		t.Errorf("parseVersion(abc) = %v, want nil", got)
	}
	if got := parseVersion(""); got != nil {
		t.Errorf("parseVersion(\"\") = %v, want nil", got)
	}
}

func TestLoadVersionManifest(t *testing.T) {
	dir := t.TempDir()
	latestDir := filepath.Join(dir, "latest")
	if err := os.MkdirAll(latestDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	t.Run("valid manifest", func(t *testing.T) {
		content := `{"latest_version":"v3.0.4","min_required_version":"v3.0.0","release_notes":"fix","published_at":"2026-07-16"}`
		if err := os.WriteFile(filepath.Join(latestDir, "version.json"), []byte(content), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		m, err := loadVersionManifest(dir)
		if err != nil {
			t.Fatalf("loadVersionManifest: %v", err)
		}
		if m.LatestVersion != "v3.0.4" {
			t.Errorf("LatestVersion = %q, want v3.0.4", m.LatestVersion)
		}
		if m.MinRequiredVersion != "v3.0.0" {
			t.Errorf("MinRequiredVersion = %q, want v3.0.0", m.MinRequiredVersion)
		}
	})

	t.Run("empty latest_version rejected", func(t *testing.T) {
		content := `{"latest_version":""}`
		if err := os.WriteFile(filepath.Join(latestDir, "version.json"), []byte(content), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if _, err := loadVersionManifest(dir); err == nil {
			t.Errorf("空 latest_version 应返回错误")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		if _, err := loadVersionManifest(t.TempDir()); err == nil {
			t.Errorf("文件缺失应返回错误")
		}
	})
}

func TestVersionLatestInvalidCurrentDoesNotReportUpdate(t *testing.T) {
	dir := t.TempDir()
	latestDir := filepath.Join(dir, "latest")
	if err := os.MkdirAll(latestDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := `{"latest_version":"v3.0.4","min_required_version":"v3.0.0","release_notes":"fix","published_at":"2026-07-16"}`
	if err := os.WriteFile(filepath.Join(latestDir, "version.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	router := &Router{config: &Config{Release: Release{Dir: dir}}}
	app := fiber.New()
	app.Get("/version", router.VersionLatest)
	resp, err := app.Test(httptest.NewRequest("GET", "/version?current=not-semver", nil), -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	var body VersionLatestResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.UpdateAvailable {
		t.Fatal("非法 current 版本不应报告可更新")
	}
}
