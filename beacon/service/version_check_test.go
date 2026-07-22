package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestChecker(t *testing.T, url string) *VersionChecker {
	t.Helper()
	vc := NewVersionChecker(&Config{Update: Update{Enable: true, URL: url, CheckInterval: 3600}})
	// 测试中不依赖 ldflags 注入的版本
	vc.current = "v3.0.0"
	return vc
}

func TestCheckOnceUpdateAvailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c := r.URL.Query().Get("current"); c != "v3.0.0" {
			t.Errorf("current query = %q, want v3.0.0", c)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"latest_version":       "v3.1.0",
			"min_required_version": "v3.0.0",
			"update_available":     true,
			"release_notes":        "bug fixes",
		})
	}))
	defer srv.Close()

	vc := newTestChecker(t, srv.URL)
	vc.checkOnce()
	st := vc.Status()
	if !st.UpdateAvailable {
		t.Fatalf("expected update available, got %+v", st)
	}
	if st.LatestVersion != "v3.1.0" {
		t.Fatalf("latest = %s, want v3.1.0", st.LatestVersion)
	}
	if st.LastError != "" {
		t.Fatalf("unexpected last_error: %s", st.LastError)
	}
}

func TestCheckOnceNon200DegradesGracefully(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	vc := newTestChecker(t, srv.URL)
	vc.checkOnce()
	st := vc.Status()
	if st.UpdateAvailable {
		t.Fatal("should not flag update on 503")
	}
	if st.LastError == "" {
		t.Fatal("expected last_error on non-200")
	}
}

func TestCheckOnceBadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	vc := newTestChecker(t, srv.URL)
	vc.checkOnce()
	if vc.Status().LastError == "" {
		t.Fatal("expected parse error recorded")
	}
}

func TestStopIsIdempotent(t *testing.T) {
	vc := newTestChecker(t, "https://example.invalid")
	vc.Stop()
	vc.Stop() // 不能 panic / double close
}

func TestIntervalClampedToMin(t *testing.T) {
	vc := NewVersionChecker(&Config{Update: Update{Enable: true, URL: "https://x", CheckInterval: 10}})
	if vc.tick < minCheckInterval {
		t.Fatalf("tick = %v, want >= %v", vc.tick, minCheckInterval)
	}
}

func TestRunStopsOnStop(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"latest_version":"v3.0.0","update_available":false}`))
	}))
	defer srv.Close()

	vc := newTestChecker(t, srv.URL)
	// 缩短 tick 以便快速验证 Run 能被 Stop 中断
	vc.tick = 50 * time.Millisecond
	done := make(chan struct{})
	go func() {
		vc.Run()
		close(done)
	}()
	// 确保至少触发一次 checkOnce
	var ran atomic.Bool
	go func() {
		time.Sleep(120 * time.Millisecond)
		ran.Store(vc.Status().CheckedAt != "")
		vc.Stop()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop within 2s")
	}
	if !ran.Load() {
		t.Fatal("expected at least one check to have run")
	}
}
