package service

import (
	"context"
	"errors"
	"strings"

	apperrors "beacon/pkg/errors"
	"beacon/service/dingtalk/client"
	"beacon/service/dingtalk/repository"
	"beacon/service/model"
	"beacon/service/schema"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var DingTalkServiceSet = wire.NewSet(NewDingTalkService, wire.Bind(new(IDingTalkService), new(*DingTalkService)))

type IDingTalkService interface {
	Query(context.Context) (schema.DingTalkSetting, error)
	Update(context.Context, schema.DingTalkUpdateArgs) error
	Test(context.Context) error
}

type DingTalkService struct {
	repository repository.IDingTalkRepository
	sender     client.Sender
}

func NewDingTalkService(repository repository.IDingTalkRepository, sender client.Sender) *DingTalkService {
	return &DingTalkService{repository: repository, sender: sender}
}

func (s *DingTalkService) Query(ctx context.Context) (schema.DingTalkSetting, error) {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return schema.DingTalkSetting{}, nil
	}
	if err != nil {
		return schema.DingTalkSetting{}, err
	}
	webhookMasked := ""
	if setting.Webhook != "" {
		webhookMasked = client.MaskWebhookURL(setting.Webhook)
	}
	return schema.DingTalkSetting{
		ID:                setting.ID,
		Enabled:           setting.Enabled,
		WebhookMasked:     webhookMasked,
		WebhookConfigured: setting.Webhook != "",
		SecretConfigured:  setting.Secret != "",
		AtAll:             setting.AtAll,
	}, nil
}

func (s *DingTalkService) Update(ctx context.Context, args schema.DingTalkUpdateArgs) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting = model.DingTalk{Key: model.DefaultDingTalkConfigKey}
	} else if err != nil {
		return err
	}

	if webhook := strings.TrimSpace(args.Webhook); webhook != "" {
		if _, err := client.ValidateWebhookURL(webhook); err != nil {
			return apperrors.New400Error(err.Error())
		}
		setting.Webhook = webhook
	}
	if args.ClearSecret {
		setting.Secret = ""
	} else if secret := strings.TrimSpace(args.Secret); secret != "" {
		setting.Secret = secret
	}
	setting.Enabled = args.Enabled
	setting.AtAll = args.AtAll

	if setting.Enabled && setting.Webhook == "" {
		return apperrors.New400Error("启用钉钉告警前必须配置 Webhook")
	}
	if setting.Webhook != "" {
		if _, err := client.ValidateWebhookURL(setting.Webhook); err != nil {
			return apperrors.New400Error(err.Error())
		}
	}
	return s.repository.Save(ctx, &setting)
}

func (s *DingTalkService) Test(ctx context.Context) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.New400Error("请先配置钉钉 Webhook")
	}
	if err != nil {
		return err
	}
	if setting.Webhook == "" {
		return apperrors.New400Error("请先配置钉钉 Webhook")
	}
	return s.sender.Send(ctx, client.Config{
		Webhook: setting.Webhook,
		Secret:  setting.Secret,
		AtAll:   setting.AtAll,
	}, "这是一条 Beacon 钉钉告警测试消息")
}
