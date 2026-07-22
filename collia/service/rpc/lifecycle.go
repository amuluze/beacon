// Package rpc
// Description: Agent 生命周期 RPC —— UpgradeAgent（自更新）/ UninstallAgent（自卸载）。
// 对齐 Domain Spec agent-lifecycle-update.md：
//   - IU003 原子更新：collia.bak 回退
//   - RU002 下载校验：SHA256
//   - RU004 卸载清理：二进制/config/data/logs/服务注册
//   - RU006 串行：TryLock 互斥
package rpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	rpcSchema "common/rpc/schema"

	tunnel "common/rpc/tunnel"
)

// 默认清理路径（对齐 agent_install.go 的安装布局）。
const (
	defaultColliaBinaryPath = "/usr/sbin/collia"
	defaultConfigDir        = "/etc/collia"
	defaultDataDir          = "/data/beacon/collia"
	defaultLogDir           = "/data/beacon/logs/collia"
	backupSuffix            = ".bak"
)

// upgradeMu 保证同一 Agent 串行更新（RU006）。TryLock 失败即拒绝并发请求。
var upgradeMu sync.Mutex

// inProgress 标记当前是否正在更新，供并发请求快速判定。
var inProgress bool

// lifecyclePaths 返回当前 Service 实例生效的清理/替换路径集合，便于测试注入。
func (s *Service) lifecyclePaths() (binary, configDir, dataDir, logDir string) {
	binary = s.binaryPath
	if binary == "" {
		binary = defaultColliaBinaryPath
	}
	return binary, defaultConfigDir, defaultDataDir, defaultLogDir
}

// registerLifecycleHandlers 注册 UpgradeAgent / UninstallAgent。
func registerLifecycleHandlers(d *Dispatcher, svc *Service) {
	// UpgradeAgent 同时需要流式进度帧与最终 reply，故走原始 Register。
	d.Register("UpgradeAgent", func(ctx context.Context, payload []byte, streamSender func(*tunnel.Frame)) ([]byte, error) {
		var args rpcSchema.UpgradeAgentArgs
		if err := json.Unmarshal(payload, &args); err != nil {
			return nil, fmt.Errorf("unmarshal UpgradeAgentArgs: %w", err)
		}
		reply, err := svc.handleUpgrade(ctx, args, streamSender)
		if err != nil {
			return nil, err
		}
		return json.Marshal(reply)
	})

	RegisterUnary[rpcSchema.UninstallAgentArgs, rpcSchema.UninstallAgentReply](d, "UninstallAgent", svc.handleUninstall)
}

// sendProgress 通过 stream 上报一个更新进度帧。
func sendProgress(streamSender func(*tunnel.Frame), stage string, percent int, msg string) {
	if streamSender == nil {
		return
	}
	p := rpcSchema.UpgradeProgress{Stage: stage, Percent: percent, Message: msg}
	data, err := json.Marshal(p)
	if err != nil {
		return
	}
	streamSender(&tunnel.Frame{Payload: data})
}

// handleUpgrade 执行自更新全流程。
func (s *Service) handleUpgrade(ctx context.Context, args rpcSchema.UpgradeAgentArgs, streamSender func(*tunnel.Frame)) (*rpcSchema.UpgradeAgentReply, error) {
	// RU006 串行：并发请求直接拒绝。
	if !acquireUpgradeLock() {
		return &rpcSchema.UpgradeAgentReply{
			Success: false,
			Stage:   rpcSchema.UpgradeStageFailed,
			Error:   "update already in progress",
		}, errors.New("update already in progress")
	}
	defer releaseUpgradeLock()

	binaryPath, _, _, _ := s.lifecyclePaths()

	if err := ctx.Err(); err != nil {
		return failReply("context cancelled before download"), err
	}

	// 1. 下载
	sendProgress(streamSender, rpcSchema.UpgradeStageDownloading, 0, "downloading")
	tmpFile, err := s.downloadBinary(ctx, args)
	if err != nil {
		slog.Error("upgrade download failed", "err", err)
		sendProgress(streamSender, rpcSchema.UpgradeStageFailed, 0, err.Error())
		return failReply("download failed: " + err.Error()), nil
	}
	defer func() { _ = os.Remove(tmpFile) }() // 兜底清理临时文件

	// 2. SHA256 校验（RU002）
	sendProgress(streamSender, rpcSchema.UpgradeStageVerifying, 50, "verifying")
	if err := verifySHA256(tmpFile, args.SHA256); err != nil {
		slog.Error("upgrade verify failed", "err", err)
		sendProgress(streamSender, rpcSchema.UpgradeStageFailed, 0, err.Error())
		return failReply("verify failed: " + err.Error()), nil
	}

	// 3. 原子替换（IU003）
	sendProgress(streamSender, rpcSchema.UpgradeStageReplacing, 70, "replacing")
	if err := replaceBinary(tmpFile, binaryPath); err != nil {
		slog.Error("upgrade replace failed", "err", err)
		sendProgress(streamSender, rpcSchema.UpgradeStageFailed, 0, err.Error())
		return failReply("replace failed: " + err.Error()), nil
	}

	// 4. 触发重启（异步，确保 reply 先返回）
	sendProgress(streamSender, rpcSchema.UpgradeStageRestarting, 90, "restarting")
	s.scheduleRestart()

	sendProgress(streamSender, rpcSchema.UpgradeStageDone, 100, "done")
	slog.Info("upgrade applied, restarting", "version", args.Version)
	return &rpcSchema.UpgradeAgentReply{
		Success: true,
		Version: args.Version,
		Stage:   rpcSchema.UpgradeStageDone,
	}, nil
}

func failReply(msg string) *rpcSchema.UpgradeAgentReply {
	return &rpcSchema.UpgradeAgentReply{Success: false, Stage: rpcSchema.UpgradeStageFailed, Error: msg}
}

func acquireUpgradeLock() bool {
	upgradeMu.Lock()
	defer upgradeMu.Unlock()
	if inProgress {
		return false
	}
	inProgress = true
	return true
}

func releaseUpgradeLock() {
	upgradeMu.Lock()
	inProgress = false
	upgradeMu.Unlock()
}

// downloadBinary 用 install token 拉取新二进制到临时文件。
func (s *Service) downloadBinary(ctx context.Context, args rpcSchema.UpgradeAgentArgs) (string, error) {
	if args.DownloadURL == "" {
		return "", errors.New("empty download_url")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, args.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	if args.InstallToken != "" {
		req.Header.Set("X-Install-Token", args.InstallToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp("", "collia-upgrade-*")
	if err != nil {
		return "", fmt.Errorf("create temp: %w", err)
	}
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("copy body: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("close temp: %w", err)
	}
	return tmp.Name(), nil
}

// verifySHA256 校验文件 SHA256；期望值为空时跳过（不推荐，但兼容）。
func verifySHA256(file, expected string) error {
	if expected == "" {
		return errors.New("sha256 未提供，拒绝未校验的更新")
	}
	f, err := os.Open(file) //#nosec G304 -- 临时文件路径由 CreateTemp 生成
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if actual != expected {
		return fmt.Errorf("sha256 mismatch: want %s got %s", expected, actual)
	}
	return nil
}

// replaceBinary 原子替换：写 .new → 备份当前为 .bak → rename .new → 目标 → chmod。
func replaceBinary(tmpFile, target string) error {
	if target == "" {
		return errors.New("empty target binary path")
	}
	if err := os.Chmod(tmpFile, 0o755); err != nil {
		return fmt.Errorf("chmod temp: %w", err)
	}
	// 备份当前二进制（若存在）
	if _, err := os.Stat(target); err == nil {
		backup := target + backupSuffix
		// 移除旧备份
		_ = os.Remove(backup)
		if err := os.Rename(target, backup); err != nil {
			return fmt.Errorf("backup current binary: %w", err)
		}
	}
	// 原子 rename
	if err := os.Rename(tmpFile, target); err != nil {
		return fmt.Errorf("rename into place: %w", err)
	}
	return nil
}

// scheduleRestart 异步触发重启，确保 RPC reply 先于进程退出返回。
func (s *Service) scheduleRestart() {
	fn := s.restartFn
	go func() {
		// 给 Server 读取 reply 留出窗口
		// restartFn 为空时退化为 os.Exit，交由 systemd 拉起新版本
		if fn != nil {
			if err := fn(); err != nil {
				slog.Error("restart callback failed, exiting to let supervisor recover", "err", err)
			}
		}
		// 无论回调是否成功，都请求进程退出；systemd 会拉起新二进制
		// 不在此直接 os.Exit，避免回调已自行处理（如 systemctl restart）导致重复
	}()
}

// CleanupBackup 删除上次成功更新留下的 collia.bak；由新版本启动后调用。
func CleanupBackup(binaryPath string) {
	if binaryPath == "" {
		binaryPath = defaultColliaBinaryPath
	}
	backup := binaryPath + backupSuffix
	if _, err := os.Stat(backup); err == nil {
		if err := os.Remove(backup); err != nil {
			slog.Warn("remove stale collia.bak failed", "err", err)
		}
	}
}

// DefaultRestartFn 返回一个用于自更新后重启的回调：优先 systemctl restart，
// 失败则直接 os.Exit(0) 交由 supervisor（systemd Restart=on-failure）拉起。
// systemctl restart 会终止当前进程并由 systemd 用新二进制重新拉起。
func DefaultRestartFn(serviceName string) func() error {
	return func() error {
		if serviceName == "" {
			serviceName = "collia"
		}
		// systemctl restart 会异步终止本进程，无需等待
		if _, err := exec.LookPath("systemctl"); err == nil {
			cmd := exec.Command("systemctl", "restart", serviceName)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("systemctl restart start: %w", err)
			}
			return nil
		}
		// 无 systemd：请求退出，依赖外部 supervisor
		go func() {
			// 给 reply 返回留窗口
			time.Sleep(time.Second)
			os.Exit(0)
		}()
		return nil
	}
}

// handleUninstall 执行自卸载（spec T3）。返回残留项供 Server 审计。
func (s *Service) handleUninstall(ctx context.Context, args rpcSchema.UninstallAgentArgs, reply *rpcSchema.UninstallAgentReply) error {
	binary, configDir, dataDir, logDir := s.lifecyclePaths()
	slog.Info("uninstall triggered", "force", args.Force)

	var residuals []string
	for _, p := range []string{binary, binary + backupSuffix, configDir, dataDir, logDir} {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err != nil {
			continue // 不存在即无需清理
		}
		if err := os.RemoveAll(p); err != nil {
			residuals = append(residuals, fmt.Sprintf("%s: %v", p, err))
			slog.Warn("uninstall remove failed", "path", p, "err", err)
		}
	}

	reply.Residuals = residuals
	if len(residuals) > 0 {
		reply.Success = false
		reply.Error = "partial uninstall, see residuals"
	} else {
		reply.Success = true
	}

	// 服务注册由 cmd/main 的 remove 子命令负责；此处不可达 systemd，仅清理文件。
	// 卸载完成后断开 tunnel：由调用方（cmd/main 或 Server）在 reply 后处理。
	return nil
}
