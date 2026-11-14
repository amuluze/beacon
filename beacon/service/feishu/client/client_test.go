package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestValidateWebhookURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		webhook string
		valid   bool
	}{
		{name: "official webhook", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/abc123def456", valid: true},
		{name: "official webhook on tls port", webhook: "https://open.feishu.cn:443/open-apis/bot/v2/hook/abc123def456", valid: true},
		{name: "plain http", webhook: "http://open.feishu.cn/open-apis/bot/v2/hook/abc123def456"},
		{name: "localhost", webhook: "https://127.0.0.1/open-apis/bot/v2/hook/abc123def456"},
		{name: "lookalike host", webhook: "https://open.feishu.cn.example.com/open-apis/bot/v2/hook/abc123def456"},
		{name: "userinfo confusion", webhook: "https://open.feishu.cn@127.0.0.1/open-apis/bot/v2/hook/abc123def456"},
		{name: "custom port", webhook: "https://open.feishu.cn:8443/open-apis/bot/v2/hook/abc123def456"},
		{name: "wrong path", webhook: "https://open.feishu.cn/other/hook/abc123def456"},
		{name: "missing hook id", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/"},
		{name: "short hook id", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/ab"},
		{name: "nested path", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/abc/def"},
		{name: "query not allowed", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/abc123def456?redirect=http://127.0.0.1"},
		{name: "fragment", webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/abc123def456#fragment"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateWebhookURL(tt.webhook)
			if (err == nil) != tt.valid {
				t.Fatalf("ValidateWebhookURL() error = %v, valid = %v", err, tt.valid)
			}
		})
	}
}

func TestMaskWebhookURLKeepsHookSuffix(t *testing.T) {
	masked := MaskWebhookURL("https://open.feishu.cn/open-apis/bot/v2/hook/very-secret-hook-id")
	if !strings.HasPrefix(masked, "https://open.feishu.cn/open-apis/bot/v2/hook/****") {
		t.Fatalf("MaskWebhookURL() = %q, want masked hook id prefix", masked)
	}
	if strings.Contains(masked, "very-secret-hook-id") {
		t.Fatalf("MaskWebhookURL() = %q, must not contain full hook id", masked)
	}
}

func TestWebhookSenderBuildsSignedRequest(t *testing.T) {
	t.Parallel()

	fixedTime := time.Unix(1_700_000_000, 0)
	secret := "feishu-sign-secret"
	var receivedPayload struct {
		Timestamp string `json:"timestamp"`
		Sign      string `json:"sign"`
		MsgType   string `json:"msg_type"`
		Content   struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.RawQuery; got != "" {
			t.Errorf("RawQuery = %q, want empty (feishu signs in body)", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("decode payload: %v", err)
		}
		_, _ = w.Write([]byte(`{"code":0,"msg":"success"}`))
	}))
	defer server.Close()

	testURL := server.URL + "/open-apis/bot/v2/hook/test-hook-id"
	sender := &webhookSender{
		httpClient: server.Client(),
		now:        func() time.Time { return fixedTime },
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
	if err := sender.Send(context.Background(), Config{
		Webhook:      testURL,
		Secret:       secret,
		MessageTitle: "生产环境主机告警",
	}, "CPU 告警"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	timestamp := "1700000000"
	// 飞书签名：以 timestamp+"\n"+secret 作为 HMAC-SHA256 密钥，对空消息签名。
	mac := hmac.New(sha256.New, []byte(timestamp+"\n"+secret))
	_, _ = mac.Write([]byte(""))
	expectedSign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if got := receivedPayload.Timestamp; got != timestamp {
		t.Fatalf("timestamp = %q, want %q", got, timestamp)
	}
	if got := receivedPayload.Sign; got != expectedSign {
		t.Fatalf("sign = %q, want %q", got, expectedSign)
	}
	if receivedPayload.MsgType != "text" {
		t.Fatalf("msg_type = %q, want text", receivedPayload.MsgType)
	}
	if !strings.Contains(receivedPayload.Content.Text, "生产环境主机告警") || !strings.Contains(receivedPayload.Content.Text, "CPU 告警") {
		t.Fatalf("content = %q, want title and message", receivedPayload.Content.Text)
	}
}

func TestWebhookSenderRejectsNonZeroCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"code":19021,"msg":"Key Words Not Found"}`))
	}))
	defer server.Close()

	sender := &webhookSender{
		httpClient: server.Client(),
		now:        time.Now,
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
	err := sender.Send(context.Background(), Config{Webhook: server.URL + "/open-apis/bot/v2/hook/test-hook-id"}, "msg")
	if err == nil || !strings.Contains(err.Error(), "19021") {
		t.Fatalf("Send() error = %v, want non-zero code error", err)
	}
}
