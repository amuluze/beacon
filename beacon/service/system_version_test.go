package service

import "testing"

// 未注册 VersionChecker 时，版本接口只返回当前构建版本，不声明存在更新。
func TestGlobalSystemVersionWithoutChecker(t *testing.T) {
	original := globalVersionChecker
	globalVersionChecker = nil
	defer func() { globalVersionChecker = original }()

	got := GlobalSystemVersion()
	if got.Version != BuildVersion {
		t.Fatalf("version = %q, want %q", got.Version, BuildVersion)
	}
	if got.UpdateAvailable {
		t.Fatalf("update_available = true, want false")
	}
	if got.LatestVersion != "" || got.ReleaseNotes != "" || got.CheckedAt != "" {
		t.Fatalf("unexpected update fields: %+v", got)
	}
}
