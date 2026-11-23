package service

import (
	"context"
	"encoding/base64"
	"errors"
	"path/filepath"
	"testing"

	"beacon/pkg/secret"
	"beacon/service/model"
	"common/database"
)

func TestNewDataCipherRequiresValidEnvironmentKey(t *testing.T) {
	t.Setenv(dataEncryptionKeyEnv, "")
	if _, err := NewDataCipher(&Config{}); !errors.Is(err, ErrDataEncryptionKeyRequired) {
		t.Fatalf("missing key error = %v, want ErrDataEncryptionKeyRequired", err)
	}
	t.Setenv(dataEncryptionKeyEnv, base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef")))
	if _, err := NewDataCipher(&Config{}); err != nil {
		t.Fatalf("valid key rejected: %v", err)
	}
}

// TestMigrateSensitiveDataEncryptsLegacyRowsAcrossChannels 锁定启动迁移覆盖面：
// SMTP 密码、钉钉 webhook/secret、飞书 webhook/secret、企微 webhook 的遗留
// 明文行必须在启动事务内统一加密，而不是等各自配置被读取时才迁移。
func TestMigrateSensitiveDataEncryptsLegacyRowsAcrossChannels(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "encrypt-migrate")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Mail), new(model.DingTalk), new(model.Feishu), new(model.WeCom)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}

	if err := db.Create(&model.Mail{Server: "smtp.example.com", Password: "legacy-mail-password"}).Error; err != nil {
		t.Fatalf("seed mail: %v", err)
	}
	if err := db.Create(&model.DingTalk{Key: model.DefaultDingTalkConfigKey, Webhook: "legacy-ding-webhook", Secret: "legacy-ding-secret"}).Error; err != nil {
		t.Fatalf("seed dingtalk: %v", err)
	}
	if err := db.Create(&model.Feishu{Key: model.DefaultFeishuConfigKey, Webhook: "legacy-feishu-webhook", Secret: "legacy-feishu-secret"}).Error; err != nil {
		t.Fatalf("seed feishu: %v", err)
	}
	if err := db.Create(&model.WeCom{Key: model.DefaultWeComConfigKey, Webhook: "legacy-wecom-webhook"}).Error; err != nil {
		t.Fatalf("seed wecom: %v", err)
	}

	if err := migrateSensitiveData(context.Background(), db, cipher); err != nil {
		t.Fatalf("migrateSensitiveData: %v", err)
	}

	var mail model.Mail
	if err := db.First(&mail).Error; err != nil {
		t.Fatalf("reload mail: %v", err)
	}
	if !secret.IsEncrypted(mail.Password) {
		t.Fatalf("mail password not migrated: %q", mail.Password)
	}
	var ding model.DingTalk
	if err := db.First(&ding).Error; err != nil {
		t.Fatalf("reload dingtalk: %v", err)
	}
	if !secret.IsEncrypted(ding.Webhook) || !secret.IsEncrypted(ding.Secret) {
		t.Fatalf("dingtalk credentials not migrated: webhook=%q secret=%q", ding.Webhook, ding.Secret)
	}
	var feishu model.Feishu
	if err := db.First(&feishu).Error; err != nil {
		t.Fatalf("reload feishu: %v", err)
	}
	if !secret.IsEncrypted(feishu.Webhook) || !secret.IsEncrypted(feishu.Secret) {
		t.Fatalf("feishu credentials not migrated: webhook=%q secret=%q", feishu.Webhook, feishu.Secret)
	}
	var wecom model.WeCom
	if err := db.First(&wecom).Error; err != nil {
		t.Fatalf("reload wecom: %v", err)
	}
	if !secret.IsEncrypted(wecom.Webhook) {
		t.Fatalf("wecom webhook not migrated: %q", wecom.Webhook)
	}
}
