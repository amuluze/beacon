package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"beacon/service/model"
	"beacon/service/schema"
	"beacon/service/testutil"
	"beacon/service/wecom/client"

	"gorm.io/gorm"
)

type fakeRepository struct {
	setting model.WeCom
	err     error
	saved   *model.WeCom
}

func (r *fakeRepository) Query(context.Context) (model.WeCom, error) {
	return r.setting, r.err
}

func (r *fakeRepository) Save(_ context.Context, setting *model.WeCom) error {
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

func TestWeComServiceQueryMasksCredentials(t *testing.T) {
	repository := &fakeRepository{setting: model.WeCom{
		Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=very-secret-key-value",
		Enabled: true,
	}}
	service := NewWeComService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !got.WebhookConfigured || !got.Enabled {
		t.Fatalf("Query() = %+v, want configured status", got)
	}
	if strings.Contains(got.WebhookMasked, "very-secret-key-value") || !strings.Contains(got.WebhookMasked, "alue") {
		t.Fatalf("WebhookMasked = %q, want only masked key suffix", got.WebhookMasked)
	}
}

func TestWeComServiceQueryEmptyConfigDoesNotLookConfigured(t *testing.T) {
	service := NewWeComService(&fakeRepository{setting: model.WeCom{Key: model.DefaultWeComConfigKey}}, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.WebhookConfigured || got.WebhookMasked != "" {
		t.Fatalf("Query() = %+v, want visibly unconfigured result", got)
	}
}

func TestWeComServiceUpdatePreservesBlankWebhook(t *testing.T) {
	const webhook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=existing-key-value"
	repository := &fakeRepository{setting: model.WeCom{Webhook: webhook, Key: model.DefaultWeComConfigKey}}
	service := NewWeComService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	if err := service.Update(context.Background(), schema.WeComUpdateArgs{Enabled: true, AtAll: true}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.saved == nil || repository.saved.Webhook != webhook {
		t.Fatalf("saved = %+v, want webhook preserved", repository.saved)
	}
	if !repository.saved.Enabled || !repository.saved.AtAll {
		t.Fatalf("saved = %+v, want enabled and at_all", repository.saved)
	}
}

func TestWeComServiceRejectsSSRFWebhook(t *testing.T) {
	repository := &fakeRepository{err: gorm.ErrRecordNotFound}
	service := NewWeComService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})
	err := service.Update(context.Background(), schema.WeComUpdateArgs{
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

func TestWeComServiceTestUsesStoredCredentialsWhenDisabled(t *testing.T) {
	repository := &fakeRepository{setting: model.WeCom{
		Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=existing-key-value",
		AtAll:   true,
		Enabled: false,
	}}
	sender := &fakeSender{}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{MessageTitle: "  生产环境主机告警  "}, nil
	}
	service := NewWeComService(repository, alarmRepository, sender)
	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || !sender.configs[0].AtAll {
		t.Fatalf("sender configs = %+v, want stored config", sender.configs)
	}
	if sender.configs[0].MessageTitle != "生产环境主机告警" {
		t.Fatalf("message title = %q, want shared custom title", sender.configs[0].MessageTitle)
	}
}

func TestWeComServiceTestPropagatesSenderError(t *testing.T) {
	repository := &fakeRepository{setting: model.WeCom{
		Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=existing-key-value",
	}}
	sender := &fakeSender{err: errors.New("send failed")}
	service := NewWeComService(repository, testutil.NewFakeAlarmRepo(), sender)
	if err := service.Test(context.Background()); !errors.Is(err, sender.err) {
		t.Fatalf("Test() error = %v, want sender error", err)
	}
}

func TestWeComServiceTestFallsBackToDefaultTitle(t *testing.T) {
	repository := &fakeRepository{setting: model.WeCom{
		Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=existing-key-value",
	}}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{}, errors.New("database unavailable")
	}
	sender := &fakeSender{}
	service := NewWeComService(repository, alarmRepository, sender)

	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || sender.configs[0].MessageTitle != model.DefaultAlarmMessageTitle {
		t.Fatalf("sender configs = %+v, want default title", sender.configs)
	}
}
