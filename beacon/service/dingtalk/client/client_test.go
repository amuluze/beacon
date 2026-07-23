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
		{name: "official webhook", webhook: "https://oapi.dingtalk.com/robot/send?access_token=abc123", valid: true},
		{name: "official webhook on tls port", webhook: "https://oapi.dingtalk.com:443/robot/send?access_token=abc123", valid: true},
		{name: "plain http", webhook: "http://oapi.dingtalk.com/robot/send?access_token=abc123"},
		{name: "localhost", webhook: "https://127.0.0.1/robot/send?access_token=abc123"},
		{name: "lookalike host", webhook: "https://oapi.dingtalk.com.example.com/robot/send?access_token=abc123"},
		{name: "userinfo confusion", webhook: "https://oapi.dingtalk.com@127.0.0.1/robot/send?access_token=abc123"},
		{name: "custom port", webhook: "https://oapi.dingtalk.com:8443/robot/send?access_token=abc123"},
		{name: "wrong path", webhook: "https://oapi.dingtalk.com/other?access_token=abc123"},
		{name: "missing token", webhook: "https://oapi.dingtalk.com/robot/send"},
		{name: "extra query", webhook: "https://oapi.dingtalk.com/robot/send?access_token=abc123&target=http://127.0.0.1"},
		{name: "fragment", webhook: "https://oapi.dingtalk.com/robot/send?access_token=abc123#fragment"},
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

func TestWebhookSenderBuildsSignedRequest(t *testing.T) {
	t.Parallel()

	fixedTime := time.UnixMilli(1_700_000_000_123)
	secret := "SEC-test-secret"
	var receivedQuery url.Values
	var receivedPayload struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		At struct {
			IsAtAll bool `json:"isAtAll"`
		} `json:"at"`
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("decode payload: %v", err)
		}
		_, _ = w.Write([]byte(`{"errcode":0,"errmsg":"ok"}`))
	}))
	defer server.Close()

	testURL := server.URL + "/robot/send?access_token=test-token"
	sender := &webhookSender{
		httpClient: server.Client(),
		now:        func() time.Time { return fixedTime },
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
	if err := sender.Send(context.Background(), Config{Webhook: testURL, Secret: secret, AtAll: true}, "CPU 告警"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	timestamp := "1700000000123"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "\n" + secret))
	expectedSign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if got := receivedQuery.Get("timestamp"); got != timestamp {
		t.Fatalf("timestamp = %q, want %q", got, timestamp)
	}
	if got := receivedQuery.Get("sign"); got != expectedSign {
		t.Fatalf("sign = %q, want %q", got, expectedSign)
	}
	if got := receivedQuery.Get("access_token"); got != "test-token" {
		t.Fatalf("access_token = %q, want test-token", got)
	}
	if receivedPayload.MsgType != "text" || receivedPayload.Text.Content != "Beacon 主机告警\nCPU 告警" {
		t.Fatalf("unexpected payload: %+v", receivedPayload)
	}
	if !receivedPayload.At.IsAtAll {
		t.Fatal("isAtAll = false, want true")
	}
}

func TestWebhookSenderHandlesDingTalkErrorResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"errcode":310000,"errmsg":"invalid token"}`))
	}))
	defer server.Close()
	sender := testSender(server, time.Second)
	err := sender.Send(context.Background(), Config{Webhook: server.URL + "?access_token=top-secret"}, "test")
	if err == nil || !strings.Contains(err.Error(), "errcode=310000") {
		t.Fatalf("Send() error = %v, want DingTalk errcode", err)
	}
	if strings.Contains(err.Error(), "top-secret") {
		t.Fatalf("Send() leaked Webhook token: %v", err)
	}
	if strings.Contains(err.Error(), "invalid token") {
		t.Fatalf("Send() leaked DingTalk response details: %v", err)
	}
}

func TestWebhookSenderRejectsHTTPErrorResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}))
	defer server.Close()
	sender := testSender(server, time.Second)
	err := sender.Send(context.Background(), Config{Webhook: server.URL + "?access_token=top-secret"}, "test")
	if err == nil || !strings.Contains(err.Error(), "HTTP 502") {
		t.Fatalf("Send() error = %v, want HTTP status error", err)
	}
	if strings.Contains(err.Error(), "top-secret") || strings.Contains(err.Error(), "upstream unavailable") {
		t.Fatalf("Send() leaked request or response content: %v", err)
	}
}

func TestWebhookSenderTimeoutDoesNotLeakWebhook(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte(`{"errcode":0,"errmsg":"ok"}`))
	}))
	defer server.Close()
	sender := testSender(server, 20*time.Millisecond)
	err := sender.Send(context.Background(), Config{Webhook: server.URL + "?access_token=top-secret"}, "test")
	if err == nil {
		t.Fatal("Send() error = nil, want timeout")
	}
	if strings.Contains(err.Error(), "top-secret") {
		t.Fatalf("Send() leaked Webhook token: %v", err)
	}
}

func testSender(server *httptest.Server, timeout time.Duration) *webhookSender {
	httpClient := server.Client()
	httpClient.Timeout = timeout
	return &webhookSender{
		httpClient: httpClient,
		now:        time.Now,
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
}
