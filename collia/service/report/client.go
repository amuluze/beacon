// Package report
// HTTP client that pushes monitoring data from Agent to Server.
package report

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	rpcSchema "common/rpc/schema"
)

var ErrInsecureTokenTransport = errors.New("refusing to send report token over non-HTTPS transport")

// Client pushes monitoring data to the Server via HTTP POST.
type Client struct {
	url        string
	token      string
	httpClient *http.Client
}

// NewClient creates an HTTP push client. When tlsConfig is non-nil it is
// applied to the underlying transport (e.g. a custom RootCAs loaded from the
// local CA plus a ServerName matching the control-plane certificate), so the
// report target can use HTTPS against a self-signed internal CA.
func NewClient(url, token string, tlsConfig *tls.Config) *Client {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	if tlsConfig != nil {
		httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	return &Client{
		url:        url,
		token:      token,
		httpClient: httpClient,
	}
}

// Push sends a batch of monitoring data to the Server via HTTP POST.
func (c *Client) Push(ctx context.Context, args rpcSchema.MonitorReportArgs) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.token != "" {
		parsed, err := url.Parse(c.url)
		if err != nil || !strings.EqualFold(parsed.Scheme, "https") || parsed.Host == "" || parsed.User != nil {
			return ErrInsecureTokenTransport
		}
	}

	body, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-Install-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("report push failed", "agent", args.AgentID, "error", err)
		return fmt.Errorf("http post: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	slog.Info("report pushed", "agent", args.AgentID,
		"cpu", args.CPU.CPUPercent,
		"mem", args.Memory.MemPercent,
		"disks", len(args.Disks),
		"containers", len(args.Containers),
	)
	return nil
}

// Close is a no-op for the HTTP client.
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
