// Package database tests for Ping().
package database

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newMemoryDB(t *testing.T) *DB {
	t.Helper()
	raw, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open memory sqlite: %v", err)
	}
	return &DB{raw}
}

func TestDB_Ping_Success(t *testing.T) {
	db := newMemoryDB(t)
	if err := db.Ping(); err != nil {
		t.Fatalf("Ping on healthy DB should succeed, got %v", err)
	}
}

func TestDB_Ping_NotInitialized(t *testing.T) {
	var nilDB *DB
	if err := nilDB.Ping(); err == nil {
		t.Fatal("Ping on nil DB should return an error")
	} else if !errorIsNotInitialized(err) {
		t.Fatalf("expected 'database not initialized', got %v", err)
	}

	// Zero-value struct (non-nil pointer but nil gorm.DB).
	z := &DB{}
	if err := z.Ping(); err == nil {
		t.Fatal("Ping on zero-value DB should return an error")
	} else if !errorIsNotInitialized(err) {
		t.Fatalf("expected 'database not initialized', got %v", err)
	}
}

func errorIsNotInitialized(err error) bool {
	return err != nil && err.Error() == "database not initialized"
}

// Silence unused errors import when only used via errors.New in db.go for nil-check.
var _ = errors.Is
