package schema

type WeComSetting struct {
	ID                uint   `json:"id"`
	Enabled           bool   `json:"enabled"`
	WebhookMasked     string `json:"webhook_masked"`
	WebhookConfigured bool   `json:"webhook_configured"`
	AtAll             bool   `json:"at_all"`
}

type WeComUpdateArgs struct {
	Enabled bool   `json:"enabled"`
	Webhook string `json:"webhook" validate:"omitempty,lte=2048"`
	AtAll   bool   `json:"at_all"`
}
