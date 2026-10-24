package repository

import (
	"bytes"
	"context"
	"encoding/base64"
	"mime"
	stdmail "net/mail"
	"path/filepath"
	"testing"

	"beacon/pkg/secret"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
)

func newMailRepository(t *testing.T) (*MailRepository, *database.DB) {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "mail")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.Mail)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := secret.NewFromEncodedKey(encoded)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}
	return NewMailRepository(db, cipher), db
}

func TestMailRepositoryEncryptsAndDecryptsPassword(t *testing.T) {
	repository, db := newMailRepository(t)
	args := schema.MailCreateArgs{
		Server: "smtp.example.com", Port: 465, Sender: "sender@example.com",
		Password: "smtp-secret", Receiver: "receiver@example.com",
	}
	if err := repository.MailCreate(context.Background(), args); err != nil {
		t.Fatalf("MailCreate: %v", err)
	}
	var stored model.Mail
	if err := db.First(&stored).Error; err != nil {
		t.Fatalf("query raw mail: %v", err)
	}
	if !secret.IsEncrypted(stored.Password) {
		t.Fatalf("password stored in plaintext: %q", stored.Password)
	}
	got, err := repository.MailQuery(context.Background())
	if err != nil {
		t.Fatalf("MailQuery: %v", err)
	}
	if got.Password != args.Password {
		t.Fatalf("decrypted password = %q, want original", got.Password)
	}
}

func TestMailRepositoryMigratesPlaintextPassword(t *testing.T) {
	repository, db := newMailRepository(t)
	legacy := model.Mail{Server: "smtp.example.com", Password: "legacy-password"}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatalf("create legacy mail: %v", err)
	}
	if _, err := repository.MailQuery(context.Background()); err != nil {
		t.Fatalf("MailQuery: %v", err)
	}
	if err := db.First(&legacy, legacy.ID).Error; err != nil {
		t.Fatalf("reload mail: %v", err)
	}
	if !secret.IsEncrypted(legacy.Password) {
		t.Fatal("legacy plaintext password was not migrated")
	}
}

func TestBuildMailTestMessageUsesSharedAlarmTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{name: "custom", title: "生产环境主机告警", want: "生产环境主机告警"},
		{name: "plain passthrough", title: "服务器告警", want: "服务器告警"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := buildMailTestMessage(
				model.Mail{Sender: "beacon@example.com"},
				schema.MailTestArgs{Receiver: "admin@example.com"},
				tt.title,
			)
			var raw bytes.Buffer
			if _, err := message.WriteTo(&raw); err != nil {
				t.Fatalf("write message: %v", err)
			}
			parsed, err := stdmail.ReadMessage(&raw)
			if err != nil {
				t.Fatalf("parse message: %v", err)
			}
			subject, err := new(mime.WordDecoder).DecodeHeader(parsed.Header.Get("Subject"))
			if err != nil {
				t.Fatalf("decode subject: %v", err)
			}
			if subject != tt.want {
				t.Fatalf("subject = %q, want %q", subject, tt.want)
			}
		})
	}
}
