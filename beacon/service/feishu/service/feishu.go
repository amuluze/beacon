package service

import (
	"context"
	"errors"
	"strings"

	apperrors "beacon/pkg/errors"
	alarmrepository "beacon/service/alarm/repository"
	alarmtitle "beacon/service/alarm/title"
	"beacon/service/feishu/client"
	"beacon/service/feishu/repository"
	"beacon/service/model"
	"beacon/service/schema"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var FeishuServiceSet = wire.NewSet(NewFeishuService, wire.Bind(new(IFeishuService), new(*FeishuService)))

type IFeishuService interface {
	Query(context.Context) (schema.FeishuSetting, error)
	Update(context.Context, schema.FeishuUpdateArgs) error
	Test(context.Context) error
}

type FeishuService struct {
	repository      repository.IFeishuRepository
	alarmRepository alarmrepository.IAlarmRepository
	sender          client.Sender
}

func NewFeishuService(repository repository.IFeishuRepository, alarmRepository alarmrepository.IAlarmRepository, sender client.Sender) *FeishuService {
	return &FeishuService{repository: repository, alarmRepository: alarmRepository, sender: sender}
}

func (s *FeishuService) Query(ctx context.Context) (schema.FeishuSetting, error) {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return schema.FeishuSetting{}, nil
	}
	if err != nil {
		return schema.FeishuSetting{}, err
	}
	webhookMasked := ""
	if setting.Webhook != "" {
		webhookMasked = client.MaskWebhookURL(setting.Webhook)
	}
	return schema.FeishuSetting{
		ID:                setting.ID,
		Enabled:           setting.Enabled,
		WebhookMasked:     webhookMasked,
		WebhookConfigured: setting.Webhook != "",
		SecretConfigured:  setting.Secret != "",
	}, nil
}

func (s *FeishuService) Update(ctx context.Context, args schema.FeishuUpdateArgs) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting = model.Feishu{Key: model.DefaultFeishuConfigKey}
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

	if setting.Enabled && setting.Webhook == "" {
		return apperrors.New400Error("启用飞书告警前必须配置 Webhook")
	}
	if setting.Webhook != "" {
		if _, err := client.ValidateWebhookURL(setting.Webhook); err != nil {
			return apperrors.New400Error(err.Error())
		}
	}
	return s.repository.Save(ctx, &setting)
}

func (s *FeishuService) Test(ctx context.Context) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.New400Error("请先配置飞书 Webhook")
	}
	if err != nil {
		return err
	}
	if setting.Webhook == "" {
		return apperrors.New400Error("请先配置飞书 Webhook")
	}
	messageSetting, messageSettingErr := s.alarmRepository.AlarmMessageSettingQuery(ctx)
	messageTitle := ""
	if messageSettingErr == nil {
		messageTitle = messageSetting.MessageTitle
	}
	return s.sender.Send(ctx, client.Config{
		Webhook:      setting.Webhook,
		Secret:       setting.Secret,
		MessageTitle: alarmtitle.OrDefault(messageTitle),
	}, "这是一条 Beacon 飞书告警测试消息")
}
