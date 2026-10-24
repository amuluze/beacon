// Package schema
// Agent lifecycle RPC types: UpgradeAgent / UninstallAgent.
// 对齐 Domain Spec agent-lifecycle-update.md（IU003 原子更新、RU002 校验、RU004 清理、RU006 串行）。
package schema

// UpgradeAgentArgs 由 Server 经反向 tunnel 下发给 Agent，触发自更新。
// DownloadURL 复用现有 /api/v1/host/install/package 端点；InstallToken 为 agent 安装令牌。
type UpgradeAgentArgs struct {
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
	Version     string `json:"version"`
	InstallToken string `json:"install_token"`
}

// UpgradeAgentReply 返回更新结果。Stage 携带最终状态供 Server 审计。
type UpgradeAgentReply struct {
	Success bool   `json:"success"`
	Version string `json:"version"`
	Stage   string `json:"stage,omitempty"`
	Error   string `json:"error,omitempty"`
}

// UpgradeStage 标识更新进度帧的阶段。
const (
	UpgradeStageDownloading = "downloading"
	UpgradeStageVerifying   = "verifying"
	UpgradeStageReplacing   = "replacing"
	UpgradeStageRestarting  = "restarting"
	UpgradeStageDone        = "done"
	UpgradeStageFailed      = "failed"
)

// UpgradeProgress 更新过程的进度帧（通过 stream 帧上报）。
type UpgradeProgress struct {
	Stage   string `json:"stage"`
	Percent int    `json:"percent,omitempty"`
	Message string `json:"message,omitempty"`
}

// UninstallAgentArgs 触发 Agent 自卸载。Force=true 表示跳过二次确认（Server 侧已确认）。
type UninstallAgentArgs struct {
	Force bool `json:"force"`
}

// UninstallAgentReply 返回卸载结果；Residuals 列出无法删除的残留路径（T3-02）。
type UninstallAgentReply struct {
	Success   bool     `json:"success"`
	Error     string   `json:"error,omitempty"`
	Residuals []string `json:"residuals,omitempty"`
}
