package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const installReportTimeout = 5 * time.Second
const defaultInstallReportIDFile = "/app/data/install.id"

type installReportPayload struct {
	InstallID     string `json:"install_id"`
	Image         string `json:"image"`
	Version       string `json:"version"`
	PublicBaseURL string `json:"public_base_url"`
	InstallDir    string `json:"install_dir"`
	HTTPPort      string `json:"http_port"`
	ControlPort   string `json:"control_port"`
	ContainerName string `json:"container_name"`
	Hostname      string `json:"hostname"`
}

func ReportInstallation(ctx context.Context, config *Config) {
	if config == nil || !config.InstallReport.Enable {
		return
	}
	reportURL := strings.TrimSpace(config.InstallReport.URL)
	if reportURL == "" {
		return
	}

	payload, err := buildInstallReportPayload(config)
	if err != nil {
		slog.Warn("build install report payload failed", "err", err)
		return
	}
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("marshal install report failed", "err", err)
		return
	}

	timeout := installReportTimeout
	if config.InstallReport.Timeout > 0 {
		timeout = time.Duration(config.InstallReport.Timeout) * time.Second
	}
	reportCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reportCtx, http.MethodPost, reportURL, bytes.NewReader(body))
	if err != nil {
		slog.Warn("create install report request failed", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "beacon-install-reporter")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn("send install report failed", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		slog.Warn("install report returned non-2xx status", "status", resp.StatusCode)
		return
	}
	slog.Info("install report sent", "url", reportURL)
}

func buildInstallReportPayload(config *Config) (installReportPayload, error) {
	hostname, _ := os.Hostname()
	installID, err := loadOrCreateInstallID(config.InstallReport.IDFile)
	if err != nil {
		return installReportPayload{}, err
	}
	return installReportPayload{
		InstallID:     installID,
		Image:         firstEnv("BEACON_IMAGE", "AMPROBE_IMAGE"),
		Version:       firstEnv("BEACON_VERSION", "AMPROBE_VERSION"),
		PublicBaseURL: firstNonEmpty(firstEnv("BEACON_PUBLIC_BASE_URL", "AMPROBE_PUBLIC_BASE_URL"), config.AgentInstall.PublicBaseURL),
		InstallDir:    config.InstallReport.InstallDir,
		HTTPPort:      firstNonEmpty(firstEnv("BEACON_HTTP_PORT", "AMPROBE_HTTP_PORT"), intString(config.Fiber.Port)),
		ControlPort:   firstNonEmpty(firstEnv("BEACON_CONTROL_PORT", "AMPROBE_CONTROL_PORT"), intString(config.AgentInstall.ControlPort)),
		ContainerName: firstEnv("BEACON_CONTAINER_NAME", "AMPROBE_CONTAINER_NAME"),
		Hostname:      hostname,
	}, nil
}

func loadOrCreateInstallID(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		path = defaultInstallReportIDFile
	}
	if data, err := os.ReadFile(path); err == nil {
		id := strings.TrimSpace(string(data))
		if id != "" {
			return id, nil
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	id, err := randomInstallID()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(id+"\n"), 0600); err != nil {
		return "", err
	}
	return id, nil
}

func randomInstallID() (string, error) {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", fmt.Errorf("generate install id: %w", err)
	}
	return hex.EncodeToString(data[:]), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func intString(value int) string {
	if value == 0 {
		return ""
	}
	return strconv.Itoa(value)
}
