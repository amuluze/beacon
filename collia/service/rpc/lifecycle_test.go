package rpc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	rpcSchema "common/rpc/schema"
	tunnel "common/rpc/tunnel"
)

// newTestService 构造一个不依赖 docker manager 的 Service，binaryPath 指向临时目录。
func newTestService(t *testing.T) (*Service, string) {
	t.Helper()
	dir := t.TempDir()
	binary := filepath.Join(dir, "collia")
	s := &Service{
		rootDir:    dir,
		binaryPath: binary,
	}
	return s, dir
}

func sha256Hex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// noopStream 返回一个不做事的 streamSender（满足 func(*tunnel.Frame)）。
func noopStream() func(*tunnel.Frame) { return func(*tunnel.Frame) {} }

func TestVerifySHA256RejectsEmpty(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "f")
	if err := os.WriteFile(tmp, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifySHA256(tmp, ""); err == nil {
		t.Fatal("expected error when sha256 empty")
	}
}

func TestVerifySHA256Mismatch(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "f")
	if err := os.WriteFile(tmp, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifySHA256(tmp, "deadbeef"); err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestReplaceBinaryAtomicWithBackup(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "collia")
	old := []byte("old-binary")
	if err := os.WriteFile(target, old, 0o755); err != nil {
		t.Fatal(err)
	}
	tmp := filepath.Join(dir, "collia.new")
	if err := os.WriteFile(tmp, []byte("new-binary"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := replaceBinary(tmp, target); err != nil {
		t.Fatalf("replace: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if !bytes.Equal(got, []byte("new-binary")) {
		t.Fatalf("target not replaced: %q", got)
	}
	bak, err := os.ReadFile(target + backupSuffix)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if !bytes.Equal(bak, old) {
		t.Fatalf("backup content mismatch: %q", bak)
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Fatalf("temp file should be gone: %v", err)
	}
}

func TestReplaceBinaryChmods0755(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "collia")
	if err := os.WriteFile(target, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	tmp := filepath.Join(dir, "new")
	if err := os.WriteFile(tmp, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := replaceBinary(tmp, target); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("perm = %o, want 755", info.Mode().Perm())
	}
}

func TestCleanupBackup(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "collia")
	backup := binary + backupSuffix
	if err := os.WriteFile(binary, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backup, []byte("bak"), 0o644); err != nil {
		t.Fatal(err)
	}
	CleanupBackup(binary)
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatalf("backup should be removed")
	}
}

func TestHandleUpgradeSuccess(t *testing.T) {
	s, dir := newTestService(t)
	binary := filepath.Join(dir, "collia")
	if err := os.WriteFile(binary, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	payload := []byte("brand-new-binary")
	var restartCalled atomic.Bool
	s.restartFn = func() error { restartCalled.Store(true); return nil }

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Install-Token") != "tok" {
			t.Errorf("missing install token header")
		}
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	args := rpcSchema.UpgradeAgentArgs{
		DownloadURL:  srv.URL,
		SHA256:       sha256Hex(payload),
		Version:      "v9.9.9",
		InstallToken: "tok",
	}
	reply, err := s.handleUpgrade(context.Background(), args, noopStream())
	if err != nil {
		t.Fatalf("handleUpgrade: %v", err)
	}
	if !reply.Success {
		t.Fatalf("expected success, got %+v", reply)
	}
	got, _ := os.ReadFile(binary)
	if !bytes.Equal(got, payload) {
		t.Fatalf("binary not replaced: %q", got)
	}
}

func TestHandleUpgradeSHA256MismatchKeepsOld(t *testing.T) {
	s, dir := newTestService(t)
	binary := filepath.Join(dir, "collia")
	old := []byte("old")
	if err := os.WriteFile(binary, old, 0o755); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("tampered"))
	}))
	defer srv.Close()

	args := rpcSchema.UpgradeAgentArgs{
		DownloadURL: srv.URL,
		SHA256:      "bogus-not-matching",
		Version:     "v9.9.9",
	}
	reply, err := s.handleUpgrade(context.Background(), args, noopStream())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if reply.Success {
		t.Fatal("should fail on sha mismatch")
	}
	got, _ := os.ReadFile(binary)
	if !bytes.Equal(got, old) {
		t.Fatalf("old binary must be preserved: %q", got)
	}
}

func TestHandleUpgradeSerialRejection(t *testing.T) {
	s, dir := newTestService(t)
	binary := filepath.Join(dir, "collia")
	if err := os.WriteFile(binary, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	if !acquireUpgradeLock() {
		t.Fatal("should acquire lock")
	}
	defer releaseUpgradeLock()

	payload := []byte("new")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	args := rpcSchema.UpgradeAgentArgs{DownloadURL: srv.URL, SHA256: sha256Hex(payload)}
	reply, err := s.handleUpgrade(context.Background(), args, noopStream())
	if err == nil {
		t.Fatal("expected serial-rejection error")
	}
	if reply == nil || reply.Success {
		t.Fatal("reply should indicate failure")
	}
}

func TestHandleUninstallReplyConsistency(t *testing.T) {
	s, _ := newTestService(t)
	var reply rpcSchema.UninstallAgentReply
	if err := s.handleUninstall(context.Background(), rpcSchema.UninstallAgentArgs{}, &reply); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	// success 与 residuals 不能自相矛盾
	if reply.Success && len(reply.Residuals) != 0 {
		t.Fatalf("inconsistent reply: success but residuals present")
	}
}
