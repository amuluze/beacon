package repository

import (
	"context"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"

	"github.com/google/wire"
)

var WeComRepositorySet = wire.NewSet(NewWeComRepository, wire.Bind(new(IWeComRepository), new(*WeComRepository)))

type IWeComRepository interface {
	Query(context.Context) (model.WeCom, error)
	Save(context.Context, *model.WeCom) error
}

type WeComRepository struct {
	DB     *database.DB
	Cipher *secret.Cipher
}

func NewWeComRepository(db *database.DB, cipher *secret.Cipher) *WeComRepository {
	return &WeComRepository{DB: db, Cipher: cipher}
}

func (r *WeComRepository) Query(ctx context.Context) (model.WeCom, error) {
	var setting model.WeCom
	err := r.DB.WithContext(ctx).
		Where("key = ?", model.DefaultWeComConfigKey).
		First(&setting).Error
	if err != nil {
		return setting, err
	}
	webhook, webhookLegacy, err := r.Cipher.Decrypt(setting.Webhook)
	if err != nil {
		return setting, err
	}
	if webhookLegacy {
		encrypted, encErr := r.Cipher.Encrypt(webhook)
		if encErr != nil {
			return setting, encErr
		}
		if err := r.DB.WithContext(ctx).Model(&setting).Update("webhook", encrypted).Error; err != nil {
			return setting, err
		}
	}
	setting.Webhook = webhook
	return setting, nil
}

func (r *WeComRepository) Save(ctx context.Context, setting *model.WeCom) error {
	stored := *setting
	var err error
	stored.Webhook, err = r.Cipher.Encrypt(setting.Webhook)
	if err != nil {
		return err
	}
	if err := r.DB.WithContext(ctx).Save(&stored).Error; err != nil {
		return err
	}
	setting.Model = stored.Model
	return nil
}
