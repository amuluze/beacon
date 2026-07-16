// Package service
// Date: 2026/07/16
// Author: Amu
// Description: 版本检查接口 —— 读取发布清单 version.json，供 beacon 定期轮询判断是否有新版本。
package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// versionManifestFile 发布清单文件名，位于 <Release.Dir>/latest/ 下，作为版本信息 SSOT。
const versionManifestFile = "latest/version.json"

// VersionManifest 发布清单，由 release/latest/version.json 提供。
type VersionManifest struct {
	LatestVersion      string `json:"latest_version"`
	MinRequiredVersion string `json:"min_required_version"`
	ReleaseNotes       string `json:"release_notes"`
	PublishedAt        string `json:"published_at"`
}

// VersionLatestResponse 版本检查响应。
type VersionLatestResponse struct {
	LatestVersion      string `json:"latest_version"`
	MinRequiredVersion string `json:"min_required_version"`
	UpdateAvailable    bool   `json:"update_available"`
	ReleaseNotes       string `json:"release_notes"`
	PublishedAt        string `json:"published_at"`
}

// VersionLatest GET /api/v1/version/latest?current=vX.Y.Z
// 读取发布清单并结合请求方的当前版本判断是否需要更新。
// current 缺省或非法时 update_available 恒为 false（无法判定）。
func (a *Router) VersionLatest(ctx *fiber.Ctx) error {
	manifest, err := loadVersionManifest(a.config.Release.Dir)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, fmt.Sprintf("load version manifest: %v", err))
	}

	resp := VersionLatestResponse{
		LatestVersion:      manifest.LatestVersion,
		MinRequiredVersion: manifest.MinRequiredVersion,
		ReleaseNotes:       manifest.ReleaseNotes,
		PublishedAt:        manifest.PublishedAt,
	}

	if current := strings.TrimSpace(ctx.Query("current")); current != "" {
		resp.UpdateAvailable = compareVersions(current, manifest.LatestVersion) < 0
	}

	return ctx.JSON(resp)
}

// loadVersionManifest 从发布目录读取并解析版本清单。
func loadVersionManifest(releaseDir string) (VersionManifest, error) {
	path := filepath.Join(releaseDir, versionManifestFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return VersionManifest{}, fmt.Errorf("read %s: %w", path, err)
	}

	var manifest VersionManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return VersionManifest{}, fmt.Errorf("unmarshal %s: %w", path, err)
	}
	if strings.TrimSpace(manifest.LatestVersion) == "" {
		return VersionManifest{}, fmt.Errorf("%s: latest_version 为空", path)
	}
	return manifest, nil
}

// compareVersions 语义化版本比较（ vX.Y.Z 形式）。
// 返回 -1 / 0 / 1 表示 current 小于 / 等于 / 大于 latest。
// 任一端非法（无法解析为数字段）时回退字符串比较，保证可判定且不 panic。
func compareVersions(current, latest string) int {
	c := parseVersion(current)
	l := parseVersion(latest)

	// 两端均可解析为数字段时按段比较
	if c != nil && l != nil {
		for i := 0; i < len(c) || i < len(l); i++ {
			var ci, li int
			if i < len(c) {
				ci = c[i]
			}
			if i < len(l) {
				li = l[i]
			}
			if ci < li {
				return -1
			}
			if ci > li {
				return 1
			}
		}
		return 0
	}

	// 回退：字符串比较
	switch {
	case c == nil && l == nil:
		return strings.Compare(current, latest)
	default:
		// 可解析的一端视为较新
		if c != nil {
			return 1
		}
		return -1
	}
}

// parseVersion 解析 "v3.0.4" 为 [3,0,4]；非法返回 nil。
func parseVersion(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ".")
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil
		}
		nums = append(nums, n)
	}
	return nums
}
