package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"beacon/pkg/contextx"
	"beacon/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
	tunnelpkg "common/rpc/tunnel"
)

type capturedUpgradeCall struct {
	agentID string
	method  string
	args    rpcSchema.UpgradeAgentArgs
}

type fakeUpdateCaller struct {
	mu      sync.Mutex
	calls   []capturedUpgradeCall
	callErr error
	reply   rpcSchema.UpgradeAgentReply
	entered chan struct{}
	release chan struct{}
}

func (f *fakeUpdateCaller) Call(ctx context.Context, method string, args interface{}, reply interface{}) error {
	upgradeArgs := args.(rpcSchema.UpgradeAgentArgs)
	f.mu.Lock()
	f.calls = append(f.calls, capturedUpgradeCall{
		agentID: contextx.FromAgentID(ctx),
		method:  method,
		args:    upgradeArgs,
	})
	f.mu.Unlock()
	if f.entered != nil {
		select {
		case f.entered <- struct{}{}:
		default:
		}
	}
	if f.release != nil {
		select {
		case <-f.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if f.callErr != nil {
		return f.callErr
	}
	result := f.reply
	if !result.Success && result.Error == "" {
		result.Success = true
		result.Stage = rpcSchema.UpgradeStageDone
		result.Version = upgradeArgs.Version
	}
	*(reply.(*rpcSchema.UpgradeAgentReply)) = result
	return nil
}

func (f *fakeUpdateCaller) snapshot() []capturedUpgradeCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]capturedUpgradeCall(nil), f.calls...)
}

func newUpdaterTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "updater")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.Agent{}, &model.Audit{}, &model.MonitorHost{}, &model.MonitorDocker{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newUpdaterTestConfig(t *testing.T, packageContents string) UpdateConfig {
	t.Helper()
	root := t.TempDir()
	for _, arch := range []string{"amd64", "arm64"} {
		dir := filepath.Join(root, arch)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatalf("mkdir package dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "collia"), []byte(packageContents+arch), 0o700); err != nil {
			t.Fatalf("write package: %v", err)
		}
	}
	return UpdateConfig{
		Enable:        true,
		TargetVersion: "v3.0.6",
		Timeout:       2 * time.Second,
		MaxConcurrent: 2,
		RetryInterval: time.Hour,
		PublicBaseURL: "https://beacon.example.com",
		InstallToken:  "secret-install-token",
		PackageDir:    root,
	}
}

func TestUpdaterTriggerComputesChecksumCallsAgentAndAudits(t *testing.T) {
	db := newUpdaterTestDB(t)
	config := newUpdaterTestConfig(t, "binary-")
	caller := &fakeUpdateCaller{}
	updater, err := NewUpdater(config, caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	result, err := updater.Trigger(context.Background(), tunnelpkg.AgentInfo{
		AgentID: "agent-1",
		Version: "v3.0.5",
		OS:      "linux",
		Arch:    "x86_64",
	})
	if err != nil {
		t.Fatalf("Trigger: %v", err)
	}
	if !result.Triggered || result.Arch != "amd64" {
		t.Fatalf("result = %#v", result)
	}

	calls := caller.snapshot()
	if len(calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(calls))
	}
	call := calls[0]
	if call.agentID != "agent-1" || call.method != "UpgradeAgent" {
		t.Fatalf("call target = %q method = %q", call.agentID, call.method)
	}
	if call.args.DownloadURL != "https://beacon.example.com/api/v1/host/install/package?arch=amd64" {
		t.Fatalf("DownloadURL = %q", call.args.DownloadURL)
	}
	if call.args.InstallToken != config.InstallToken || call.args.Version != config.TargetVersion {
		t.Fatalf("args = %#v", call.args)
	}
	wantHash := sha256.Sum256([]byte("binary-amd64"))
	if call.args.SHA256 != hex.EncodeToString(wantHash[:]) {
		t.Fatalf("SHA256 = %q, want %q", call.args.SHA256, hex.EncodeToString(wantHash[:]))
	}

	var audits []model.Audit
	if err := db.Order("id asc").Find(&audits).Error; err != nil {
		t.Fatalf("query audits: %v", err)
	}
	if len(audits) != 2 {
		t.Fatalf("audit count = %d, want start and success", len(audits))
	}
	for _, audit := range audits {
		if audit.Username != "system" || audit.AgentID != "agent-1" {
			t.Fatalf("audit = %#v", audit)
		}
		if strings.Contains(audit.Operate, config.InstallToken) {
			t.Fatal("audit leaked install token")
		}
	}
}

func TestUpdaterVersionDecisions(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		wantCalls  int
		wantAudits int64
		wantErr    error
	}{
		{name: "empty legacy version upgrades", version: "", wantCalls: 1, wantAudits: 2},
		{name: "same version is quiet", version: "v3.0.6", wantCalls: 0, wantAudits: 0},
		{name: "newer version never downgrades", version: "v3.1.0", wantCalls: 0, wantAudits: 1, wantErr: ErrDowngradeBlocked},
		{name: "invalid version is not guessed", version: "dev", wantCalls: 0, wantAudits: 1, wantErr: ErrInvalidAgentVersion},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newUpdaterTestDB(t)
			caller := &fakeUpdateCaller{}
			updater, err := NewUpdater(newUpdaterTestConfig(t, "version-test-"), caller, db)
			if err != nil {
				t.Fatalf("NewUpdater: %v", err)
			}
			_, gotErr := updater.Trigger(context.Background(), tunnelpkg.AgentInfo{
				AgentID: "agent-version",
				Version: tc.version,
				OS:      "linux",
				Arch:    "amd64",
			})
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("Trigger: %v", gotErr)
			}
			if tc.wantErr != nil && !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("Trigger error = %v, want %v", gotErr, tc.wantErr)
			}
			if got := len(caller.snapshot()); got != tc.wantCalls {
				t.Fatalf("calls = %d, want %d", got, tc.wantCalls)
			}
			var count int64
			if err := db.Model(&model.Audit{}).Count(&count).Error; err != nil {
				t.Fatalf("count audits: %v", err)
			}
			if count != tc.wantAudits {
				t.Fatalf("audits = %d, want %d", count, tc.wantAudits)
			}
		})
	}
}

func TestUpdaterInfersMissingArchitectureFromLatestHostReport(t *testing.T) {
	db := newUpdaterTestDB(t)
	if err := db.Create(&model.MonitorHost{
		AgentID: "legacy-agent", Timestamp: time.Now(), KernelArch: "aarch64",
	}).Error; err != nil {
		t.Fatalf("create host report: %v", err)
	}
	caller := &fakeUpdateCaller{}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "arch-test-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	result, err := updater.Trigger(context.Background(), tunnelpkg.AgentInfo{AgentID: "legacy-agent"})
	if err != nil {
		t.Fatalf("Trigger: %v", err)
	}
	if result.Arch != "arm64" {
		t.Fatalf("Arch = %q, want arm64", result.Arch)
	}
	if got := caller.snapshot()[0].args.DownloadURL; !strings.HasSuffix(got, "arch=arm64") {
		t.Fatalf("DownloadURL = %q", got)
	}
}

func TestUpdaterUnknownArchitectureAuditsAndSkips(t *testing.T) {
	db := newUpdaterTestDB(t)
	caller := &fakeUpdateCaller{}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "unknown-arch-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	_, err = updater.Trigger(context.Background(), tunnelpkg.AgentInfo{
		AgentID: "agent-unknown", Version: "v3.0.5", OS: "linux", Arch: "riscv64",
	})
	if !errors.Is(err, ErrUnsupportedArch) {
		t.Fatalf("Trigger error = %v, want ErrUnsupportedArch", err)
	}
	if len(caller.snapshot()) != 0 {
		t.Fatal("unknown architecture must not call UpgradeAgent")
	}
	var audit model.Audit
	if err := db.Last(&audit).Error; err != nil {
		t.Fatalf("query audit: %v", err)
	}
	if audit.AgentID != "agent-unknown" || !strings.Contains(audit.Operate, "unsupported") {
		t.Fatalf("audit = %#v", audit)
	}
}

func TestUpdaterStartAuditFailurePreventsRPC(t *testing.T) {
	db := newUpdaterTestDB(t)
	if err := db.Migrator().DropTable(&model.Audit{}); err != nil {
		t.Fatalf("drop audit table: %v", err)
	}
	caller := &fakeUpdateCaller{}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "audit-fail-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	_, err = updater.Trigger(context.Background(), tunnelpkg.AgentInfo{
		AgentID: "agent-audit", Version: "v3.0.5", OS: "linux", Arch: "amd64",
	})
	if err == nil || !strings.Contains(err.Error(), "start audit") {
		t.Fatalf("Trigger error = %v", err)
	}
	if len(caller.snapshot()) != 0 {
		t.Fatal("RPC must not run when start audit cannot be written")
	}
}

func TestUpdaterRPCFailureWritesTruncatedFailureAudit(t *testing.T) {
	db := newUpdaterTestDB(t)
	caller := &fakeUpdateCaller{callErr: errors.New(strings.Repeat("x", 400))}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "rpc-fail-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	_, err = updater.Trigger(context.Background(), tunnelpkg.AgentInfo{
		AgentID: "agent-rpc", Version: "v3.0.5", OS: "linux", Arch: "amd64",
	})
	if err == nil {
		t.Fatal("Trigger should return RPC error")
	}
	var audits []model.Audit
	if err := db.Order("id asc").Find(&audits).Error; err != nil {
		t.Fatalf("query audits: %v", err)
	}
	if len(audits) != 2 || !strings.Contains(audits[1].Operate, "failed") {
		t.Fatalf("audits = %#v", audits)
	}
	if got := len([]rune(audits[1].Operate)); got > maxAuditOperateRunes {
		t.Fatalf("failure audit length = %d", got)
	}
}

func TestUpdaterAutomaticScheduleDeduplicatesAndCoolsDownFailures(t *testing.T) {
	db := newUpdaterTestDB(t)
	caller := &fakeUpdateCaller{
		callErr: errors.New("offline"),
		entered: make(chan struct{}, 2),
		release: make(chan struct{}),
	}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "schedule-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	info := tunnelpkg.AgentInfo{AgentID: "agent-auto", Version: "v3.0.5", OS: "linux", Arch: "amd64"}

	updater.Schedule(info)
	updater.Schedule(info)
	select {
	case <-caller.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("automatic call did not start")
	}
	if got := len(caller.snapshot()); got != 1 {
		t.Fatalf("calls while running = %d, want 1", got)
	}
	close(caller.release)
	waitForUpdater(t, func() bool {
		updater.mu.Lock()
		defer updater.mu.Unlock()
		_, running := updater.running[info.AgentID]
		_, cooling := updater.failedAt[info.AgentID]
		return !running && cooling
	})

	updater.Schedule(info)
	time.Sleep(50 * time.Millisecond)
	if got := len(caller.snapshot()); got != 1 {
		t.Fatalf("calls during cooldown = %d, want 1", got)
	}

	caller.callErr = nil
	if _, err := updater.Trigger(context.Background(), info); err != nil {
		t.Fatalf("manual Trigger should bypass cooldown: %v", err)
	}
	if got := len(caller.snapshot()); got != 2 {
		t.Fatalf("calls after manual retry = %d, want 2", got)
	}
}

func TestUpdaterGlobalConcurrencyLimit(t *testing.T) {
	db := newUpdaterTestDB(t)
	caller := &fakeUpdateCaller{
		entered: make(chan struct{}, 2),
		release: make(chan struct{}),
	}
	config := newUpdaterTestConfig(t, "concurrency-")
	config.MaxConcurrent = 1
	updater, err := NewUpdater(config, caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	updater.Schedule(tunnelpkg.AgentInfo{AgentID: "agent-a", Version: "v3.0.5", OS: "linux", Arch: "amd64"})
	updater.Schedule(tunnelpkg.AgentInfo{AgentID: "agent-b", Version: "v3.0.5", OS: "linux", Arch: "amd64"})
	select {
	case <-caller.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first automatic call did not start")
	}
	time.Sleep(50 * time.Millisecond)
	if got := len(caller.snapshot()); got != 1 {
		t.Fatalf("concurrent calls = %d, want 1", got)
	}

	close(caller.release)
	select {
	case <-caller.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("second automatic call did not start after slot release")
	}
	waitForUpdater(t, func() bool { return len(caller.snapshot()) == 2 })
}

func TestServiceOnAgentConnectPersistsMetadataAndSchedulesUpdate(t *testing.T) {
	db := newUpdaterTestDB(t)
	caller := &fakeUpdateCaller{entered: make(chan struct{}, 1)}
	updater, err := NewUpdater(newUpdaterTestConfig(t, "connect-"), caller, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	service := NewAgentService(NewAgentRepo(db))
	service.SetUpdater(updater)

	service.OnAgentConnect(tunnelpkg.AgentInfo{
		AgentID: "agent-connect", Version: "v3.0.5", OS: "linux", Arch: "amd64",
	})
	select {
	case <-caller.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("OnAgentConnect did not asynchronously schedule update")
	}
	var stored model.Agent
	if err := db.Where("agent_id = ?", "agent-connect").First(&stored).Error; err != nil {
		t.Fatalf("query agent: %v", err)
	}
	if stored.Version != "v3.0.5" || stored.OS != "linux" || stored.Arch != "amd64" {
		t.Fatalf("stored metadata = %#v", stored)
	}
}

func TestCompareSemVersions(t *testing.T) {
	tests := []struct {
		current string
		target  string
		want    int
	}{
		{"v3.0.5", "v3.0.6", -1},
		{"3.0.6", "v3.0.6", 0},
		{"v3.1.0", "v3.0.6", 1},
		{"v3.0.6-rc.1", "v3.0.6", -1},
		{"", "v3.0.6", -1},
	}
	for _, tc := range tests {
		got, err := compareSemVersions(tc.current, tc.target)
		if err != nil {
			t.Fatalf("compareSemVersions(%q, %q): %v", tc.current, tc.target, err)
		}
		if got != tc.want {
			t.Fatalf("compareSemVersions(%q, %q) = %d, want %d", tc.current, tc.target, got, tc.want)
		}
	}
}

func waitForUpdater(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not reached before timeout")
}

// TestUpdaterPurgesStaleFailures 验证 Schedule 入口会清理超过 maxFailureRetention
// 的失败记录，防止长期离线 agent 的 entry 在内存 map 中永久累积。
func TestUpdaterPurgesStaleFailures(t *testing.T) {
	db := newUpdaterTestDB(t)
	updater, err := NewUpdater(newUpdaterTestConfig(t, "purge-"), &fakeUpdateCaller{}, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}

	now := time.Now()
	updater.now = func() time.Time { return now }

	// 直接注入 failedAt：一条已超期（>24h），一条仍在保留期内。
	updater.mu.Lock()
	updater.failedAt["stale-agent"] = now.Add(-25 * time.Hour)
	updater.failedAt["fresh-agent"] = now.Add(-time.Hour)
	updater.mu.Unlock()

	updater.purgeStaleFailures()

	updater.mu.Lock()
	defer updater.mu.Unlock()
	if _, ok := updater.failedAt["stale-agent"]; ok {
		t.Fatalf("stale-agent entry should have been purged (older than %v)", maxFailureRetention)
	}
	if _, ok := updater.failedAt["fresh-agent"]; !ok {
		t.Fatalf("fresh-agent entry should be retained (within %v)", maxFailureRetention)
	}
}

// TestUpdaterPurgeStaleFailuresKeepsAllWhenNoneExpired 验证所有 entry 都在保留期内时
// 不会误删。
func TestUpdaterPurgeStaleFailuresKeepsAllWhenNoneExpired(t *testing.T) {
	db := newUpdaterTestDB(t)
	updater, err := NewUpdater(newUpdaterTestConfig(t, "purge-keep-"), &fakeUpdateCaller{}, db)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	now := time.Now()
	updater.now = func() time.Time { return now }

	updater.mu.Lock()
	updater.failedAt["a"] = now.Add(-maxFailureRetention + time.Hour)
	updater.failedAt["b"] = now.Add(-time.Minute)
	updater.mu.Unlock()

	updater.purgeStaleFailures()

	updater.mu.Lock()
	defer updater.mu.Unlock()
	if len(updater.failedAt) != 2 {
		t.Fatalf("failedAt size = %d, want 2 (all within retention)", len(updater.failedAt))
	}
}
