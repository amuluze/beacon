// Package database
// Date: 2026/07/16
// Author: Amu
// Description: tests for query options
package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newOptionDB(t *testing.T) *DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite 失败: %v", err)
	}
	return &DB{gdb}
}

// OptionDB 在传入 nil 或 no-op 选项时不应 panic。
func TestOptionDB_NilOptionSafe(t *testing.T) {
	db := newOptionDB(t)
	got := OptionDB(db, WithLimit(0), WithOffset(-1), nil)
	if got == nil {
		t.Fatal("OptionDB 不应返回 nil")
	}
}

// WithOffset / WithLimit 在非法入参时应返回 no-op（非 nil），避免下游 panic。
func TestWithOffsetAndLimit_NoNil(t *testing.T) {
	db := newOptionDB(t)
	if got := WithOffset(-1)(db); got == nil {
		t.Fatal("WithOffset(-1) 不应返回 nil")
	}
	if got := WithLimit(0)(db); got == nil {
		t.Fatal("WithLimit(0) 不应返回 nil")
	}
}
