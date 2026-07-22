// Package service
// Description: VersionChecker 定期轮询 beacon-help 的版本清单接口，
// 判断 beacon 是否有新版本。检测到更新时仅记日志 + 暴露状态供 UI/API 读取，
// 不自动重建容器（决策：仅提示，由管理员手动执行 manager.sh update）。
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultUpdateURL      = "https://beacon.amuluze.com"
	defaultCheckInterval  = 6 * time.Hour
	minCheckInterval      = time.Hour
	versionCheckTimeout   = 10 * time.Second
	versionCheckUserAgent = "beacon-version-checker"
)

// BuildVersion 是当前 beacon 构建版本，由 main 通过 ldflags 注入（-X service.BuildVersion=...）。
// 缺省 dev；为 dev 时 beacon-help 无法判定版本序，update_available 恒为 false。
var BuildVersion = "dev"

// versionCheckEndpoint 来自 beacon-help 的版本清单接口路径（website/server/service/version.go）。
const versionCheckEndpoint = "/api/v1/version/latest"

// VersionChecker 周期性向 beacon-help 询问最新版本。
type VersionChecker struct {
	client  *http.Client
	url     string
	current string
	tick    time.Duration
	stopCh  chan struct{}
	stopped bool

	mu     sync.RWMutex
	latest UpdateStatus
}

// UpdateStatus 是一次版本检查的快照，供 API/UI 读取。
type UpdateStatus struct {
	UpdateAvailable  bool   `json:"update_available"`
	CurrentVersion   string `json:"current_version"`
	LatestVersion    string `json:"latest_version"`
	MinRequired      string `json:"min_required_version,omitempty"`
	ReleaseNotes     string `json:"release_notes,omitempty"`
	CheckedAt        string `json:"checked_at,omitempty"`
	LastError        string `json:"last_error,omitempty"`
}

// NewVersionChecker 构造检查器，并应用 Update 配置的默认值（未配置时启用、6h、官方 URL）。
func NewVersionChecker(cfg *Config) *VersionChecker {
	interval := time.Duration(cfg.Update.CheckInterval) * time.Second
	if interval <= 0 {
		interval = defaultCheckInterval
	}
	if interval < minCheckInterval {
		interval = minCheckInterval
	}
	url := strings.TrimRight(cfg.Update.URL, "/")
	if url == "" {
		url = defaultUpdateURL
	}
	return &VersionChecker{
		client:  &http.Client{Timeout: versionCheckTimeout},
		url:     url,
		current: BuildVersion,
		tick:    interval,
		stopCh:  make(chan struct{}),
	}
}

// Enabled 报告是否启用版本检查。
func (vc *VersionChecker) Enabled() bool { return vc != nil }

// Run 阻塞执行周期检查，直到 Stop 被调用。首次执行立即触发。
func (vc *VersionChecker) Run() {
	if vc == nil {
		return
	}
	ticker := time.NewTicker(vc.tick)
	defer ticker.Stop()
	vc.checkOnce()
	for {
		select {
		case <-vc.stopCh:
			return
		case <-ticker.C:
			vc.checkOnce()
		}
	}
}

// Stop 终止检查循环。
func (vc *VersionChecker) Stop() {
	if vc == nil {
		return
	}
	vc.mu.Lock()
	defer vc.mu.Unlock()
	if vc.stopped {
		return
	}
	vc.stopped = true
	close(vc.stopCh)
}

// Status 返回最近一次检查快照（线程安全）。
func (vc *VersionChecker) Status() UpdateStatus {
	if vc == nil {
		return UpdateStatus{CurrentVersion: BuildVersion}
	}
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.latest
}

func (vc *VersionChecker) setStatus(s UpdateStatus) {
	vc.mu.Lock()
	vc.latest = s
	vc.mu.Unlock()
}

// checkOnce 执行一次 HTTP 轮询并更新状态。
func (vc *VersionChecker) checkOnce() {
	status := UpdateStatus{CurrentVersion: vc.current, CheckedAt: time.Now().UTC().Format(time.RFC3339)}

	endpoint := vc.url + versionCheckEndpoint + "?current=" + vc.current
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		status.LastError = err.Error()
		vc.setStatus(status)
		slog.Debug("version check request build failed", "err", err)
		return
	}
	req.Header.Set("User-Agent", versionCheckUserAgent)

	resp, err := vc.client.Do(req)
	if err != nil {
		status.LastError = err.Error()
		vc.setStatus(status)
		slog.Debug("version check failed", "err", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		status.LastError = fmt.Sprintf("beacon-help returned %d", resp.StatusCode)
		vc.setStatus(status)
		slog.Debug("version check non-200", "status", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		status.LastError = err.Error()
		vc.setStatus(status)
		return
	}

	var manifest struct {
		LatestVersion      string `json:"latest_version"`
		MinRequiredVersion string `json:"min_required_version"`
		UpdateAvailable    bool   `json:"update_available"`
		ReleaseNotes       string `json:"release_notes"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		status.LastError = "parse manifest: " + err.Error()
		vc.setStatus(status)
		return
	}

	status.LatestVersion = manifest.LatestVersion
	status.MinRequired = manifest.MinRequiredVersion
	status.ReleaseNotes = manifest.ReleaseNotes
	status.UpdateAvailable = manifest.UpdateAvailable
	vc.setStatus(status)

	if status.UpdateAvailable {
		slog.Warn("beacon update available",
			"current", status.CurrentVersion,
			"latest", status.LatestVersion,
			"run", "manager.sh update")
	} else {
		slog.Info("version check ok", "current", status.CurrentVersion, "latest", status.LatestVersion)
	}
}

// globalVersionChecker 由 Init 注入，供 Router 的只读 API 端点读取，
// 避免改动 wire 生成器。nil 表示未启用。
var globalVersionChecker *VersionChecker

// SetGlobalVersionChecker 注册包级单例，供 API handler 读取状态。
func SetGlobalVersionChecker(vc *VersionChecker) {
	globalVersionChecker = vc
}

// GlobalUpdateStatus 返回全局检查器的状态快照；未配置时返回仅含当前版本的最小状态。
func GlobalUpdateStatus() UpdateStatus {
	if globalVersionChecker == nil {
		return UpdateStatus{CurrentVersion: BuildVersion}
	}
	return globalVersionChecker.Status()
}
