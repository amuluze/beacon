package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"beacon/pkg/utils/hash"
	"beacon/pkg/utils/uuid"
	authrepo "beacon/service/auth/repository"
	"beacon/service/model"
	"beacon/service/schema"
	"common/database"

	"gorm.io/gorm"
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

func newAccountRoleFixture(t *testing.T, db *database.DB, name string) *model.Role {
	t.Helper()
	role := &model.Role{ID: uuid.MustUUID(), Name: name, Status: 1}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("seed role %s: %v", name, err)
	}
	return role
}

func TestUserCreateBindsAdminRoleFlag(t *testing.T) {
	db, err := database.NewDB(database.WithType("sqlite"), database.WithDBName(filepath.Join(t.TempDir(), "create-roles.db")))
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	adminRole := newAccountRoleFixture(t, db, model.DefaultAdminRoleName)
	normalRole := newAccountRoleFixture(t, db, "普通用户")
	repo := NewAccountRepository(db)
	ctx := context.Background()

	adminUser, err := repo.UserCreate(ctx, schema.UserCreateArgs{Username: "flag-admin", Password: "Strong-Secret-9", Status: 1, RoleIDs: []string{adminRole.ID.String()}})
	if err != nil {
		t.Fatalf("create admin-flag user: %v", err)
	}
	if adminUser.IsAdmin != model.IsAdminYes {
		t.Fatalf("admin role binding: is_admin = %d, want %d", adminUser.IsAdmin, model.IsAdminYes)
	}

	normalUser, err := repo.UserCreate(ctx, schema.UserCreateArgs{Username: "flag-normal", Password: "Strong-Secret-9", Status: 1, RoleIDs: []string{normalRole.ID.String()}})
	if err != nil {
		t.Fatalf("create normal-flag user: %v", err)
	}
	if normalUser.IsAdmin != model.IsAdminNo {
		t.Fatalf("normal role binding: is_admin = %d, want %d", normalUser.IsAdmin, model.IsAdminNo)
	}
}

func TestUserUpdateRoleSemantics(t *testing.T) {
	db, err := database.NewDB(database.WithType("sqlite"), database.WithDBName(filepath.Join(t.TempDir(), "update-roles.db")))
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	adminRole := newAccountRoleFixture(t, db, model.DefaultAdminRoleName)
	normalRole := newAccountRoleFixture(t, db, "普通用户")
	repo := NewAccountRepository(db)
	ctx := context.Background()
	user, err := repo.UserCreate(ctx, schema.UserCreateArgs{Username: "role-user", Password: "Strong-Secret-9", Status: 1})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	loadUser := func() model.User {
		var stored model.User
		if err := db.Preload("Roles").First(&stored, "id = ?", user.ID).Error; err != nil {
			t.Fatalf("load user: %v", err)
		}
		return stored
	}

	// 1) 显式传入管理员角色：关联被替换且 is_admin=1。
	if _, err := repo.UserUpdate(ctx, schema.UserUpdateArgs{ID: user.ID.String(), RoleIDs: []string{adminRole.ID.String()}}); err != nil {
		t.Fatalf("assign admin role: %v", err)
	}
	stored := loadUser()
	if len(stored.Roles) != 1 || stored.Roles[0].ID != adminRole.ID {
		t.Fatalf("assign admin role: got roles %v", stored.Roles)
	}
	if stored.IsAdmin != model.IsAdminYes {
		t.Fatalf("assign admin role: is_admin = %d, want %d", stored.IsAdmin, model.IsAdminYes)
	}

	// 2) 换绑普通角色：关联被替换且 is_admin 降为 2。
	if _, err := repo.UserUpdate(ctx, schema.UserUpdateArgs{ID: user.ID.String(), RoleIDs: []string{normalRole.ID.String()}}); err != nil {
		t.Fatalf("replace with normal role: %v", err)
	}
	stored = loadUser()
	if len(stored.Roles) != 1 || stored.Roles[0].ID != normalRole.ID {
		t.Fatalf("replace roles: got roles %v", stored.Roles)
	}
	if stored.IsAdmin != model.IsAdminNo {
		t.Fatalf("replace roles: is_admin = %d, want %d", stored.IsAdmin, model.IsAdminNo)
	}

	// 3) 不传 RoleIDs：角色与 is_admin 保持现状。
	if _, err := repo.UserUpdate(ctx, schema.UserUpdateArgs{ID: user.ID.String(), Status: 1}); err != nil {
		t.Fatalf("update without roles: %v", err)
	}
	stored = loadUser()
	if len(stored.Roles) != 1 || stored.Roles[0].ID != normalRole.ID {
		t.Fatalf("update without roles: roles changed to %v", stored.Roles)
	}
	if stored.IsAdmin != model.IsAdminNo {
		t.Fatalf("update without roles: is_admin = %d, want %d", stored.IsAdmin, model.IsAdminNo)
	}
}

func TestRoleDeleteRemovesOnlyTargetRoles(t *testing.T) {
	db, err := database.NewDB(database.WithType("sqlite"), database.WithDBName(filepath.Join(t.TempDir(), "role-delete.db")))
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	defer db.Close()
	if err := db.AutoMigrate(&model.User{}, &model.Role{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	roleA := newAccountRoleFixture(t, db, "role-a")
	roleB := newAccountRoleFixture(t, db, "role-b")
	user := &model.User{ID: uuid.MustUUID(), Username: "del-role-user", Password: "x", Status: 1, Roles: model.Roles{roleA}}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	repo := NewAccountRepository(db)
	if err := repo.RoleDelete(context.Background(), schema.RoleDeleteArgs{IDs: []string{roleA.ID.String()}}); err != nil {
		t.Fatalf("delete role: %v", err)
	}

	var gone model.Role
	if err := db.First(&gone, "id = ?", roleA.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("role A should be deleted, got err = %v", err)
	}
	var kept model.Role
	if err := db.First(&kept, "id = ?", roleB.ID).Error; err != nil {
		t.Fatalf("role B must survive: %v", err)
	}
	var storedUser model.User
	if err := db.First(&storedUser, "id = ?", user.ID).Error; err != nil {
		t.Fatalf("user must survive role deletion: %v", err)
	}
}
