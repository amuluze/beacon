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

func newRepositoryCipher(t *testing.T) *secret.Cipher {
	t.Helper()
	encoded := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}
	return cipher
}

func TestDingTalkRepositoryPersistsSingletonConfiguration(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "dingtalk")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.DingTalk)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repository := NewDingTalkRepository(db, newRepositoryCipher(t))
	setting := model.DingTalk{
		Key:     model.DefaultDingTalkConfigKey,
		Enabled: true,
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=secret-token",
		Secret:  "SEC-secret",
		AtAll:   true,
	}
	if err := repository.Save(context.Background(), &setting); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	var stored model.DingTalk
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
	if got.ID == 0 || !got.Enabled || !got.AtAll || got.Webhook != setting.Webhook || got.Secret != setting.Secret {
		t.Fatalf("Query() = %+v, want persisted configuration", got)
	}
}

func TestDingTalkRepositoryTransparentlyMigratesPlaintext(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "dingtalk-migration")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.DingTalk)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	legacy := model.DingTalk{
		Key:     model.DefaultDingTalkConfigKey,
		Webhook: "https://oapi.dingtalk.com/robot/send?access_token=legacy-token",
		Secret:  "legacy-secret",
	}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatalf("create plaintext setting: %v", err)
	}
	repository := NewDingTalkRepository(db, newRepositoryCipher(t))
	got, err := repository.Query(context.Background())
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if got.Webhook != legacy.Webhook || got.Secret != legacy.Secret {
		t.Fatalf("decrypted setting = %+v", got)
	}
	var stored model.DingTalk
	if err := db.First(&stored, legacy.ID).Error; err != nil {
		t.Fatalf("query migrated setting: %v", err)
	}
	if !secret.IsEncrypted(stored.Webhook) || !secret.IsEncrypted(stored.Secret) {
		t.Fatal("legacy plaintext setting was not migrated")
	}
}
