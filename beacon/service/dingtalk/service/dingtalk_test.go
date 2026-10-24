package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"beacon/service/dingtalk/client"
	"beacon/service/model"
	"beacon/service/schema"
	"beacon/service/testutil"

	"gorm.io/gorm"
)

type fakeRepository struct {
	setting model.DingTalk
	err     error
	saved   *model.DingTalk
}

func (r *fakeRepository) Query(context.Context) (model.DingTalk, error) {
	return r.setting, r.err
}

func (r *fakeRepository) Save(_ context.Context, setting *model.DingTalk) error {
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

func TestDingTalkServiceQueryMasksCredentials(t *testing.T) {
	repository := &fakeRepository{setting: model.DingTalk{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=very-secret-token",
		Secret:  "SEC-secret",
		Enabled: true,
	}}
	service := NewDingTalkService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if !got.WebhookConfigured || !got.SecretConfigured || !got.Enabled {
		t.Fatalf("Query() = %+v, want configured status", got)
	}
	if strings.Contains(got.WebhookMasked, "very-secret-token") || !strings.Contains(got.WebhookMasked, "oken") {
		t.Fatalf("WebhookMasked = %q, want only masked token suffix", got.WebhookMasked)
	}
}

func TestDingTalkServiceQueryEmptyConfigDoesNotLookConfigured(t *testing.T) {
	service := NewDingTalkService(&fakeRepository{setting: model.DingTalk{Key: model.DefaultDingTalkConfigKey}}, testutil.NewFakeAlarmRepo(), &fakeSender{})

	got, err := service.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.WebhookConfigured || got.SecretConfigured || got.WebhookMasked != "" {
		t.Fatalf("Query() = %+v, want visibly unconfigured result", got)
	}
}

func TestDingTalkServiceUpdatePreservesBlankCredentials(t *testing.T) {
	const webhook = "https://oapi.dingtalk.com/robot/send?access_token=existing-token"
	repository := &fakeRepository{setting: model.DingTalk{Webhook: webhook, Secret: "SEC-existing", Key: model.DefaultDingTalkConfigKey}}
	service := NewDingTalkService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})

	if err := service.Update(context.Background(), schema.DingTalkUpdateArgs{Enabled: true, AtAll: true}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.saved == nil || repository.saved.Webhook != webhook || repository.saved.Secret != "SEC-existing" {
		t.Fatalf("saved = %+v, want credentials preserved", repository.saved)
	}
	if !repository.saved.Enabled || !repository.saved.AtAll {
		t.Fatalf("saved = %+v, want enabled and at_all", repository.saved)
	}
}

func TestDingTalkServiceUpdateCanClearSecret(t *testing.T) {
	repository := &fakeRepository{setting: model.DingTalk{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=existing-token",
		Secret:  "SEC-existing",
		Key:     model.DefaultDingTalkConfigKey,
	}}
	service := NewDingTalkService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})
	if err := service.Update(context.Background(), schema.DingTalkUpdateArgs{ClearSecret: true}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.saved.Secret != "" {
		t.Fatalf("saved secret = %q, want empty", repository.saved.Secret)
	}
}

func TestDingTalkServiceRejectsSSRFWebhook(t *testing.T) {
	repository := &fakeRepository{err: gorm.ErrRecordNotFound}
	service := NewDingTalkService(repository, testutil.NewFakeAlarmRepo(), &fakeSender{})
	err := service.Update(context.Background(), schema.DingTalkUpdateArgs{
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

func TestDingTalkServiceTestUsesStoredCredentialsWhenDisabled(t *testing.T) {
	repository := &fakeRepository{setting: model.DingTalk{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=existing-token",
		Secret:  "SEC-existing",
		AtAll:   true,
		Enabled: false,
	}}
	sender := &fakeSender{}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{MessageTitle: "  生产环境主机告警  "}, nil
	}
	service := NewDingTalkService(repository, alarmRepository, sender)
	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || sender.configs[0].Secret != "SEC-existing" || !sender.configs[0].AtAll {
		t.Fatalf("sender configs = %+v, want stored config", sender.configs)
	}
	if sender.configs[0].MessageTitle != "生产环境主机告警" {
		t.Fatalf("message title = %q, want shared custom title", sender.configs[0].MessageTitle)
	}
}

func TestDingTalkServiceTestPropagatesSenderError(t *testing.T) {
	repository := &fakeRepository{setting: model.DingTalk{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=existing-token",
	}}
	sender := &fakeSender{err: errors.New("send failed")}
	service := NewDingTalkService(repository, testutil.NewFakeAlarmRepo(), sender)
	if err := service.Test(context.Background()); !errors.Is(err, sender.err) {
		t.Fatalf("Test() error = %v, want sender error", err)
	}
}

func TestDingTalkServiceTestFallsBackToDefaultTitle(t *testing.T) {
	repository := &fakeRepository{setting: model.DingTalk{
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=existing-token",
	}}
	alarmRepository := testutil.NewFakeAlarmRepo()
	alarmRepository.AlarmMessageSettingQueryFn = func(context.Context) (model.AlarmMessageSetting, error) {
		return model.AlarmMessageSetting{}, errors.New("database unavailable")
	}
	sender := &fakeSender{}
	service := NewDingTalkService(repository, alarmRepository, sender)

	if err := service.Test(context.Background()); err != nil {
		t.Fatalf("Test() error = %v", err)
	}
	if len(sender.configs) != 1 || sender.configs[0].MessageTitle != model.DefaultAlarmMessageTitle {
		t.Fatalf("sender configs = %+v, want default title", sender.configs)
	}
}
