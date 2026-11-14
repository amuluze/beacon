package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestValidateWebhookURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		webhook string
		valid   bool
	}{
		{name: "official webhook", webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123", valid: true},
		{name: "official webhook on tls port", webhook: "https://qyapi.weixin.qq.com:443/cgi-bin/webhook/send?key=abc123", valid: true},
		{name: "plain http", webhook: "http://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123"},
		{name: "localhost", webhook: "https://127.0.0.1/cgi-bin/webhook/send?key=abc123"},
		{name: "lookalike host", webhook: "https://qyapi.weixin.qq.com.example.com/cgi-bin/webhook/send?key=abc123"},
		{name: "userinfo confusion", webhook: "https://qyapi.weixin.qq.com@127.0.0.1/cgi-bin/webhook/send?key=abc123"},
		{name: "custom port", webhook: "https://qyapi.weixin.qq.com:8443/cgi-bin/webhook/send?key=abc123"},
		{name: "wrong path", webhook: "https://qyapi.weixin.qq.com/other?key=abc123"},
		{name: "missing key", webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"},
		{name: "extra query", webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123&target=http://127.0.0.1"},
		{name: "fragment", webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=abc123#fragment"},
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

func TestMaskWebhookURLKeepsKeySuffix(t *testing.T) {
	masked := MaskWebhookURL("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=very-secret-key-value")
	if !strings.HasPrefix(masked, "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=****") {
		t.Fatalf("MaskWebhookURL() = %q, want masked key prefix", masked)
	}
	if strings.Contains(masked, "very-secret-key-value") {
		t.Fatalf("MaskWebhookURL() = %q, must not contain full key", masked)
	}
}

func TestWebhookSenderBuildsMentionedPayload(t *testing.T) {
	t.Parallel()

	var receivedQuery url.Values
	var receivedPayload struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content       string   `json:"content"`
			MentionedList []string `json:"mentioned_list"`
		} `json:"text"`
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

	testURL := server.URL + "/cgi-bin/webhook/send?key=test-key"
	sender := &webhookSender{
		httpClient: server.Client(),
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
	if err := sender.Send(context.Background(), Config{
		Webhook:      testURL,
		AtAll:        true,
		MessageTitle: "生产环境主机告警",
	}, "CPU 告警"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if got := receivedQuery.Get("key"); got != "test-key" {
		t.Fatalf("key = %q, want test-key", got)
	}
	if receivedPayload.MsgType != "text" {
		t.Fatalf("msgtype = %q, want text", receivedPayload.MsgType)
	}
	if len(receivedPayload.Text.MentionedList) != 1 || receivedPayload.Text.MentionedList[0] != "@all" {
		t.Fatalf("mentioned_list = %v, want [@all]", receivedPayload.Text.MentionedList)
	}
	if !strings.Contains(receivedPayload.Text.Content, "生产环境主机告警") || !strings.Contains(receivedPayload.Text.Content, "CPU 告警") {
		t.Fatalf("content = %q, want title and message", receivedPayload.Text.Content)
	}
}

func TestWebhookSenderRejectsNonZeroErrcode(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"errcode":93000,"errmsg":"invalid webhook url"}`))
	}))
	defer server.Close()

	sender := &webhookSender{
		httpClient: server.Client(),
		validate: func(rawURL string) (*url.URL, error) {
			return url.Parse(rawURL)
		},
	}
	err := sender.Send(context.Background(), Config{Webhook: server.URL + "/cgi-bin/webhook/send?key=test-key"}, "msg")
	if err == nil || !strings.Contains(err.Error(), "93000") {
		t.Fatalf("Send() error = %v, want non-zero errcode error", err)
	}
}
