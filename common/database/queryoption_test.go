package database

import "testing"

// TestOptionDBSkipsNilOptions 个别 QueryOption 工厂（WithOffset/WithLimit）在非法
// 入参时返回 nil；OptionDB 必须跳过 nil option 而不是 panic。
func TestOptionDBSkipsNilOptions(t *testing.T) {
	db := &DB{}
	if got := OptionDB(db, nil, nil); got != db {
		t.Fatalf("OptionDB must skip nil options and return the same DB")
	}
	if got := OptionDB(db, WithOffset(-1), WithLimit(0)); got != db {
		t.Fatalf("OptionDB must skip nil-valued options, got %v", got)
	}
}
