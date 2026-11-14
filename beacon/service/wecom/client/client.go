package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	alarmtitle "beacon/service/alarm/title"

	"github.com/google/wire"
)

const (
	weComWebhookSchemeHost = "https://qyapi.weixin.qq.com"
	weComWebhookHost       = "qyapi.weixin.qq.com"
	weComWebhookPath       = "/cgi-bin/webhook/send"
	maxResponseBytes       = 64 * 1024
	maskedCredentialTail   = 4
)

var Set = wire.NewSet(NewSender)

type Config struct {
	Webhook      string
	AtAll        bool
	MessageTitle string
}

type Sender interface {
	Send(ctx context.Context, config Config, message string) error
}

type webhookSender struct {
	httpClient *http.Client
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
		validate: ValidateWebhookURL,
	}
}

// ValidateWebhookURL applies a strict allowlist so a stored Webhook cannot be
// used as an SSRF primitive. WeCom group robot Webhooks use this one HTTPS
// host, path and key query parameter.
func ValidateWebhookURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, errors.New("企业微信 Webhook 格式无效")
	}
	if parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), weComWebhookHost) {
		return nil, errors.New("企业微信 Webhook 必须使用官方 HTTPS 地址")
	}
	if parsed.User != nil || (parsed.Port() != "" && parsed.Port() != "443") || parsed.Fragment != "" {
		return nil, errors.New("企业微信 Webhook 地址包含不允许的内容")
	}
	if parsed.EscapedPath() != weComWebhookPath {
		return nil, errors.New("企业微信 Webhook 路径无效")
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil || len(query) != 1 {
		return nil, errors.New("企业微信 Webhook 查询参数无效")
	}
	keys, ok := query["key"]
	if !ok || len(keys) != 1 || strings.TrimSpace(keys[0]) == "" {
		return nil, errors.New("企业微信 Webhook 缺少 key")
	}
	return parsed, nil
}

// MaskWebhookURL keeps the key tail so operators can tell which robot is
// configured without exposing the full credential.
func MaskWebhookURL(rawURL string) string {
	parsed, err := ValidateWebhookURL(rawURL)
	if err != nil {
		return "已配置"
	}
	key := parsed.Query().Get("key")
	maskedKey := "****"
	if len(key) > maskedCredentialTail {
		maskedKey += key[len(key)-maskedCredentialTail:]
	}
	// 手工拼接：url.Values.Encode 会把掩码星号百分号编码，破坏可读性。
	return weComWebhookSchemeHost + weComWebhookPath + "?key=" + maskedKey
}

func (s *webhookSender) Send(ctx context.Context, config Config, message string) error {
	endpoint, err := s.validate(config.Webhook)
	if err != nil {
		return err
	}

	payload := struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content       string   `json:"content"`
			MentionedList []string `json:"mentioned_list,omitempty"`
		} `json:"text"`
	}{MsgType: "text"}
	payload.Text.Content = alarmtitle.OrDefault(config.MessageTitle) + "\n" + message
	if config.AtAll {
		payload.Text.MentionedList = []string{"@all"}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("编码企业微信告警消息失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return errors.New("创建企业微信告警请求失败")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("发送企业微信告警失败或超时")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
		return fmt.Errorf("企业微信 Webhook 返回 HTTP %d", resp.StatusCode)
	}

	var result struct {
		ErrCode int `json:"errcode"`
	}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes))
	if err := decoder.Decode(&result); err != nil {
		return errors.New("企业微信 Webhook 返回了无效响应")
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("企业微信 Webhook 拒绝请求（errcode=%d）", result.ErrCode)
	}
	return nil
}
