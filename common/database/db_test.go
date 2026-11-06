// Package database tests for Ping().
package database

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
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

// Parse with the same pgx driver as GORM: verify TLS modes and ensure password
// characters cannot inject a second connection option.
func TestPostgresSSLModeAndQuotedSettings(t *testing.T) {
	for _, mode := range []string{"", "disable", "require", "verify-full"} {
		t.Run(mode, func(t *testing.T) {
			opt := &option{Type: "postgres", Host: "127.0.0.1", Port: "5432", UserName: "test-user", Password: "space ' slash\\ sslmode=disable", DBName: "test db"}
			WithSSLMode(mode)(opt)
			driver := dial(opt).(*postgres.Dialector)
			config, err := pgx.ParseConfig(driver.Config.DSN)
			if err != nil {
				t.Fatal(err)
			}
			if config.Password != opt.Password || config.Database != opt.DBName {
				t.Fatal("connection values were not preserved")
			}
			if mode == "" || mode == "disable" {
				if config.TLSConfig != nil {
					t.Fatal("default/disable must preserve historical behavior")
				}
			} else if config.TLSConfig == nil {
				t.Fatal("explicit TLS mode was lost")
			}
			if mode == "verify-full" && config.TLSConfig.InsecureSkipVerify {
				t.Fatal("verify-full must verify the server")
			}
		})
	}
	opt := &option{Type: "postgres", Host: "127.0.0.1", Port: "5432", UserName: "test", DBName: "test", SSLMode: "require user=injected"}
	if _, err := pgx.ParseConfig(dial(opt).(*postgres.Dialector).Config.DSN); err == nil {
		t.Fatal("invalid SSL mode should fail, never become extra DSN settings")
	}
}
