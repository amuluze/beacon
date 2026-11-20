package rpc

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

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

func TestReplaceBinaryChmods0700(t *testing.T) {
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
	if info.Mode().Perm() != 0o700 {
		t.Fatalf("perm = %o, want 700", info.Mode().Perm())
	}
}

func TestReplaceBinaryFailureKeepsCurrentPath(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "collia")
	if err := os.WriteFile(target, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	// A directory cannot atomically replace an existing regular file, forcing
	// the final rename to fail after the backup has been created.
	tmpDir := filepath.Join(dir, "collia.new")
	if err := os.Mkdir(tmpDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := replaceBinary(tmpDir, target); err == nil {
		t.Fatal("expected rename failure")
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("current path disappeared: %v", err)
	}
	if !bytes.Equal(got, []byte("old")) {
		t.Fatalf("current binary = %q, want old", got)
	}
	backup, err := os.ReadFile(target + backupSuffix)
	if err != nil || !bytes.Equal(backup, []byte("old")) {
		t.Fatalf("backup = %q err=%v", backup, err)
	}
}

func TestScheduleRestartLeavesReplyWindow(t *testing.T) {
	s, _ := newTestService(t)
	restarted := make(chan struct{}, 1)
	s.restartFn = func() error {
		restarted <- struct{}{}
		return nil
	}
	s.scheduleRestart()

	select {
	case <-restarted:
		t.Fatal("restart ran before reply window elapsed")
	case <-time.After(100 * time.Millisecond):
	}
	select {
	case <-restarted:
	case <-time.After(2 * time.Second):
		t.Fatal("restart callback was not invoked")
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

func TestVerifyUpgradeSignatureMatrix(t *testing.T) {
	dir := t.TempDir()
	tmp := filepath.Join(dir, "collia")
	if err := os.WriteFile(tmp, []byte("signed-binary"), 0o644); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("signed-binary"))
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, otherPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	validSig := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, digest[:]))
	wrongSig := base64.StdEncoding.EncodeToString(ed25519.Sign(otherPriv, digest[:]))

	s, _ := newTestService(t)
	// 未固定公钥 + 无签名：兼容模式放行（存量部署）。
	if err := s.verifyUpgradeSignature(tmp, ""); err != nil {
		t.Fatalf("legacy unsigned upgrade should pass, got %v", err)
	}
	// 未固定公钥 + 有签名：无法验证，拒绝。
	if err := s.verifyUpgradeSignature(tmp, validSig); err == nil {
		t.Fatal("signature without a pinned public key must be rejected")
	}

	s.upgradePubKeys = []ed25519.PublicKey{pub}
	// 固定公钥 + 缺失签名：fail-closed 拒绝。
	if err := s.verifyUpgradeSignature(tmp, ""); err == nil {
		t.Fatal("pinned public key must reject unsigned upgrades")
	}
	// 固定公钥 + 有效签名：放行。
	if err := s.verifyUpgradeSignature(tmp, validSig); err != nil {
		t.Fatalf("valid signature should pass, got %v", err)
	}
	// 固定公钥 + 非本公钥签名：拒绝。
	if err := s.verifyUpgradeSignature(tmp, wrongSig); err == nil {
		t.Fatal("signature from another key must be rejected")
	}
	// 固定公钥 + 非法 base64：拒绝。
	if err := s.verifyUpgradeSignature(tmp, "not-base64!!!"); err == nil {
		t.Fatal("malformed signature must be rejected")
	}
}

func TestSetUpgradePublicKeys(t *testing.T) {
	s, _ := newTestService(t)
	if err := s.SetUpgradePublicKeys(nil); err != nil || s.upgradePubKeys != nil {
		t.Fatalf("empty key set should leave the agent unpinned, err=%v", err)
	}
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	encoded := base64.StdEncoding.EncodeToString(pub)
	if err := s.SetUpgradePublicKeys([]string{encoded}); err != nil {
		t.Fatalf("valid public key: %v", err)
	}
	if len(s.upgradePubKeys) != 1 || len(s.upgradePubKeys[0]) != ed25519.PublicKeySize {
		t.Fatalf("pinned key set = %v", s.upgradePubKeys)
	}
	// 重复键去重，空串忽略。
	if err := s.SetUpgradePublicKeys([]string{encoded, encoded, "  "}); err != nil {
		t.Fatalf("dedupe keys: %v", err)
	}
	if len(s.upgradePubKeys) != 1 {
		t.Fatalf("duplicate keys must be deduped, got %d", len(s.upgradePubKeys))
	}
	if err := s.SetUpgradePublicKeys([]string{"not-base64!!!"}); err == nil {
		t.Fatal("malformed key must fail, not silently downgrade")
	}
	if err := s.SetUpgradePublicKeys([]string{base64.StdEncoding.EncodeToString([]byte("short"))}); err == nil {
		t.Fatal("wrong-size key must fail, not silently downgrade")
	}
}

// TestVerifyUpgradeSignatureAcceptsRotatedKeyRing 验证轮换重叠窗口：
// 旧钥与新钥同时固定时两种签名都放行，集合外的第三把钥签名拒绝。
func TestVerifyUpgradeSignatureAcceptsRotatedKeyRing(t *testing.T) {
	dir := t.TempDir()
	tmp := filepath.Join(dir, "collia")
	if err := os.WriteFile(tmp, []byte("rotate-binary"), 0o644); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("rotate-binary"))
	oldPub, oldPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	newPub, newPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, roguePriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	oldSig := base64.StdEncoding.EncodeToString(ed25519.Sign(oldPriv, digest[:]))
	newSig := base64.StdEncoding.EncodeToString(ed25519.Sign(newPriv, digest[:]))
	rogueSig := base64.StdEncoding.EncodeToString(ed25519.Sign(roguePriv, digest[:]))

	s, _ := newTestService(t)
	s.upgradePubKeys = []ed25519.PublicKey{oldPub, newPub}
	if err := s.verifyUpgradeSignature(tmp, oldSig); err != nil {
		t.Fatalf("old key during rotation must verify: %v", err)
	}
	if err := s.verifyUpgradeSignature(tmp, newSig); err != nil {
		t.Fatalf("new key during rotation must verify: %v", err)
	}
	if err := s.verifyUpgradeSignature(tmp, rogueSig); err == nil {
		t.Fatal("signature outside the key ring must be rejected")
	}
}

func TestHandleUpgradeWithPinnedKeyRequiresValidSignature(t *testing.T) {
	s, dir := newTestService(t)
	binary := filepath.Join(dir, "collia")
	if err := os.WriteFile(binary, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	s.upgradePubKeys = []ed25519.PublicKey{pub}

	payload := []byte("new-signed-binary")
	digest := sha256.Sum256(payload)
	validSig := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, digest[:]))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	// 缺签名的升级必须被拒绝且保留旧二进制。
	args := rpcSchema.UpgradeAgentArgs{DownloadURL: srv.URL, SHA256: sha256Hex(payload)}
	reply, err := s.handleUpgrade(context.Background(), args, noopStream())
	if err != nil {
		t.Fatalf("handleUpgrade: %v", err)
	}
	if reply.Success {
		t.Fatal("unsigned upgrade must fail with a pinned key")
	}
	got, _ := os.ReadFile(binary)
	if string(got) != "old" {
		t.Fatalf("old binary must be preserved, got %q", got)
	}

	// 携带有效签名的同一升级成功。
	args.Signature = validSig
	releaseUpgradeLockForTest(t)
	reply, err = s.handleUpgrade(context.Background(), args, noopStream())
	if err != nil {
		t.Fatalf("handleUpgrade signed: %v", err)
	}
	if !reply.Success {
		t.Fatalf("signed upgrade should succeed, got %+v", reply)
	}
	got, _ = os.ReadFile(binary)
	if !bytes.Equal(got, payload) {
		t.Fatalf("binary not replaced: %q", got)
	}
}

// releaseUpgradeLockForTest 清理由上一次 handleUpgrade 正常释放的串行锁；
// 若锁被意外占用则直接失败，保证两次升级调用互不干扰。
func releaseUpgradeLockForTest(t *testing.T) {
	t.Helper()
	if inProgress {
		t.Fatal("upgrade lock is still held after a completed upgrade")
	}
}
