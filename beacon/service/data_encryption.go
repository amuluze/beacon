package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"
	"gorm.io/gorm"
)

const dataEncryptionKeyEnv = "BEACON_DATA_ENCRYPTION_KEY"

var ErrDataEncryptionKeyRequired = errors.New("data encryption key is required")

func NewDataCipher(_ *Config) (*secret.Cipher, error) {
	encoded := os.Getenv(dataEncryptionKeyEnv)
	if encoded == "" {
		return nil, fmt.Errorf("%w: set %s to a base64 or hex encoded 32-byte key", ErrDataEncryptionKeyRequired, dataEncryptionKeyEnv)
	}
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("initialize data encryption: %w", err)
	}
	return cipher, nil
}

func migrateSensitiveData(ctx context.Context, db *database.DB, cipher *secret.Cipher) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var mails []model.Mail
		if err := tx.Find(&mails).Error; err != nil {
			return err
		}
		for i := range mails {
			if err := migrateEncryptedColumn(tx, cipher, &mails[i], "password", mails[i].Password); err != nil {
				return fmt.Errorf("migrate SMTP password: %w", err)
			}
		}

		var settings []model.DingTalk
		if err := tx.Find(&settings).Error; err != nil {
			return err
		}
		for i := range settings {
			if err := migrateEncryptedColumn(tx, cipher, &settings[i], "webhook", settings[i].Webhook); err != nil {
				return fmt.Errorf("migrate DingTalk webhook: %w", err)
			}
			if err := migrateEncryptedColumn(tx, cipher, &settings[i], "secret", settings[i].Secret); err != nil {
				return fmt.Errorf("migrate DingTalk secret: %w", err)
			}
		}
		return nil
	})
}

func migrateEncryptedColumn(tx *gorm.DB, cipher *secret.Cipher, target any, column, value string) error {
	plaintext, legacy, err := cipher.Decrypt(value)
	if err != nil || !legacy {
		return err
	}
	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		return err
	}
	return tx.Model(target).Update(column, encrypted).Error
}
