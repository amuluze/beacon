package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	alarmtitle "beacon/service/alarm/title"

	"github.com/google/wire"
)

const (
	feishuWebhookHost  = "open.feishu.cn"
	feishuWebhookPath  = "/open-apis/bot/v2/hook"
	maxResponseBytes   = 64 * 1024
	minHookIDLength    = 8
	maskedCredentialTail = 4
)

var Set = wire.NewSet(NewSender)

type Config struct {
	Webhook      string
	Secret       string
	MessageTitle string
}

type Sender interface {
	Send(ctx context.Context, config Config, message string) error
}

type webhookSender struct {
	httpClient *http.Client
	now        func() time.Time
	validate   func(string) (*url.URL, error)
}

func NewSender() Sender {
	return &webhookSender{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		now:      time.Now,
		validate: ValidateWebhookURL,
	}
}

// ValidateWebhookURL applies a strict allowlist so a stored Webhook cannot be
// used as an SSRF primitive. Feishu group robot Webhooks use this one HTTPS
// host, the bot hook path and a hook id segment; no query parameters allowed.
func ValidateWebhookURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, errors.New("飞书 Webhook 格式无效")
	}
	if parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), feishuWebhookHost) {
		return nil, errors.New("飞书 Webhook 必须使用官方 HTTPS 地址")
	}
	if parsed.User != nil || (parsed.Port() != "" && parsed.Port() != "443") || parsed.Fragment != "" {
		return nil, errors.New("飞书 Webhook 地址包含不允许的内容")
	}
	if !strings.HasPrefix(parsed.EscapedPath(), feishuWebhookPath+"/") {
		return nil, errors.New("飞书 Webhook 路径无效")
	}
	hookID := strings.TrimPrefix(parsed.EscapedPath(), feishuWebhookPath+"/")
	if len(hookID) < minHookIDLength || strings.Contains(hookID, "/") {
		return nil, errors.New("飞书 Webhook 缺少有效的 hook 标识")
	}
	if parsed.RawQuery != "" {
		return nil, errors.New("飞书 Webhook 不允许携带查询参数")
	}
	return parsed, nil
}

// MaskWebhookURL keeps the hook id tail so operators can tell which robot is
// configured without exposing the full credential.
func MaskWebhookURL(rawURL string) string {
	parsed, err := ValidateWebhookURL(rawURL)
	if err != nil {
		return "已配置"
	}
	hookID := strings.TrimPrefix(parsed.EscapedPath(), feishuWebhookPath+"/")
	maskedHook := "****"
	if len(hookID) > maskedCredentialTail {
		maskedHook += hookID[len(hookID)-maskedCredentialTail:]
	}
	// 手工拼接：url.URL.String 会把掩码星号百分号编码，破坏可读性。
	return "https://" + feishuWebhookHost + feishuWebhookPath + "/" + maskedHook
}

func (s *webhookSender) Send(ctx context.Context, config Config, message string) error {
	endpoint, err := s.validate(config.Webhook)
	if err != nil {
		return err
	}

	payload := struct {
		Timestamp string `json:"timestamp,omitempty"`
		Sign      string `json:"sign,omitempty"`
		MsgType   string `json:"msg_type"`
		Content   struct {
			Text string `json:"text"`
		} `json:"content"`
	}{MsgType: "text"}
	payload.Content.Text = alarmtitle.OrDefault(config.MessageTitle) + "\n" + message
	if config.Secret != "" {
		timestamp := strconv.FormatInt(s.now().Unix(), 10)
		payload.Timestamp = timestamp
		payload.Sign = feishuSignature(config.Secret, timestamp)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("编码飞书告警消息失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return errors.New("创建飞书告警请求失败")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("发送飞书告警失败或超时")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
		return fmt.Errorf("飞书 Webhook 返回 HTTP %d", resp.StatusCode)
	}

	var result struct {
		Code       int `json:"code"`
		StatusCode int `json:"StatusCode"`
	}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes))
	if err := decoder.Decode(&result); err != nil {
		return errors.New("飞书 Webhook 返回了无效响应")
	}
	if result.Code != 0 || result.StatusCode != 0 {
		return fmt.Errorf("飞书 Webhook 拒绝请求（code=%d）", result.Code)
	}
	return nil
}

// feishuSignature implements the Feishu custom bot signing scheme: the HMAC
// key is "timestamp\nsecret" over an empty message, Base64 encoded.
func feishuSignature(secret string, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(timestamp+"\n"+secret))
	_, _ = mac.Write([]byte(""))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
