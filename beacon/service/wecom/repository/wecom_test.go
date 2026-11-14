package repository

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"testing"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"
)

func newWeComRepositoryCipher(t *testing.T) *secret.Cipher {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}
	return cipher
}

func TestWeComRepositoryPersistsSingletonConfiguration(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "wecom")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.WeCom)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repository := NewWeComRepository(db, newWeComRepositoryCipher(t))
	setting := model.WeCom{
		Key:     model.DefaultWeComConfigKey,
		Enabled: true,
		Webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=secret-key-value",
		AtAll:   true,
	}
	if err := repository.Save(context.Background(), &setting); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	var stored model.WeCom
	if err := db.First(&stored, setting.ID).Error; err != nil {
		t.Fatalf("query raw setting: %v", err)
	}
	if !secret.IsEncrypted(stored.Webhook) {
		t.Fatalf("credentials were stored as plaintext: webhook=%q", stored.Webhook)
	}

	got, err := repository.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.Webhook != setting.Webhook || !got.Enabled || !got.AtAll {
		t.Fatalf("Query() = %+v, want decrypted original", got)
	}
}

func TestWeComRepositoryQueryMigratesLegacyPlaintext(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "wecom-legacy")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.WeCom)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	cipher := newWeComRepositoryCipher(t)
	repository := NewWeComRepository(db, cipher)
	ctx := context.Background()

	plaintextWebhook := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=legacy-key-value"
	if err := db.WithContext(ctx).Create(&model.WeCom{
		Key:     model.DefaultWeComConfigKey,
		Enabled: true,
		Webhook: plaintextWebhook,
	}).Error; err != nil {
		t.Fatalf("seed legacy row: %v", err)
	}

	got, err := repository.Query(ctx)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.Webhook != plaintextWebhook {
		t.Fatalf("Query() = %+v, want legacy plaintext", got)
	}

	var stored model.WeCom
	if err := db.WithContext(ctx).First(&stored, got.ID).Error; err != nil {
		t.Fatalf("query raw setting: %v", err)
	}
	if !secret.IsEncrypted(stored.Webhook) {
		t.Fatalf("legacy credentials were not migrated to ciphertext")
	}
}
