package repository

import (
	"context"
	"path/filepath"
	"testing"

	"beacon/pkg/utils/hash"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"
)

func TestUserCreateStoresBcryptPassword(t *testing.T) {
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "account.db")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	const password = "Strong-Secret-9"
	user, err := NewAccountRepository(db).UserCreate(context.Background(), schema.UserCreateArgs{
		Username: "bcrypt-user",
		Password: password,
		Status:   1,
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.Password == password {
		t.Fatal("password was stored as plaintext")
	}
	if err := hash.BcryptVerify(password, user.Password); err != nil {
		t.Fatalf("stored password is not a bcrypt hash of the submitted password: %v", err)
	}
}
