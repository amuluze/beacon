// Package report
// HTTP client that pushes monitoring data from Agent to Server.
package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	rpcSchema "common/rpc/schema"
)

// Client pushes monitoring data to the Server via HTTP POST.
type Client struct {
	url        string
	token      string
	httpClient *http.Client
}

// NewClient creates an HTTP push client.
func NewClient(url, token string) *Client {
	return &Client{
		url:   url,
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Push sends a batch of monitoring data to the Server via HTTP POST.
func (c *Client) Push(ctx context.Context, args rpcSchema.MonitorReportArgs) error {
	if ctx == nil {
		ctx = context.Background()
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
		slog.Error("report push failed", "agent", args.AgentID, "url", c.url, "error", err)
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

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
