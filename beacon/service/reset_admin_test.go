package service

import (
	"errors"
	"path/filepath"
	"testing"

	"beacon/pkg/auth"
	"beacon/pkg/utils/hash"
	"beacon/pkg/utils/uuid"
	"beacon/service/model"
	"common/database"

	"golang.org/x/crypto/bcrypt"
)

// newResetAdminTestDB 构造带 admin 账号的临时库（不含 admin 通过 seedAdmin 控制）。
func newResetAdminTestDB(t *testing.T, seedAdmin bool) *database.DB {
	t.Helper()
	db, err := database.NewDB(database.WithDBName(filepath.Join(t.TempDir(), "reset-admin")))
	if err != nil {
		t.Fatalf("new db: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.User{}, &model.Audit{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if seedAdmin {
		original, err := hash.BcryptHash("Original-Pass-1!")
		if err != nil {
			t.Fatal(err)
		}
		admin := &model.User{
			ID:       uuid.MustUUID(),
			Username: "admin",
			Password: original,
			Status:   1,
		}
		if err := db.Create(admin).Error; err != nil {
			t.Fatalf("seed admin: %v", err)
		}
	}
	return db
}

func TestResetAdminPasswordGeneratesTemporaryPassword(t *testing.T) {
	db := newResetAdminTestDB(t, true)

	password, generated, err := ResetAdminPassword(db, "")
	if err != nil {
		t.Fatalf("ResetAdminPassword: %v", err)
	}
	if !generated {
		t.Fatal("empty password must produce a generated temporary password")
	}
	if err := auth.ValidatePassword(password); err != nil {
		t.Fatalf("generated password must satisfy the policy: %v", err)
	}

	var admin model.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		t.Fatalf("stored password does not verify against the generated one: %v", err)
	}
	if !admin.MustChangePassword {
		t.Fatal("reset-admin must force a password change at next login")
	}

	var audits []model.Audit
	if err := db.Where("operate = ?", "reset admin password").Find(&audits).Error; err != nil {
		t.Fatal(err)
	}
	if len(audits) != 1 || audits[0].Username != "system" || audits[0].Object != admin.ID.String() {
		t.Fatalf("audit row = %#v", audits)
	}
}

func TestResetAdminPasswordExplicitPassword(t *testing.T) {
	db := newResetAdminTestDB(t, true)
	const explicit = "New-Admin-Pass-2026!"

	password, generated, err := ResetAdminPassword(db, explicit)
	if err != nil {
		t.Fatalf("ResetAdminPassword: %v", err)
	}
	if generated || password != explicit {
		t.Fatalf("explicit password must be used as-is: generated=%v password=%q", generated, password)
	}
	var admin model.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(explicit)); err != nil {
		t.Fatalf("stored password does not verify against the explicit one: %v", err)
	}
}

func TestResetAdminPasswordRejectsWeakPassword(t *testing.T) {
	db := newResetAdminTestDB(t, true)
	before, err := hash.BcryptHash("Original-Pass-1!")
	if err != nil {
		t.Fatal(err)
	}
	// 用已知原口令覆盖，确保弱口令被拒时数据库完全未变。
	if err := db.Model(&model.User{}).Where("username = ?", "admin").Update("password", before).Error; err != nil {
		t.Fatal(err)
	}

	if _, _, err := ResetAdminPassword(db, "short"); !errors.Is(err, auth.ErrWeakPassword) {
		t.Fatalf("weak password error = %v, want ErrWeakPassword", err)
	}
	var admin model.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if admin.Password != before {
		t.Fatal("rejected reset must not modify the stored password")
	}
}

func TestResetAdminPasswordMissingAdmin(t *testing.T) {
	db := newResetAdminTestDB(t, false)

	if _, _, err := ResetAdminPassword(db, ""); !errors.Is(err, ErrAdminMissing) {
		t.Fatalf("error = %v, want ErrAdminMissing", err)
	}
}
