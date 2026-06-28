package service

import "testing"

func TestBuildInstallReportPayloadUsesConfigAndPersistentInstallID(t *testing.T) {
	t.Setenv("AMPROBE_IMAGE", "amprobe:test")
	t.Setenv("AMPROBE_VERSION", "v1")
	t.Setenv("AMPROBE_PUBLIC_BASE_URL", "")
	t.Setenv("AMPROBE_HTTP_PORT", "")
	t.Setenv("AMPROBE_CONTROL_PORT", "18000")
	t.Setenv("AMPROBE_CONTAINER_NAME", "amprobe")

	config := &Config{
		Fiber: Fiber{
			Port: 8000,
		},
		AgentInstall: AgentInstall{
			PublicBaseURL: "http://127.0.0.1:1443",
			ControlPort:   17000,
		},
		InstallReport: InstallReport{
			InstallDir: "/data/amprobe",
			IDFile:     t.TempDir() + "/install.id",
		},
	}

	payload, err := buildInstallReportPayload(config)
	if err != nil {
		t.Fatal(err)
	}

	if payload.InstallID == "" {
		t.Fatal("expected install id")
	}
	if payload.Image != "amprobe:test" {
		t.Fatalf("image mismatch: %q", payload.Image)
	}
	if payload.PublicBaseURL != "http://127.0.0.1:1443" {
		t.Fatalf("public base url mismatch: %q", payload.PublicBaseURL)
	}
	if payload.HTTPPort != "8000" {
		t.Fatalf("http port mismatch: %q", payload.HTTPPort)
	}
	if payload.ControlPort != "18000" {
		t.Fatalf("control port mismatch: %q", payload.ControlPort)
	}

	nextPayload, err := buildInstallReportPayload(config)
	if err != nil {
		t.Fatal(err)
	}
	if nextPayload.InstallID != payload.InstallID {
		t.Fatalf("install id should be stable: %q != %q", nextPayload.InstallID, payload.InstallID)
	}
}
