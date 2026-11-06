package repository

import (
	"context"
	"path/filepath"
	"testing"

	"beacon/pkg/utils/hash"
	authrepo "beacon/service/auth/repository"
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

func TestAccountStatusPersistence(t *testing.T) {
	db, err := database.NewDB(database.WithType("sqlite"), database.WithDBName(filepath.Join(t.TempDir(), "status.db")))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Resource{}); err != nil {
		t.Fatal(err)
	}
	repo := NewAccountRepository(db)
	ctx := context.Background()
	user, err := repo.UserCreate(ctx, schema.UserCreateArgs{Username: "disabled-user", Password: "Strong-Secret-9", Status: 2})
	if err != nil {
		t.Fatal(err)
	}
	role, err := repo.RoleCreate(ctx, schema.RoleCreateArgs{Name: "disabled-role", Status: 2})
	if err != nil {
		t.Fatal(err)
	}
	// Read the stored values after create, then preserve, enable, and disable again.
	for _, status := range []int{0, 1, 2} {
		if _, err := repo.UserUpdate(ctx, schema.UserUpdateArgs{ID: user.ID.String(), Status: status}); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.RoleUpdate(ctx, schema.RoleUpdateArgs{ID: role.ID.String(), Status: status}); err != nil {
			t.Fatal(err)
		}
		var storedUser model.User
		var storedRole model.Role
		if err := db.First(&storedUser, "id = ?", user.ID).Error; err != nil {
			t.Fatal(err)
		}
		if err := db.First(&storedRole, "id = ?", role.ID).Error; err != nil {
			t.Fatal(err)
		}
		want := status
		if want == 0 {
			want = 2
		}
		if storedUser.Status != want || storedRole.Status != want {
			t.Fatalf("status user=%d role=%d want=%d", storedUser.Status, storedRole.Status, want)
		}
		_, loginErr := authrepo.NewAuthRepo(db).Login(ctx, schema.LoginArgs{Username: "disabled-user", Password: "Strong-Secret-9"})
		if (loginErr == nil) != (want == 1) {
			t.Fatalf("status %d: unexpected login acceptance: %v", want, loginErr)
		}
	}
}
