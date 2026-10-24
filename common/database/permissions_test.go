package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSQLiteDatabaseUsesOwnerOnlyPermissions(t *testing.T) {
	dbName := filepath.Join(t.TempDir(), "storage", "beacon")
	db, err := NewDB(WithType("sqlite"), WithDBName(dbName))
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(db.Close)

	info, err := os.Stat(sqliteFilePath(dbName))
	if err != nil {
		t.Fatalf("stat sqlite db: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("sqlite permissions = %o, want 600", got)
	}
}

func TestSecureSQLiteFilesRestrictsWALAndSHM(t *testing.T) {
	dbName := filepath.Join(t.TempDir(), "beacon")
	base := sqliteFilePath(dbName)
	for _, path := range []string{base, base + "-wal", base + "-shm"} {
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("write %s: %v", filepath.Base(path), err)
		}
		if err := os.Chmod(path, 0o644); err != nil {
			t.Fatalf("chmod %s: %v", filepath.Base(path), err)
		}
	}

	if err := SecureSQLiteFiles(dbName); err != nil {
		t.Fatalf("SecureSQLiteFiles: %v", err)
	}

	for _, path := range []string{base, base + "-wal", base + "-shm"} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", filepath.Base(path), err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("%s permissions = %o, want 600", filepath.Base(path), got)
		}
	}
}
