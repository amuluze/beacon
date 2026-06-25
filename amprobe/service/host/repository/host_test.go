package repository

import (
	"context"
	"path/filepath"
	"testing"

	"amprobe/service/model"
	"common/database"
	rpcSchema "common/rpc/schema"
)

func TestNetUsageReturnsDBError(t *testing.T) {
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "probe")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	if err := db.AutoMigrate(new(model.MonitorNet)); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	db.Close()

	repo := &HostRepo{DB: db}
	if _, err := repo.NetUsage(context.Background(), rpcSchema.NetUsageArgs{}); err == nil {
		t.Fatal("expected db error")
	}
}
