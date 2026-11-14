package service

import (
	"context"
	"errors"
	"strings"

	apperrors "beacon/pkg/errors"
	alarmrepository "beacon/service/alarm/repository"
	alarmtitle "beacon/service/alarm/title"
	"beacon/service/model"
	"beacon/service/schema"
	"beacon/service/wecom/client"
	"beacon/service/wecom/repository"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var WeComServiceSet = wire.NewSet(NewWeComService, wire.Bind(new(IWeComService), new(*WeComService)))

type IWeComService interface {
	Query(context.Context) (schema.WeComSetting, error)
	Update(context.Context, schema.WeComUpdateArgs) error
	Test(context.Context) error
}

type WeComService struct {
	repository      repository.IWeComRepository
	alarmRepository alarmrepository.IAlarmRepository
	sender          client.Sender
}

func NewWeComService(repository repository.IWeComRepository, alarmRepository alarmrepository.IAlarmRepository, sender client.Sender) *WeComService {
	return &WeComService{repository: repository, alarmRepository: alarmRepository, sender: sender}
}

func (s *WeComService) Query(ctx context.Context) (schema.WeComSetting, error) {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return schema.WeComSetting{}, nil
	}
	if err != nil {
		return schema.WeComSetting{}, err
	}
	webhookMasked := ""
	if setting.Webhook != "" {
		webhookMasked = client.MaskWebhookURL(setting.Webhook)
	}
	return schema.WeComSetting{
		ID:                setting.ID,
		Enabled:           setting.Enabled,
		WebhookMasked:     webhookMasked,
		WebhookConfigured: setting.Webhook != "",
		AtAll:             setting.AtAll,
	}, nil
}

func (s *WeComService) Update(ctx context.Context, args schema.WeComUpdateArgs) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		setting = model.WeCom{Key: model.DefaultWeComConfigKey}
	} else if err != nil {
		return err
	}

	if webhook := strings.TrimSpace(args.Webhook); webhook != "" {
		if _, err := client.ValidateWebhookURL(webhook); err != nil {
			return apperrors.New400Error(err.Error())
		}
		setting.Webhook = webhook
	}
	setting.Enabled = args.Enabled
	setting.AtAll = args.AtAll

	if setting.Enabled && setting.Webhook == "" {
		return apperrors.New400Error("启用企业微信告警前必须配置 Webhook")
	}
	if setting.Webhook != "" {
		if _, err := client.ValidateWebhookURL(setting.Webhook); err != nil {
			return apperrors.New400Error(err.Error())
		}
	}
	return s.repository.Save(ctx, &setting)
}

func (s *WeComService) Test(ctx context.Context) error {
	setting, err := s.repository.Query(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.New400Error("请先配置企业微信 Webhook")
	}
	if err != nil {
		return err
	}
	if setting.Webhook == "" {
		return apperrors.New400Error("请先配置企业微信 Webhook")
	}
	messageSetting, messageSettingErr := s.alarmRepository.AlarmMessageSettingQuery(ctx)
	messageTitle := ""
	if messageSettingErr == nil {
		messageTitle = messageSetting.MessageTitle
	}
	return s.sender.Send(ctx, client.Config{
		Webhook:      setting.Webhook,
		AtAll:        setting.AtAll,
		MessageTitle: alarmtitle.OrDefault(messageTitle),
	}, "这是一条 Beacon 企业微信告警测试消息")
}
