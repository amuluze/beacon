// Package service
// Date: 2024/3/27 17:04
// Author: Amu
// Description:
package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"beacon/pkg/auth"
	"beacon/pkg/secret"
	"beacon/pkg/utils/hash"
	"beacon/pkg/utils/uuid"
	"beacon/service/model"
	"common/database"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var PrepareSet = wire.NewSet(wire.Struct(new(Prepare), "*"))

var notAdminResourcesMap = map[string]struct{}{
	"登录":              {},
	"登出":              {},
	"更新密码":            {},
	"更新 token":        {},
	"获取 WebSocket 票据": {},
}

var adminOnlyResourcesMap = map[string]struct{}{
	"查询钉钉告警配置":  {},
	"更新钉钉告警配置":  {},
	"测试钉钉告警":    {},
	"升级 Collia": {},
}

var thresholds = []model.AlarmThreshold{
	{
		Type:      "cpu",
		Duration:  2,
		Threshold: 80,
	},
	{
		Type:      "memory",
		Duration:  2,
		Threshold: 80,
	},
	{
		Type:      "disk",
		Duration:  2,
		Threshold: 85,
	},
}

const bootstrapAdminPasswordEnv = "BEACON_BOOTSTRAP_ADMIN_PASSWORD"

var (
	ErrBootstrapAdminPasswordRequired = errors.New("bootstrap administrator password is required")
	ErrBootstrapAdminPasswordWeak     = errors.New("bootstrap administrator password does not satisfy password policy")
)

// Historical public defaults are retained only as one-way hashes so existing
// installations can be migrated into a mandatory password-change flow. They
// are never used to create a new account.
var legacyPasswordMigrations = map[string]string{
	"f865b53623b121fd34ee5426c792e5c33af8c227": "$2a$10$iNfHGSdhRHO1wcG9/2gv7.0ZJxN3YciHxTIQUgoyNUOS0SoXL4vLe",
	"7c4a8d09ca3762af61e59520943dc26494f8941b": "$2a$10$yZhikKtaqyLpyrQOJE9WN.AAtWigY7XK145c3mFo5SU7/NYY/RyIK",
}

type Prepare struct {
	db       *database.DB
	cipher   *secret.Cipher
	enforcer *casbin.SyncedEnforcer
}

type NamePolicy struct {
	RoleID string
	Path   string
	Method string
}

type GroupPolicy struct {
	UserID string
	RoleID string
}

func (a *Prepare) Init(app *fiber.App) error {
	if err := migrateSensitiveData(context.Background(), a.db, a.cipher); err != nil {
		return err
	}
	// init account
	if err := a.InitAccount(app); err != nil {
		return err
	}

	// init alarm threshold
	if err := a.InitAlarmThreshold(); err != nil {
		return err
	}
	if err := a.InitAlarmMessageSetting(); err != nil {
		return err
	}

	// init casbin rules
	return a.InitCasbinRules()
}

func (a *Prepare) InitAccount(app *fiber.App) error {
	var notAdminResources []*model.Resource
	var adminResources []*model.Resource

	for _, routers := range app.Stack() {
		for _, router := range routers {
			if router.Path == "/" || (router.Method != "GET" && router.Method != "POST") {
				continue
			}

			resource := &model.Resource{
				ID:     uuid.MustUUID(),
				Name:   router.Name,
				Path:   router.Path,
				Method: router.Method,
				Status: 1,
			}
			adminResources = append(adminResources, resource)
			if _, ok := adminOnlyResourcesMap[router.Name]; ok {
				continue
			}
			if router.Method == "GET" {
				notAdminResources = append(notAdminResources, resource)
			}
			if _, ok := notAdminResourcesMap[router.Name]; ok {
				notAdminResources = append(notAdminResources, resource)
			}
		}
	}

	return a.db.RunInTransaction(func(tx *gorm.DB) error {
		// 更新 resource，并将数据库中的真实 ID 回填给角色关联。
		if err := syncRouteResources(tx, adminResources); err != nil {
			return err
		}

		adminRole, err := ensureDefaultRole(tx, "管理员")
		if err != nil {
			return err
		}
		normalRole, err := ensureDefaultRole(tx, "普通用户")
		if err != nil {
			return err
		}
		if err := syncDefaultRoleResources(tx, adminRole.Name, adminResources); err != nil {
			return err
		}
		if err := syncDefaultRoleResources(tx, normalRole.Name, notAdminResources); err != nil {
			return err
		}

		if err := migrateLegacyDefaultPasswords(tx); err != nil {
			return err
		}

		var adminUser model.User
		err = tx.Where("username = ?", "admin").First(&adminUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			password, err := bootstrapAdminPassword()
			if err != nil {
				return err
			}
			hashed, err := hash.BcryptHash(password)
			if err != nil {
				return fmt.Errorf("hash bootstrap administrator password: %w", err)
			}
			adminUser = model.User{
				ID:                 uuid.MustUUID(),
				Username:           "admin",
				Password:           hashed,
				MustChangePassword: true,
				Remark:             "管理员",
				IsAdmin:            1,
				Status:             1,
				Roles:              model.Roles{&adminRole},
			}
			if err := tx.Create(&adminUser).Error; err != nil {
				return fmt.Errorf("create bootstrap administrator: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("query administrator account: %w", err)
		} else if err := tx.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			return fmt.Errorf("ensure administrator role: %w", err)
		}
		return nil
	})
}

func bootstrapAdminPassword() (string, error) {
	password := os.Getenv(bootstrapAdminPasswordEnv)
	if password == "" {
		return "", fmt.Errorf("%w: set %s before first startup", ErrBootstrapAdminPasswordRequired, bootstrapAdminPasswordEnv)
	}
	if err := auth.ValidatePassword(password); err != nil {
		return "", fmt.Errorf("%w: %v", ErrBootstrapAdminPasswordWeak, err)
	}
	return password, nil
}

func ensureDefaultRole(tx *gorm.DB, name string) (model.Role, error) {
	var role model.Role
	err := tx.Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = model.Role{ID: uuid.MustUUID(), Name: name, Status: 1}
		if err := tx.Create(&role).Error; err != nil {
			return role, fmt.Errorf("create default role %q: %w", name, err)
		}
		return role, nil
	}
	if err != nil {
		return role, fmt.Errorf("query default role %q: %w", name, err)
	}
	return role, nil
}

func migrateLegacyDefaultPasswords(tx *gorm.DB) error {
	knownBcryptHashes := make([]string, 0, len(legacyPasswordMigrations))
	for legacyHash, bcryptHash := range legacyPasswordMigrations {
		if err := tx.Model(&model.User{}).Where("password = ?", legacyHash).Updates(map[string]any{
			"password":             bcryptHash,
			"must_change_password": true,
		}).Error; err != nil {
			return fmt.Errorf("migrate historical default password hash: %w", err)
		}
		knownBcryptHashes = append(knownBcryptHashes, bcryptHash)
	}
	if err := tx.Model(&model.User{}).Where("password IN ?", knownBcryptHashes).Update("must_change_password", true).Error; err != nil {
		return fmt.Errorf("mark historical default passwords for change: %w", err)
	}
	return nil
}

func syncRouteResources(tx *gorm.DB, resources []*model.Resource) error {
	for _, resource := range resources {
		var persisted model.Resource
		err := tx.Where("name = ?", resource.Name).First(&persisted).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(resource).Error; err != nil {
				return fmt.Errorf("create resource %q: %w", resource.Name, err)
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("query resource %q: %w", resource.Name, err)
		}

		if err := tx.Model(&persisted).Updates(map[string]any{
			"path":   resource.Path,
			"method": resource.Method,
			"status": resource.Status,
		}).Error; err != nil {
			return fmt.Errorf("update resource %q: %w", resource.Name, err)
		}
		persisted.Path = resource.Path
		persisted.Method = resource.Method
		persisted.Status = resource.Status
		*resource = persisted
	}
	return nil
}

// syncDefaultRoleResources keeps built-in roles aligned with the current HTTP
// route set. FirstOrCreate only writes associations for a newly-created role,
// so an upgrade that adds routes would otherwise leave existing admins with
// stale Casbin policies and return 403 for the new APIs.
func syncDefaultRoleResources(tx *gorm.DB, roleName string, resources []*model.Resource) error {
	var role model.Role
	if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("query default role %q: %w", roleName, err)
	}
	if err := tx.Model(&role).Association("Resources").Replace(resources); err != nil {
		return fmt.Errorf("sync default role %q resources: %w", roleName, err)
	}
	return nil
}

func (a *Prepare) InitAlarmThreshold() error {
	return a.db.RunInTransaction(syncAlarmThresholds)
}

func syncAlarmThresholds(tx *gorm.DB) error {
	for _, defaultThreshold := range thresholds {
		types := []string{defaultThreshold.Type}
		if defaultThreshold.Type == "memory" {
			types = append(types, "mem")
		}

		var records []model.AlarmThreshold
		if err := tx.Where("type IN ?", types).Order("id ASC").Find(&records).Error; err != nil {
			return fmt.Errorf("query %q alarm thresholds: %w", defaultThreshold.Type, err)
		}
		if len(records) == 0 {
			threshold := defaultThreshold
			if err := tx.Create(&threshold).Error; err != nil {
				return fmt.Errorf("create %q alarm threshold: %w", defaultThreshold.Type, err)
			}
			continue
		}

		keeper := records[0]
		if keeper.Type != defaultThreshold.Type {
			if err := tx.Model(&keeper).Update("type", defaultThreshold.Type).Error; err != nil {
				return fmt.Errorf("migrate %q alarm threshold: %w", keeper.Type, err)
			}
		}

		if len(records) > 1 {
			duplicateIDs := make([]uint, 0, len(records)-1)
			for _, record := range records[1:] {
				duplicateIDs = append(duplicateIDs, record.ID)
			}
			if err := tx.Where("id IN ?", duplicateIDs).Delete(&model.AlarmThreshold{}).Error; err != nil {
				return fmt.Errorf("delete duplicate %q alarm thresholds: %w", defaultThreshold.Type, err)
			}
		}
	}
	return nil
}

func (a *Prepare) InitAlarmMessageSetting() error {
	var setting model.AlarmMessageSetting
	if err := a.db.
		Where("key = ?", model.DefaultAlarmMessageSettingKey).
		Attrs(model.AlarmMessageSetting{MessageTitle: model.DefaultAlarmMessageTitle}).
		FirstOrCreate(&setting, model.AlarmMessageSetting{Key: model.DefaultAlarmMessageSettingKey}).Error; err != nil {
		return fmt.Errorf("initialize alarm message setting: %w", err)
	}
	return nil
}

func (a *Prepare) InitCasbinRules() error {
	var users []*model.User
	if err := a.db.Preload(clause.Associations).Preload("Roles").Find(&users).Error; err != nil {
		return fmt.Errorf("get all users: %w", err)
	}
	var roles []*model.Role
	if err := a.db.Preload("Resources").Find(&roles).Error; err != nil {
		return fmt.Errorf("get all role resources: %w", err)
	}

	desiredGroupPolicies := make(map[string][]string)
	desiredPolicies := make(map[string][]string)
	for _, user := range users {
		for _, role := range user.Roles {
			groupPolicy := []string{user.ID.String(), role.ID.String()}
			desiredGroupPolicies[casbinRuleKey(groupPolicy)] = groupPolicy
		}
	}
	for _, role := range roles {
		for _, resource := range role.Resources {
			policy := []string{role.ID.String(), resource.Path, resource.Method}
			desiredPolicies[casbinRuleKey(policy)] = policy
		}
	}

	if err := syncCasbinGroupingPolicies(a.enforcer, desiredGroupPolicies); err != nil {
		return err
	}
	return syncCasbinPolicies(a.enforcer, desiredPolicies)
}

func casbinRuleKey(rule []string) string {
	return strings.Join(rule, "\x00")
}

func syncCasbinGroupingPolicies(enforcer *casbin.SyncedEnforcer, desired map[string][]string) error {
	existing, err := enforcer.GetNamedGroupingPolicy("g")
	if err != nil {
		return fmt.Errorf("get grouping policies: %w", err)
	}
	for _, policy := range existing {
		if _, ok := desired[casbinRuleKey(policy)]; ok {
			continue
		}
		params := make([]interface{}, len(policy))
		for i, value := range policy {
			params[i] = value
		}
		if _, err := enforcer.RemoveNamedGroupingPolicy("g", params...); err != nil {
			return fmt.Errorf("remove stale grouping policy %v: %w", policy, err)
		}
	}
	for _, policy := range desired {
		if _, err := enforcer.AddNamedGroupingPolicy("g", policy[0], policy[1]); err != nil {
			return fmt.Errorf("add grouping policy %v: %w", policy, err)
		}
	}
	return nil
}

func syncCasbinPolicies(enforcer *casbin.SyncedEnforcer, desired map[string][]string) error {
	existing, err := enforcer.GetNamedPolicy("p")
	if err != nil {
		return fmt.Errorf("get policies: %w", err)
	}
	for _, policy := range existing {
		if _, ok := desired[casbinRuleKey(policy)]; ok {
			continue
		}
		params := make([]interface{}, len(policy))
		for i, value := range policy {
			params[i] = value
		}
		if _, err := enforcer.RemoveNamedPolicy("p", params...); err != nil {
			return fmt.Errorf("remove stale policy %v: %w", policy, err)
		}
	}
	for _, policy := range desired {
		if _, err := enforcer.AddNamedPolicy("p", policy[0], policy[1], policy[2]); err != nil {
			return fmt.Errorf("add policy %v: %w", policy, err)
		}
	}
	return nil
}
