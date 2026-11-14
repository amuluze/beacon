package schema

type FeishuSetting struct {
	ID                uint   `json:"id"`
	Enabled           bool   `json:"enabled"`
	WebhookMasked     string `json:"webhook_masked"`
	WebhookConfigured bool   `json:"webhook_configured"`
	SecretConfigured  bool   `json:"secret_configured"`
}

type FeishuUpdateArgs struct {
	Enabled     bool   `json:"enabled"`
	Webhook     string `json:"webhook" validate:"omitempty,lte=2048"`
	Secret      string `json:"secret" validate:"omitempty,lte=255"`
	ClearSecret bool   `json:"clear_secret"`
}
