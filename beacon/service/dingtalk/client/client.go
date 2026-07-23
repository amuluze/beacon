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

	"github.com/google/wire"
)

const (
	dingTalkWebhookHost = "oapi.dingtalk.com"
	dingTalkWebhookPath = "/robot/send"
	maxResponseBytes    = 64 * 1024
)

var Set = wire.NewSet(NewSender)

type Config struct {
	Webhook string
	Secret  string
	AtAll   bool
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
// used as an SSRF primitive. DingTalk group robot Webhooks use this one HTTPS
// host, path and access_token query parameter.
func ValidateWebhookURL(rawURL string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return nil, errors.New("钉钉 Webhook 格式无效")
	}
	if parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), dingTalkWebhookHost) {
		return nil, errors.New("钉钉 Webhook 必须使用官方 HTTPS 地址")
	}
	if parsed.User != nil || (parsed.Port() != "" && parsed.Port() != "443") || parsed.Fragment != "" {
		return nil, errors.New("钉钉 Webhook 地址包含不允许的内容")
	}
	if parsed.Path != dingTalkWebhookPath || parsed.EscapedPath() != dingTalkWebhookPath {
		return nil, errors.New("钉钉 Webhook 路径无效")
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil || len(query) != 1 {
		return nil, errors.New("钉钉 Webhook 查询参数无效")
	}
	tokens, ok := query["access_token"]
	if !ok || len(tokens) != 1 || strings.TrimSpace(tokens[0]) == "" {
		return nil, errors.New("钉钉 Webhook 缺少 access_token")
	}
	return parsed, nil
}

func MaskWebhookURL(rawURL string) string {
	parsed, err := ValidateWebhookURL(rawURL)
	if err != nil {
		return "已配置"
	}
	token := parsed.Query().Get("access_token")
	maskedToken := "****"
	if len(token) > 4 {
		maskedToken += token[len(token)-4:]
	}
	parsed.RawQuery = url.Values{"access_token": []string{maskedToken}}.Encode()
	return parsed.String()
}

func (s *webhookSender) Send(ctx context.Context, config Config, message string) error {
	endpoint, err := s.validate(config.Webhook)
	if err != nil {
		return err
	}

	endpoint = cloneURL(endpoint)
	if config.Secret != "" {
		addSignature(endpoint, config.Secret, s.now())
	}

	payload := struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		At struct {
			IsAtAll bool `json:"isAtAll"`
		} `json:"at"`
	}{MsgType: "text"}
	payload.Text.Content = "Beacon 主机告警\n" + message
	payload.At.IsAtAll = config.AtAll

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("编码钉钉告警消息失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return errors.New("创建钉钉告警请求失败")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("发送钉钉告警失败或超时")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponseBytes))
		return fmt.Errorf("钉钉 Webhook 返回 HTTP %d", resp.StatusCode)
	}

	var result struct {
		ErrCode int `json:"errcode"`
	}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes))
	if err := decoder.Decode(&result); err != nil {
		return errors.New("钉钉 Webhook 返回了无效响应")
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("钉钉 Webhook 拒绝请求（errcode=%d）", result.ErrCode)
	}
	return nil
}

func cloneURL(source *url.URL) *url.URL {
	cloned := *source
	return &cloned
}

func addSignature(endpoint *url.URL, secret string, now time.Time) {
	timestamp := strconv.FormatInt(now.UnixMilli(), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "\n" + secret))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	query := endpoint.Query()
	query.Set("timestamp", timestamp)
	query.Set("sign", sign)
	endpoint.RawQuery = query.Encode()
}
