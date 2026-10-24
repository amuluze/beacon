package service

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"beacon/pkg/utils/hash"
	"beacon/pkg/utils/uuid"
	"beacon/service/model"
	"common/database"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func TestInitAlarmThresholdCreatesAllDefaultsAndPreservesOverrides(t *testing.T) {
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.AlarmThreshold{}); err != nil {
		t.Fatalf("migrate alarm threshold: %v", err)
	}

	prepare := &Prepare{db: db}
	if err := prepare.InitAlarmThreshold(); err != nil {
		t.Fatalf("initialize alarm thresholds: %v", err)
	}

	var got []model.AlarmThreshold
	if err := db.Order("type").Find(&got).Error; err != nil {
		t.Fatalf("query alarm thresholds: %v", err)
	}
	if len(got) != len(thresholds) {
		t.Fatalf("alarm threshold count = %d, want %d", len(got), len(thresholds))
	}

	if err := db.Model(&model.AlarmThreshold{}).
		Where("type = ?", "memory").
		Updates(map[string]any{"duration": 5, "threshold": 90}).Error; err != nil {
		t.Fatalf("override memory alarm threshold: %v", err)
	}
	if err := prepare.InitAlarmThreshold(); err != nil {
		t.Fatalf("reinitialize alarm thresholds: %v", err)
	}

	var count int64
	if err := db.Model(&model.AlarmThreshold{}).Count(&count).Error; err != nil {
		t.Fatalf("count alarm thresholds: %v", err)
	}
	if count != int64(len(thresholds)) {
		t.Fatalf("alarm threshold count after reinitialization = %d, want %d", count, len(thresholds))
	}

	var memory model.AlarmThreshold
	if err := db.Where("type = ?", "memory").First(&memory).Error; err != nil {
		t.Fatalf("query memory alarm threshold: %v", err)
	}
	if memory.Duration != 5 || memory.Threshold != 90 {
		t.Fatalf("memory override was replaced: duration=%d threshold=%d", memory.Duration, memory.Threshold)
	}
}

func TestInitAccountRequiresBootstrapPasswordForNewDatabase(t *testing.T) {
	db := newPrepareTestDB(t)
	prepare := &Prepare{db: db}

	t.Setenv(bootstrapAdminPasswordEnv, "")
	err := prepare.InitAccount(fiber.New())
	if !errors.Is(err, ErrBootstrapAdminPasswordRequired) {
		t.Fatalf("InitAccount() error = %v, want ErrBootstrapAdminPasswordRequired", err)
	}
}

func TestInitAccountRejectsWeakBootstrapPassword(t *testing.T) {
	db := newPrepareTestDB(t)
	prepare := &Prepare{db: db}

	t.Setenv(bootstrapAdminPasswordEnv, "too-short")
	err := prepare.InitAccount(fiber.New())
	if !errors.Is(err, ErrBootstrapAdminPasswordWeak) {
		t.Fatalf("InitAccount() error = %v, want ErrBootstrapAdminPasswordWeak", err)
	}
}

func TestInitAccountCreatesOnlyMandatoryChangeAdministrator(t *testing.T) {
	db := newPrepareTestDB(t)
	prepare := &Prepare{db: db}
	const bootstrapPassword = "Bootstrap-Secure-9!"
	t.Setenv(bootstrapAdminPasswordEnv, bootstrapPassword)

	if err := prepare.InitAccount(fiber.New()); err != nil {
		t.Fatalf("InitAccount() error = %v", err)
	}

	var admin model.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("query bootstrap administrator: %v", err)
	}
	if !admin.MustChangePassword {
		t.Fatal("bootstrap administrator must be required to change password")
	}
	if err := hash.BcryptVerify(bootstrapPassword, admin.Password); err != nil {
		t.Fatalf("bootstrap password does not match stored hash: %v", err)
	}

	var normalCount int64
	if err := db.Model(&model.User{}).Where("username = ?", "beacon").Count(&normalCount).Error; err != nil {
		t.Fatalf("count removed default user: %v", err)
	}
	if normalCount != 0 {
		t.Fatalf("default normal user count = %d, want 0", normalCount)
	}
}

func newPrepareTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Resource{}); err != nil {
		t.Fatalf("migrate account models: %v", err)
	}
	return db
}

func TestInitAccountMarksHistoricalDefaultForMandatoryChange(t *testing.T) {
	db := newPrepareTestDB(t)
	var legacyHash, migratedHash string
	for legacyHash, migratedHash = range legacyPasswordMigrations {
		break
	}
	admin := model.User{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: legacyHash,
		Status:   1,
		IsAdmin:  1,
	}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create legacy administrator: %v", err)
	}

	if err := (&Prepare{db: db}).InitAccount(fiber.New()); err != nil {
		t.Fatalf("InitAccount() error = %v", err)
	}
	if err := db.First(&admin, "id = ?", admin.ID).Error; err != nil {
		t.Fatalf("reload migrated administrator: %v", err)
	}
	if admin.Password != migratedHash {
		t.Fatalf("migrated password hash = %q, want historical bcrypt hash", admin.Password)
	}
	if !admin.MustChangePassword {
		t.Fatal("historical default password was not marked for mandatory change")
	}
}

func TestInitAccountSyncsNewRoutesToExistingAdminRole(t *testing.T) {
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Resource{}); err != nil {
		t.Fatalf("migrate account models: %v", err)
	}
	existingResources := model.Resources{
		{ID: uuid.MustUUID(), Name: "", Path: "/health", Method: fiber.MethodGet, Status: 1},
		{ID: uuid.MustUUID(), Name: "已有查询接口", Path: "/api/v1/existing", Method: fiber.MethodGet, Status: 1},
		{ID: uuid.MustUUID(), Name: "查询钉钉告警配置", Path: "/api/v1/dingtalk/query", Method: fiber.MethodGet, Status: 1},
		{ID: uuid.MustUUID(), Name: "更新钉钉告警配置", Path: "/api/v1/dingtalk/update", Method: fiber.MethodPost, Status: 1},
		{ID: uuid.MustUUID(), Name: "测试钉钉告警", Path: "/api/v1/dingtalk/test", Method: fiber.MethodPost, Status: 1},
	}
	if err := db.Create(&existingResources).Error; err != nil {
		t.Fatalf("create existing route resources: %v", err)
	}

	adminRole := &model.Role{ID: uuid.MustUUID(), Name: "管理员", Status: 1}
	normalRole := &model.Role{ID: uuid.MustUUID(), Name: "普通用户", Status: 1}
	adminPassword, err := hash.BcryptHash("Existing-Admin-9!")
	if err != nil {
		t.Fatalf("hash existing admin password: %v", err)
	}
	normalPassword, err := hash.BcryptHash("Existing-Normal-9!")
	if err != nil {
		t.Fatalf("hash existing normal password: %v", err)
	}
	admin := &model.User{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: adminPassword,
		IsAdmin:  1,
		Status:   1,
		Roles:    model.Roles{adminRole},
	}
	normal := &model.User{
		ID:       uuid.MustUUID(),
		Username: "beacon",
		Password: normalPassword,
		IsAdmin:  2,
		Status:   1,
		Roles:    model.Roles{normalRole},
	}
	if err := db.Create(admin).Error; err != nil {
		t.Fatalf("create existing admin: %v", err)
	}
	if err := db.Create(normal).Error; err != nil {
		t.Fatalf("create existing normal user: %v", err)
	}

	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	app.Get("/api/v1/existing", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }).Name("已有查询接口")
	app.Get("/api/v1/dingtalk/query", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }).Name("查询钉钉告警配置")
	app.Post("/api/v1/dingtalk/update", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }).Name("更新钉钉告警配置")
	app.Post("/api/v1/dingtalk/test", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }).Name("测试钉钉告警")

	enforcer, err := casbin.NewSyncedEnforcer(filepath.Join("..", "configs", "model.conf"))
	if err != nil {
		t.Fatalf("create Casbin enforcer: %v", err)
	}
	if _, err := enforcer.AddNamedPolicy("p", normalRole.ID.String(), "/api/v1/dingtalk/query", fiber.MethodGet); err != nil {
		t.Fatalf("seed stale normal-user dingtalk policy: %v", err)
	}
	prepare := &Prepare{db: db, enforcer: enforcer}
	if err := prepare.InitAccount(app); err != nil {
		t.Fatalf("first InitAccount: %v", err)
	}
	if err := prepare.InitAccount(app); err != nil {
		t.Fatalf("second InitAccount: %v", err)
	}
	if err := prepare.InitCasbinRules(); err != nil {
		t.Fatalf("InitCasbinRules: %v", err)
	}

	for _, route := range []struct {
		path   string
		method string
	}{
		{path: "/api/v1/dingtalk/query", method: fiber.MethodGet},
		{path: "/api/v1/dingtalk/update", method: fiber.MethodPost},
		{path: "/api/v1/dingtalk/test", method: fiber.MethodPost},
	} {
		allowed, err := enforcer.Enforce(admin.ID.String(), route.path, route.method)
		if err != nil {
			t.Fatalf("enforce %s %s: %v", route.method, route.path, err)
		}
		if !allowed {
			t.Fatalf("existing admin is not allowed to access %s %s", route.method, route.path)
		}
		normalAllowed, err := enforcer.Enforce(normal.ID.String(), route.path, route.method)
		if err != nil {
			t.Fatalf("enforce normal user %s %s: %v", route.method, route.path, err)
		}
		if normalAllowed {
			t.Fatalf("normal user is unexpectedly allowed to access %s %s", route.method, route.path)
		}
	}
	policies, err := enforcer.GetNamedPolicy("p")
	if err != nil {
		t.Fatalf("query synchronized Casbin policies: %v", err)
	}
	dingTalkPolicyCount := 0
	for _, policy := range policies {
		if len(policy) >= 2 && strings.HasPrefix(policy[1], "/api/v1/dingtalk/") {
			dingTalkPolicyCount++
		}
	}
	if dingTalkPolicyCount != 3 {
		t.Fatalf("dingtalk Casbin policy count = %d, want 3", dingTalkPolicyCount)
	}

	var persistedRole model.Role
	if err := db.Where("name = ?", "管理员").Preload("Resources").First(&persistedRole).Error; err != nil {
		t.Fatalf("query administrator role resources: %v", err)
	}
	linkedRoutes := make(map[string]struct{}, len(persistedRole.Resources))
	for _, resource := range persistedRole.Resources {
		linkedRoutes[resource.Method+" "+resource.Path] = struct{}{}
	}
	for _, key := range []string{
		"GET /api/v1/dingtalk/query",
		"POST /api/v1/dingtalk/update",
		"POST /api/v1/dingtalk/test",
	} {
		if _, ok := linkedRoutes[key]; !ok {
			t.Fatalf("administrator role is missing resource %s", key)
		}
	}
}

func TestInitAlarmThresholdNormalizesRecordsAndPreservesOverrides(t *testing.T) {
	db, err := database.NewDB(
		database.WithType("sqlite"),
		database.WithDBName(filepath.Join(t.TempDir(), "beacon")),
	)
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	t.Cleanup(db.Close)
	if err := db.AutoMigrate(&model.AlarmThreshold{}); err != nil {
		t.Fatalf("migrate alarm threshold: %v", err)
	}

	seed := []model.AlarmThreshold{
		{Type: "cpu", Duration: 5, Threshold: 71},
		{Type: "cpu", Duration: 2, Threshold: 80},
		{Type: "mem", Duration: 4, Threshold: 73},
		{Type: "memory", Duration: 8, Threshold: 91},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed thresholds: %v", err)
	}

	prepare := &Prepare{db: db}
	if err := prepare.InitAlarmThreshold(); err != nil {
		t.Fatalf("first normalize alarm thresholds: %v", err)
	}
	if err := prepare.InitAlarmThreshold(); err != nil {
		t.Fatalf("second normalize alarm thresholds: %v", err)
	}

	var records []model.AlarmThreshold
	if err := db.Order("id ASC").Find(&records).Error; err != nil {
		t.Fatalf("query normalized thresholds: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("active threshold count = %d, want 3", len(records))
	}

	byType := make(map[string]model.AlarmThreshold, len(records))
	for _, record := range records {
		if _, exists := byType[record.Type]; exists {
			t.Fatalf("duplicate active threshold type %q", record.Type)
		}
		byType[record.Type] = record
	}

	if cpu := byType["cpu"]; cpu.Duration != 5 || cpu.Threshold != 71 {
		t.Fatalf("cpu threshold = %#v, want preserved user override", cpu)
	}
	if memory := byType["memory"]; memory.Duration != 4 || memory.Threshold != 73 {
		t.Fatalf("memory threshold = %#v, want migrated earliest mem override", memory)
	}
	if disk := byType["disk"]; disk.Duration != 2 || disk.Threshold != 85 {
		t.Fatalf("disk threshold = %#v, want default", disk)
	}
	if _, exists := byType["mem"]; exists {
		t.Fatal("legacy mem threshold should not remain active")
	}
}
