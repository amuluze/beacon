package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	stderrors "errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"beacon/pkg/contextx"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
	tunnelpkg "common/rpc/tunnel"

	"gorm.io/gorm"
)

const (
	defaultUpdateTimeout       = 5 * time.Minute
	defaultUpdateConcurrency   = 2
	defaultUpdateRetryInterval = 30 * time.Minute
	defaultPackageDir          = "/app/downloads/collia"
	maxAuditOperateRunes       = 255
	// maxFailureRetention 是 failedAt 条目的最长保留时长。超过后视为该 agent 长期不可达，
	// 从内存 map 中删除避免长期离线 agent 的 entry 永久驻留累积。
	// 远大于 RetryInterval(默认 30min)，不影响正常冷却逻辑。
	maxFailureRetention = 24 * time.Hour
)

var (
	ErrUpdaterUnavailable   = stderrors.New("collia updater is not configured")
	ErrUpgradeInProgress    = stderrors.New("collia upgrade already in progress")
	ErrUnsupportedArch      = stderrors.New("unsupported agent architecture")
	ErrDowngradeBlocked     = stderrors.New("collia downgrade is not allowed")
	ErrInvalidAgentVersion  = stderrors.New("invalid agent version")
	ErrInvalidTargetVersion = stderrors.New("invalid target collia version")
)

type UpdateConfig struct {
	Enable        bool
	TargetVersion string
	Timeout       time.Duration
	MaxConcurrent int
	RetryInterval time.Duration
	PublicBaseURL string
	InstallToken  string
	PackageDir    string
}

type UpdateResult struct {
	AgentID        string `json:"agent_id"`
	CurrentVersion string `json:"current_version"`
	TargetVersion  string `json:"target_version"`
	Arch           string `json:"arch,omitempty"`
	Triggered      bool   `json:"triggered"`
	Message        string `json:"message"`
}

type updateCaller interface {
	Call(ctx context.Context, method string, args interface{}, reply interface{}) error
}

type Updater struct {
	config UpdateConfig
	caller updateCaller
	db     *database.DB

	mu        sync.Mutex
	running   map[string]struct{}
	failedAt  map[string]time.Time
	semaphore chan struct{}
	now       func() time.Time
}

func NewUpdater(config UpdateConfig, caller updateCaller, db *database.DB) (*Updater, error) {
	config.TargetVersion = strings.TrimSpace(config.TargetVersion)
	if config.Timeout <= 0 {
		config.Timeout = defaultUpdateTimeout
	}
	if config.MaxConcurrent <= 0 {
		config.MaxConcurrent = defaultUpdateConcurrency
	}
	if config.RetryInterval <= 0 {
		config.RetryInterval = defaultUpdateRetryInterval
	}
	if config.PackageDir == "" {
		config.PackageDir = defaultPackageDir
	}
	config.PublicBaseURL = strings.TrimRight(strings.TrimSpace(config.PublicBaseURL), "/")

	if config.Enable {
		if caller == nil || db == nil {
			return nil, ErrUpdaterUnavailable
		}
		if _, err := parseSemVersion(config.TargetVersion); err != nil {
			return nil, fmt.Errorf("%w: %q", ErrInvalidTargetVersion, config.TargetVersion)
		}
		if config.PublicBaseURL == "" {
			return nil, fmt.Errorf("collia updater public base URL is empty")
		}
		if config.InstallToken == "" {
			return nil, fmt.Errorf("collia updater install token is empty")
		}
	}

	return &Updater{
		config:    config,
		caller:    caller,
		db:        db,
		running:   make(map[string]struct{}),
		failedAt:  make(map[string]time.Time),
		semaphore: make(chan struct{}, config.MaxConcurrent),
		now:       time.Now,
	}, nil
}

// Schedule performs a non-blocking automatic update decision after tunnel
// registration. The actual RPC always runs outside the lifecycle callback.
func (u *Updater) Schedule(info tunnelpkg.AgentInfo) {
	if u == nil || !u.config.Enable || info.AgentID == "" {
		return
	}
	u.purgeStaleFailures()
	if err := u.claim(info.AgentID, false); err != nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), u.config.Timeout)
		defer cancel()
		_, err := u.execute(ctx, info)
		u.finish(info.AgentID, err)
		if err != nil {
			slog.Error("automatic collia upgrade failed", "agent_id", info.AgentID, "err", err)
		}
	}()
}

// Trigger manually runs the same update path while bypassing automatic failure
// cooldown. Single-Agent deduplication and global concurrency still apply.
func (u *Updater) Trigger(ctx context.Context, info tunnelpkg.AgentInfo) (UpdateResult, error) {
	if u == nil {
		return UpdateResult{}, ErrUpdaterUnavailable
	}
	if err := u.validateManualConfig(); err != nil {
		return UpdateResult{}, err
	}
	if err := u.claim(info.AgentID, true); err != nil {
		return UpdateResult{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()
	result, err := u.execute(ctx, info)
	u.finish(info.AgentID, err)
	return result, err
}

func (u *Updater) validateManualConfig() error {
	if _, err := parseSemVersion(u.config.TargetVersion); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidTargetVersion, u.config.TargetVersion)
	}
	if u.config.PublicBaseURL == "" {
		return fmt.Errorf("collia updater public base URL is empty")
	}
	if u.config.InstallToken == "" {
		return fmt.Errorf("collia updater install token is empty")
	}
	if u.caller == nil || u.db == nil {
		return ErrUpdaterUnavailable
	}
	return nil
}

func (u *Updater) claim(agentID string, manual bool) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.running[agentID]; ok {
		return ErrUpgradeInProgress
	}
	if !manual {
		if failedAt, ok := u.failedAt[agentID]; ok && u.now().Before(failedAt.Add(u.config.RetryInterval)) {
			return fmt.Errorf("automatic collia upgrade is cooling down")
		}
	}
	u.running[agentID] = struct{}{}
	return nil
}

func (u *Updater) finish(agentID string, err error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.running, agentID)
	if err != nil {
		u.failedAt[agentID] = u.now()
		return
	}
	delete(u.failedAt, agentID)
}

// purgeStaleFailures 删除超过 maxFailureRetention 的失败记录，防止长期离线 agent
// 的 entry 在内存 map 中永久累积。在 agent 注册（Schedule）时顺带清理。
func (u *Updater) purgeStaleFailures() {
	u.mu.Lock()
	defer u.mu.Unlock()
	cutoff := u.now().Add(-maxFailureRetention)
	for agentID, failedAt := range u.failedAt {
		if failedAt.Before(cutoff) {
			delete(u.failedAt, agentID)
		}
	}
}

func (u *Updater) execute(ctx context.Context, info tunnelpkg.AgentInfo) (UpdateResult, error) {
	result := UpdateResult{
		AgentID:        info.AgentID,
		CurrentVersion: strings.TrimSpace(info.Version),
		TargetVersion:  u.config.TargetVersion,
	}

	comparison, err := compareSemVersions(info.Version, u.config.TargetVersion)
	if err != nil {
		err = fmt.Errorf("%w: %q", ErrInvalidAgentVersion, info.Version)
		return result, u.auditPreflightFailure(info.AgentID, err)
	}
	if comparison == 0 {
		result.Message = "collia is already up to date"
		return result, nil
	}
	if comparison > 0 {
		result.Message = "agent version is newer than target"
		err = fmt.Errorf("%w: current=%s target=%s", ErrDowngradeBlocked, info.Version, u.config.TargetVersion)
		if auditErr := u.writeAudit(info.AgentID, "Collia upgrade skipped: "+err.Error()); auditErr != nil {
			return result, fmt.Errorf("%w; write audit: %v", err, auditErr)
		}
		return result, err
	}

	if info.OS != "" && !strings.EqualFold(strings.TrimSpace(info.OS), "linux") {
		err = fmt.Errorf("unsupported agent OS %q", info.OS)
		return result, u.auditPreflightFailure(info.AgentID, err)
	}

	arch, err := u.resolveArch(info.AgentID, info.Arch)
	if err != nil {
		return result, u.auditPreflightFailure(info.AgentID, err)
	}
	result.Arch = arch

	select {
	case u.semaphore <- struct{}{}:
		defer func() { <-u.semaphore }()
	case <-ctx.Done():
		return result, u.auditPreflightFailure(info.AgentID, fmt.Errorf("wait for upgrade slot: %w", ctx.Err()))
	}

	checksum, err := sha256File(filepath.Join(u.config.PackageDir, arch, "collia"))
	if err != nil {
		return result, u.auditPreflightFailure(info.AgentID, fmt.Errorf("calculate collia sha256: %w", err))
	}
	downloadURL, err := u.downloadURL(arch)
	if err != nil {
		return result, u.auditPreflightFailure(info.AgentID, err)
	}

	startMessage := fmt.Sprintf("Collia upgrade started: current=%s target=%s arch=%s", displayVersion(info.Version), u.config.TargetVersion, arch)
	if err := u.writeAudit(info.AgentID, startMessage); err != nil {
		return result, fmt.Errorf("write collia upgrade start audit: %w", err)
	}

	args := rpcSchema.UpgradeAgentArgs{
		DownloadURL:  downloadURL,
		SHA256:       checksum,
		Version:      u.config.TargetVersion,
		InstallToken: u.config.InstallToken,
	}
	var reply rpcSchema.UpgradeAgentReply
	callCtx := contextx.NewAgentID(ctx, info.AgentID)
	if err := u.caller.Call(callCtx, "UpgradeAgent", args, &reply); err != nil {
		return result, u.auditRPCFailure(info.AgentID, fmt.Errorf("upgrade RPC failed: %w", err))
	}
	if !reply.Success {
		message := strings.TrimSpace(reply.Error)
		if message == "" {
			message = "agent rejected upgrade"
		}
		return result, u.auditRPCFailure(info.AgentID, stderrors.New(message))
	}

	result.Triggered = true
	result.Message = "collia upgrade triggered"
	if err := u.writeAudit(info.AgentID, fmt.Sprintf("Collia upgrade succeeded: target=%s arch=%s", u.config.TargetVersion, arch)); err != nil {
		return result, fmt.Errorf("upgrade succeeded but write success audit failed: %w", err)
	}
	return result, nil
}

func (u *Updater) auditPreflightFailure(agentID string, cause error) error {
	if auditErr := u.writeAudit(agentID, "Collia upgrade skipped: "+cause.Error()); auditErr != nil {
		return fmt.Errorf("%w; write audit: %v", cause, auditErr)
	}
	return cause
}

func (u *Updater) auditRPCFailure(agentID string, cause error) error {
	if auditErr := u.writeAudit(agentID, "Collia upgrade failed: "+cause.Error()); auditErr != nil {
		return fmt.Errorf("%w; write failure audit: %v", cause, auditErr)
	}
	return cause
}

func (u *Updater) writeAudit(agentID, operation string) error {
	if u.db == nil {
		return ErrUpdaterUnavailable
	}
	return u.db.Create(&model.Audit{
		Username: "system",
		AgentID:  agentID,
		Operate:  truncateRunes(operation, maxAuditOperateRunes),
	}).Error
}

func (u *Updater) resolveArch(agentID, reported string) (string, error) {
	if strings.TrimSpace(reported) != "" {
		if arch, ok := normalizeArch(reported); ok {
			return arch, nil
		}
		return "", fmt.Errorf("%w: %q", ErrUnsupportedArch, reported)
	}

	var host model.MonitorHost
	if err := u.db.Where("agent_id = ?", agentID).Order("timestamp desc").First(&host).Error; err == nil {
		if arch, ok := normalizeArch(host.KernelArch); ok {
			return arch, nil
		}
	} else if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("query latest host architecture: %w", err)
	}

	var docker model.MonitorDocker
	if err := u.db.Where("agent_id = ?", agentID).Order("timestamp desc").First(&docker).Error; err == nil {
		if arch, ok := normalizeArch(docker.Arch); ok {
			return arch, nil
		}
	} else if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("query latest docker architecture: %w", err)
	}

	return "", fmt.Errorf("%w: no recognized registration or monitoring value", ErrUnsupportedArch)
}

func normalizeArch(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amd64", "x86_64":
		return "amd64", true
	case "arm64", "aarch64":
		return "arm64", true
	default:
		return "", false
	}
}

func (u *Updater) downloadURL(arch string) (string, error) {
	base, err := url.Parse(u.config.PublicBaseURL)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("invalid collia public base URL")
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/api/v1/host/install/package"
	query := base.Query()
	query.Set("arch", arch)
	base.RawQuery = query.Encode()
	return base.String(), nil
}

func sha256File(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304 -- path is composed from configured package root and allowlisted arch.
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func truncateRunes(value string, max int) string {
	if max <= 0 || utf8.RuneCountInString(value) <= max {
		return value
	}
	return string([]rune(value)[:max])
}

func displayVersion(version string) string {
	if strings.TrimSpace(version) == "" {
		return "unknown"
	}
	return strings.TrimSpace(version)
}

type semVersion struct {
	major int
	minor int
	patch int
	pre   []string
}

func parseSemVersion(value string) (semVersion, error) {
	value = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(value), "v"))
	if value == "" {
		return semVersion{}, fmt.Errorf("empty version")
	}
	value = strings.SplitN(value, "+", 2)[0]
	parts := strings.SplitN(value, "-", 2)
	core := strings.Split(parts[0], ".")
	if len(core) != 3 {
		return semVersion{}, fmt.Errorf("version must contain major.minor.patch")
	}
	parsed := semVersion{}
	fields := []*int{&parsed.major, &parsed.minor, &parsed.patch}
	for index, component := range core {
		if component == "" || (len(component) > 1 && component[0] == '0') {
			return semVersion{}, fmt.Errorf("invalid numeric component")
		}
		n, err := strconv.Atoi(component)
		if err != nil || n < 0 {
			return semVersion{}, fmt.Errorf("invalid numeric component")
		}
		*fields[index] = n
	}
	if len(parts) == 2 {
		if parts[1] == "" {
			return semVersion{}, fmt.Errorf("empty prerelease")
		}
		parsed.pre = strings.Split(parts[1], ".")
		for _, identifier := range parsed.pre {
			if identifier == "" {
				return semVersion{}, fmt.Errorf("empty prerelease identifier")
			}
		}
	}
	return parsed, nil
}

func compareSemVersions(current, target string) (int, error) {
	if strings.TrimSpace(current) == "" {
		if _, err := parseSemVersion(target); err != nil {
			return 0, err
		}
		return -1, nil
	}
	left, err := parseSemVersion(current)
	if err != nil {
		return 0, err
	}
	right, err := parseSemVersion(target)
	if err != nil {
		return 0, err
	}
	return left.compare(right), nil
}

func (v semVersion) compare(other semVersion) int {
	for _, pair := range [][2]int{{v.major, other.major}, {v.minor, other.minor}, {v.patch, other.patch}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	if len(v.pre) == 0 && len(other.pre) == 0 {
		return 0
	}
	if len(v.pre) == 0 {
		return 1
	}
	if len(other.pre) == 0 {
		return -1
	}
	for index := 0; index < len(v.pre) && index < len(other.pre); index++ {
		if comparison := comparePrereleaseIdentifier(v.pre[index], other.pre[index]); comparison != 0 {
			return comparison
		}
	}
	if len(v.pre) < len(other.pre) {
		return -1
	}
	if len(v.pre) > len(other.pre) {
		return 1
	}
	return 0
}

func comparePrereleaseIdentifier(left, right string) int {
	leftNumber, leftErr := strconv.Atoi(left)
	rightNumber, rightErr := strconv.Atoi(right)
	if leftErr == nil && rightErr == nil {
		if leftNumber < rightNumber {
			return -1
		}
		if leftNumber > rightNumber {
			return 1
		}
		return 0
	}
	if leftErr == nil {
		return -1
	}
	if rightErr == nil {
		return 1
	}
	return strings.Compare(left, right)
}
