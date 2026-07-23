package service

import (
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
	prepare.InitAlarmThreshold()

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
	prepare.InitAlarmThreshold()

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

func TestDefaultUserPasswordsUseLoginHash(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{name: "admin", username: "admin", password: "admin123"},
		{name: "beacon", username: "beacon", password: "123456"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, user := range users {
				if user.Username != tc.username {
					continue
				}
				if err := hash.BcryptVerify(tc.password, user.Password); err != nil {
					t.Fatalf("default password for %q should pass bcrypt login verification: %v", tc.username, err)
				}
				return
			}
			t.Fatalf("default user %q not found", tc.username)
		})
	}
}

func TestLegacyDefaultPasswordHashesMatchPreviousDefaults(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{name: "admin", username: "admin", password: "admin123"},
		{name: "beacon", username: "beacon", password: "123456"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := legacyDefaultPasswordHashes[tc.username]
			if !ok {
				t.Fatalf("legacy hash for %q not found", tc.username)
			}
			if want := hash.SHA1String(tc.password); got != want {
				t.Fatalf("legacy hash for %q = %q, want %q", tc.username, got, want)
			}
		})
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
	admin := &model.User{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: defaultAdminPasswordHash,
		IsAdmin:  1,
		Status:   1,
		Roles:    model.Roles{adminRole},
	}
	normal := &model.User{
		ID:       uuid.MustUUID(),
		Username: "beacon",
		Password: defaultBeaconPasswordHash,
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
	prepare.InitAccount(app)
	prepare.InitAccount(app)
	prepare.InitCasbinRules()

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
