package repository

import (
	"context"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"

	"github.com/google/wire"
)

var FeishuRepositorySet = wire.NewSet(NewFeishuRepository, wire.Bind(new(IFeishuRepository), new(*FeishuRepository)))

type IFeishuRepository interface {
	Query(context.Context) (model.Feishu, error)
	Save(context.Context, *model.Feishu) error
}

type FeishuRepository struct {
	DB     *database.DB
	Cipher *secret.Cipher
}

func NewFeishuRepository(db *database.DB, cipher *secret.Cipher) *FeishuRepository {
	return &FeishuRepository{DB: db, Cipher: cipher}
}

func (r *FeishuRepository) Query(ctx context.Context) (model.Feishu, error) {
	var setting model.Feishu
	err := r.DB.WithContext(ctx).
		Where("key = ?", model.DefaultFeishuConfigKey).
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

func (r *FeishuRepository) Save(ctx context.Context, setting *model.Feishu) error {
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
