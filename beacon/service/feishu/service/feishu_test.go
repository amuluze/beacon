package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"beacon/service/feishu/client"
	"beacon/service/model"
	"beacon/service/schema"
	"beacon/service/testutil"

	"gorm.io/gorm"
)

type fakeRepository struct {
	setting model.Feishu
	err     error
	saved   *model.Feishu
}

func (r *fakeRepository) Query(context.Context) (model.Feishu, error) {
	return r.setting, r.err
}

func (r *fakeRepository) Save(_ context.Context, setting *model.Feishu) error {
	copy := *setting
	r.saved = &copy
	r.setting = copy
	r.err = nil
	return nil
}

type fakeSender struct {
	configs  []client.Config
	messages []string
	err      error
}

func (s *fakeSender) Send(_ context.Context, config client.Config, message string) error {
	s.configs = append(s.configs, config)
	s.messages = append(s.messages, message)
	return s.err
}

func TestFeishuServiceQueryMasksCredentials(t *testing.T) {
	repository := &fakeRepository{setting: model.Feishu{
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/very-secret-hook-id",
		Secret:  "feishu-sign-secret",
		Enabled: true,
	}}
	service := NewFeishuService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !got.WebhookConfigured || !got.SecretConfigured || !got.Enabled {
		t.Fatalf("Query() = %+v, want configured status", got)
	}
	if strings.Contains(got.WebhookMasked, "very-secret-hook-id") || !strings.Contains(got.WebhookMasked, "ook") {
		t.Fatalf("WebhookMasked = %q, want only masked hook suffix", got.WebhookMasked)
	}
}

func TestFeishuServiceQueryEmptyConfigDoesNotLookConfigured(t *testing.T) {
	service := NewFeishuService(&fakeRepository{setting: model.Feishu{Key: model.DefaultFeishuConfigKey}}, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.WebhookConfigured || got.SecretConfigured || got.WebhookMasked != "" {
		t.Fatalf("Query() = %+v, want visibly unconfigured result", got)
	}
}

func TestFeishuServiceUpdatePreservesBlankCredentials(t *testing.T) {
	const webhook = "https://open.feishu.cn/open-apis/bot/v2/hook/existing-hook-id"
	repository := &fakeRepository{setting: model.Feishu{Webhook: webhook, Secret: "feishu-existing-secret", Key: model.DefaultFeishuConfigKey}}
	service := NewFeishuService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	if err := service.Update(context.Background(), schema.FeishuUpdateArgs{Enabled: true}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.saved == nil || repository.saved.Webhook != webhook || repository.saved.Secret != "feishu-existing-secret" {
		t.Fatalf("saved = %+v, want credentials preserved", repository.saved)
	}
	if !repository.saved.Enabled {
		t.Fatalf("saved = %+v, want enabled", repository.saved)
	}
}

func TestFeishuServiceUpdateCanClearSecret(t *testing.T) {
	repository := &fakeRepository{setting: model.Feishu{
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/existing-hook-id",
		Secret:  "feishu-existing-secret",
		Key:     model.DefaultFeishuConfigKey,
	}}
	service := NewFeishuService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})
	if err := service.Update(context.Background(), schema.FeishuUpdateArgs{ClearSecret: true}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.saved.Secret != "" {
		t.Fatalf("saved secret = %q, want empty", repository.saved.Secret)
	}
}

func TestFeishuServiceRejectsSSRFWebhook(t *testing.T) {
	repository := &fakeRepository{err: gorm.ErrRecordNotFound}
	service := NewFeishuService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})
	err := service.Update(context.Background(), schema.FeishuUpdateArgs{
		Enabled: true,
		Webhook: "http://127.0.0.1/internal",
	})
	if err == nil {
		t.Fatal("Update() error = nil, want validation error")
	}
	if repository.saved != nil {
		t.Fatalf("Save() called with %+v for invalid Webhook", repository.saved)
	}
}

func TestFeishuServiceTestUsesStoredCredentialsWhenDisabled(t *testing.T) {
	repository := &fakeRepository{setting: model.Feishu{
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/existing-hook-id",
		Secret:  "feishu-existing-secret",
		Enabled: false,
	}}
	sender := &fakeSender{}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{MessageTitle: "  生产环境主机告警  "}, nil
	}
	service := NewFeishuService(repository, alarmRepository, sender)
	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || sender.configs[0].Secret != "feishu-existing-secret" {
		t.Fatalf("sender configs = %+v, want stored config", sender.configs)
	}
	if sender.configs[0].MessageTitle != "生产环境主机告警" {
		t.Fatalf("message title = %q, want shared custom title", sender.configs[0].MessageTitle)
	}
}

func TestFeishuServiceTestPropagatesSenderError(t *testing.T) {
	repository := &fakeRepository{setting: model.Feishu{
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/existing-hook-id",
	}}
	sender := &fakeSender{err: errors.New("send failed")}
	service := NewFeishuService(repository, testutil.NewFakeAlarmRepo(), sender)
	if err := service.Test(context.Background()); !errors.Is(err, sender.err) {
		t.Fatalf("Test() error = %v, want sender error", err)
	}
}

func TestFeishuServiceTestFallsBackToDefaultTitle(t *testing.T) {
	repository := &fakeRepository{setting: model.Feishu{
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/existing-hook-id",
	}}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{}, errors.New("database unavailable")
	}
	sender := &fakeSender{}
	service := NewFeishuService(repository, alarmRepository, sender)

	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || sender.configs[0].MessageTitle != model.DefaultAlarmMessageTitle {
		t.Fatalf("sender configs = %+v, want default title", sender.configs)
	}
}
