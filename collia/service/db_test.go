package service

import (
	"path/filepath"
	"strings"
	"testing"

	"collia/service/model"
)

// TestNewDBAppliesWALOnSQLite 验证 SQLite 后端仍启用 WAL 日志模式。
func TestNewDBAppliesWALOnSQLite(t *testing.T) {
	config := &Config{DB: DB{DBType: "sqlite", DBName: filepath.Join(t.TempDir(), "probe")}}
	db, err := NewDB(config, &model.Models{})
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(db.Close)

	var mode string
	if err := db.Raw("PRAGMA journal_mode").Scan(&mode).Error; err != nil {
		t.Fatalf("query journal mode: %v", err)
	}
	if !strings.EqualFold(mode, "wal") {
		t.Fatalf("journal mode = %q, want wal", mode)
	}
}

// TestNewDBSkipsSQLitePragmasOnNonSQLiteBackends 验证非 SQLite 后端的失败
// 原样来自数据库拨号错误，而非 SQLite 专属 PRAGMA 的语法错误。
func TestNewDBSkipsSQLitePragmasOnNonSQLiteBackends(t *testing.T) {
	config := &Config{DB: DB{
		DBType:   "mysql",
		Host:     "127.0.0.1",
		Port:     "1", // 不可达端口，拨号快速失败
		User:     "beacon",
		Password: "secret",
		DBName:   "beacon",
	}}
	_, err := NewDB(config, &model.Models{})
	if err == nil {
		t.Fatal("NewDB with unreachable mysql should fail")
	}
	if strings.Contains(err.Error(), "PRAGMA") || strings.Contains(err.Error(), "journal_mode") {
		t.Fatalf("non-sqlite backend must not execute sqlite PRAGMA, got: %v", err)
	}
}
