package database

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestDebugSQLLogsDoNotInterpolateParameters(t *testing.T) {
	var logs bytes.Buffer
	db, err := NewDB(
		WithType("sqlite"),
		WithDBName(filepath.Join(t.TempDir(), "parameterized")),
		WithDebug(true),
		WithLogWriter(&logs),
	)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(db.Close)

	type secretRecord struct {
		ID     uint `gorm:"primaryKey"`
		Secret string
	}
	if err := db.AutoMigrate(&secretRecord{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	const secret = "do-not-log-this-secret"
	if err := db.Create(&secretRecord{Secret: secret}).Error; err != nil {
		t.Fatalf("Create: %v", err)
	}
	if strings.Contains(logs.String(), secret) {
		t.Fatalf("parameterized SQL log leaked query value: %s", logs.String())
	}
}
