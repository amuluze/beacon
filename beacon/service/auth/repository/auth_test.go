package repository

import (
	"context"
	"path/filepath"
	"testing"

	"beacon/pkg/errors"
	"beacon/pkg/utils/hash"
	"beacon/pkg/utils/uuid"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
)

func newTestDB(t *testing.T) *database.DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(dbPath),
	)
	if err != nil {
		t.Fatalf("create test db failed: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate test db failed: %v", err)
	}
	return db
}

func createTestUser(t *testing.T, db *database.DB, username, password string) model.User {
	t.Helper()
	hashed, err := hash.BcryptHash(password)
	if err != nil {
		t.Fatalf("hash password failed: %v", err)
	}
	user := model.User{
		ID:       uuid.MustUUID(),
		Username: username,
		Password: hashed,
		Status:   1,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create test user failed: %v", err)
	}
	return user
}

func TestAuthRepoLogin(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	repo := NewAuthRepo(db)
	createTestUser(t, db, "login-alice", "secret123")

	cases := []struct {
		name     string
		args     schema.LoginArgs
		wantErr  bool
		errMatch string
	}{
		{
			name: "correct password",
			args: schema.LoginArgs{Username: "login-alice", Password: "secret123"},
		},
		{
			name:     "wrong password",
			args:     schema.LoginArgs{Username: "login-alice", Password: "wrong"},
			wantErr:  true,
			errMatch: "invalid password",
		},
		{
			name:     "user not found",
			args:     schema.LoginArgs{Username: "bob", Password: "secret123"},
			wantErr:  true,
			errMatch: "record not found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Login(context.Background(), tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.errMatch)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAuthRepoPassUpdate(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	repo := NewAuthRepo(db)
	createTestUser(t, db, "update-alice", "old-pass")

	ctx := context.Background()

	// Update password successfully.
	if err := repo.PassUpdate(ctx, schema.PasswordUpdateArgs{
		Username:    "update-alice",
		OldPassword: "old-pass",
		NewPassword: "new-pass",
	}); err != nil {
		t.Fatalf("update password failed: %v", err)
	}

	// Login with new password should succeed.
	if _, err := repo.Login(ctx, schema.LoginArgs{Username: "update-alice", Password: "new-pass"}); err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}

	// Login with old password should fail.
	if _, err := repo.Login(ctx, schema.LoginArgs{Username: "update-alice", Password: "old-pass"}); err == nil {
		t.Fatal("login with old password should fail")
	}

	// Wrong old password should fail.
	err := repo.PassUpdate(ctx, schema.PasswordUpdateArgs{
		Username:    "update-alice",
		OldPassword: "wrong-pass",
		NewPassword: "another-pass",
	})
	if err == nil {
		t.Fatal("update with wrong old password should fail")
	}
	serviceErr, ok := err.(errors.Error)
	if !ok || serviceErr.Status != 400 {
		t.Fatalf("wrong old password error = %#v, want status 400", err)
	}
}
