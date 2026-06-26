// Package jwtauth
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for token store
package jwtauth

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func newCache() *cache.Cache {
	return cache.New(5*time.Minute, 60*time.Second)
}

func TestStore_SetAndCheck(t *testing.T) {
	s := &Store{
		Storage: newCache(),
		Prefix:  "test:",
	}

	tests := []struct {
		name    string
		key     string
		set     bool
		want    bool
		wantErr bool
	}{
		{name: "set then check", key: "token-abc", set: true, want: true, wantErr: false},
		{name: "check without set", key: "token-missing", set: false, want: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				if err := s.Set(tt.key, 5*time.Minute); err != nil {
					t.Fatalf("Set failed: %v", err)
				}
			}
			found, err := s.Check(tt.key)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("Check error: %v", err)
			}
			if found != tt.want {
				t.Errorf("Check(%q) = %v, want %v", tt.key, found, tt.want)
			}
		})
	}
}

func TestStore_PrefixIsolation(t *testing.T) {
	s1 := &Store{Storage: newCache(), Prefix: "app1:"}
	s2 := &Store{Storage: newCache(), Prefix: "app2:"}

	if err := s1.Set("token-x", 5*time.Minute); err != nil {
		t.Fatalf("s1.Set failed: %v", err)
	}

	found1, _ := s1.Check("token-x")
	found2, _ := s2.Check("token-x")

	if !found1 {
		t.Error("s1 should find token-x")
	}
	if found2 {
		t.Error("s2 should NOT find token-x (different prefix)")
	}
}

func TestStore_Close(t *testing.T) {
	s := &Store{Storage: newCache(), Prefix: "test:"}
	if err := s.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
