package repository

import (
	"context"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"

	"github.com/google/wire"
)

var DingTalkRepositorySet = wire.NewSet(NewDingTalkRepository, wire.Bind(new(IDingTalkRepository), new(*DingTalkRepository)))

type IDingTalkRepository interface {
	Query(context.Context) (model.DingTalk, error)
	Save(context.Context, *model.DingTalk) error
}

type DingTalkRepository struct {
	DB     *database.DB
	Cipher *secret.Cipher
}

func NewDingTalkRepository(db *database.DB, cipher *secret.Cipher) *DingTalkRepository {
	return &DingTalkRepository{DB: db, Cipher: cipher}
}

func (r *DingTalkRepository) Query(ctx context.Context) (model.DingTalk, error) {
	var setting model.DingTalk
	err := r.DB.WithContext(ctx).
		Where("key = ?", model.DefaultDingTalkConfigKey).
		First(&setting).Error
	if err != nil {
		return setting, err
	}
	webhook, webhookLegacy, err := r.Cipher.Decrypt(setting.Webhook)
	if err != nil {
		return setting, err
	}
	secretValue, secretLegacy, err := r.Cipher.Decrypt(setting.Secret)
	if err != nil {
		return setting, err
	}
	updates := make(map[string]any)
	if webhookLegacy {
		updates["webhook"], err = r.Cipher.Encrypt(webhook)
		if err != nil {
			return setting, err
		}
	}
	if secretLegacy {
		updates["secret"], err = r.Cipher.Encrypt(secretValue)
		if err != nil {
			return setting, err
		}
	}
	if len(updates) > 0 {
		if err := r.DB.WithContext(ctx).Model(&setting).Updates(updates).Error; err != nil {
			return setting, err
		}
	}
	setting.Webhook = webhook
	setting.Secret = secretValue
	return setting, nil
}

func (r *DingTalkRepository) Save(ctx context.Context, setting *model.DingTalk) error {
	stored := *setting
	var err error
	stored.Webhook, err = r.Cipher.Encrypt(setting.Webhook)
	if err != nil {
		return err
	}
	stored.Secret, err = r.Cipher.Encrypt(setting.Secret)
	if err != nil {
		return err
	}
	if err := r.DB.WithContext(ctx).Save(&stored).Error; err != nil {
		return err
	}
	setting.Model = stored.Model
	return nil
}
