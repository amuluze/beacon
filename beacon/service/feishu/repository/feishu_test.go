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

func newFeishuRepositoryCipher(t *testing.T) *secret.Cipher {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}
	return cipher
}

func TestFeishuRepositoryPersistsSingletonConfiguration(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "feishu")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Feishu)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repository := NewFeishuRepository(db, newFeishuRepositoryCipher(t))
	setting := model.Feishu{
		Key:     model.DefaultFeishuConfigKey,
		Enabled: true,
		Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/secret-hook-id",
		Secret:  "feishu-sign-secret",
	}
	if err := repository.Save(context.Background(), &setting); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	var stored model.Feishu
	if err := db.First(&stored, setting.ID).Error; err != nil {
		t.Fatalf("query raw setting: %v", err)
	}
	if !secret.IsEncrypted(stored.Webhook) || !secret.IsEncrypted(stored.Secret) {
		t.Fatalf("credentials were stored as plaintext: webhook=%q secret=%q", stored.Webhook, stored.Secret)
	}

	got, err := repository.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.Webhook != setting.Webhook || got.Secret != setting.Secret || !got.Enabled {
		t.Fatalf("Query() = %+v, want decrypted original", got)
	}
}

func TestFeishuRepositoryQueryMigratesLegacyPlaintext(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "feishu-legacy")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Feishu)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	cipher := newFeishuRepositoryCipher(t)
	repository := NewFeishuRepository(db, cipher)
	ctx := context.Background()

	// 模拟历史明文行：Query 应在读取时迁移为密文并返回明文。
	plaintextWebhook := "https://open.feishu.cn/open-apis/bot/v2/hook/legacy-hook-id"
	plaintextSecret := "legacy-secret"
	if err := db.WithContext(ctx).Create(&model.Feishu{
		Key:     model.DefaultFeishuConfigKey,
		Enabled: true,
		Webhook: plaintextWebhook,
		Secret:  plaintextSecret,
	}).Error; err != nil {
		t.Fatalf("seed legacy row: %v", err)
	}

	got, err := repository.Query(ctx)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.Webhook != plaintextWebhook || got.Secret != plaintextSecret {
		t.Fatalf("Query() = %+v, want legacy plaintext", got)
	}

	var stored model.Feishu
	if err := db.WithContext(ctx).First(&stored, got.ID).Error; err != nil {
		t.Fatalf("query raw setting: %v", err)
	}
	if !secret.IsEncrypted(stored.Webhook) || !secret.IsEncrypted(stored.Secret) {
		t.Fatalf("legacy credentials were not migrated to ciphertext")
	}
}
