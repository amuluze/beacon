package repository

import (
	"context"
	"path/filepath"
	"testing"

	"beacon/service/model"
	"common/database"
)

func TestDingTalkRepositoryPersistsSingletonConfiguration(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "dingtalk")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(new(model.DingTalk)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repository := NewDingTalkRepository(db)
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

	got, err := repository.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if got.ID == 0 || !got.Enabled || !got.AtAll || got.Webhook != setting.Webhook || got.Secret != setting.Secret {
		t.Fatalf("Query() = %+v, want persisted configuration", got)
	}
}
