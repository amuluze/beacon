// Package service
// Date: 2024/3/27 17:04
// Author: Amu
// Description:
package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
	"登录":       {},
	"登出":       {},
	"更新密码":     {},
	"更新 token": {},
}

var adminOnlyResourcesMap = map[string]struct{}{
	"查询钉钉告警配置": {},
	"更新钉钉告警配置": {},
	"测试钉钉告警":   {},
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

const (
	defaultAdminPasswordHash  = "$2a$10$iNfHGSdhRHO1wcG9/2gv7.0ZJxN3YciHxTIQUgoyNUOS0SoXL4vLe"
	defaultBeaconPasswordHash = "$2a$10$yZhikKtaqyLpyrQOJE9WN.AAtWigY7XK145c3mFo5SU7/NYY/RyIK"
)

var legacyDefaultPasswordHashes = map[string]string{
	"admin":  "f865b53623b121fd34ee5426c792e5c33af8c227",
	"beacon": "7c4a8d09ca3762af61e59520943dc26494f8941b",
}

var users = []*model.User{
	{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: defaultAdminPasswordHash,
		Remark:   "管理员",
		IsAdmin:  1,
		Status:   1,
		Roles: []*model.Role{
			{
				ID:     uuid.MustUUID(),
				Name:   "管理员",
				Status: 1,
			},
		},
	},
	{
		ID:       uuid.MustUUID(),
		Username: "beacon",
		Password: defaultBeaconPasswordHash,
		Remark:   "普通用户",
		IsAdmin:  2,
		Status:   1,
		Roles: []*model.Role{
			{
				ID:     uuid.MustUUID(),
				Name:   "普通用户",
				Status: 1,
			},
		},
	},
}

type Prepare struct {
	db       *database.DB
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

func (a *Prepare) Init(app *fiber.App) {
	// init account
	a.InitAccount(app)

	// init alarm threshold
	a.InitAlarmThreshold()

	// init casbin rules
	a.InitCasbinRules()
}

func (a *Prepare) InitAccount(app *fiber.App) {
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

	if err := a.db.RunInTransaction(func(tx *gorm.DB) error {
		// 更新 resource，并将数据库中的真实 ID 回填给角色关联。
		if err := syncRouteResources(tx, adminResources); err != nil {
			return err
		}

		// 更新 user role
		for _, u := range users {
			resources := notAdminResources
			if u.Username == "admin" {
				resources = adminResources
			}
			if legacyHash, ok := legacyDefaultPasswordHashes[u.Username]; ok {
				if err := tx.Model(&model.User{}).Where("username = ? AND password = ?", u.Username, legacyHash).Update("password", u.Password).Error; err != nil {
					return fmt.Errorf("upgrade default user %q password hash: %w", u.Username, err)
				}
			}

			// 使用副本创建默认账号，避免查询已有账号时修改全局模板。
			roleTemplate := *u.Roles[0]
			roleTemplate.Resources = resources
			userTemplate := *u
			userTemplate.Roles = []*model.Role{&roleTemplate}
			var existingUser model.User
			err := tx.Where("username = ?", u.Username).First(&existingUser).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = tx.Create(&userTemplate).Error
			}
			if err != nil {
				return fmt.Errorf("search or create user %q: %w", u.Username, err)
			}
			if err := syncDefaultRoleResources(tx, roleTemplate.Name, resources); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		slog.Error("initialize accounts and role resources failed", "error", err)
	}
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

func (a *Prepare) InitAlarmThreshold() {
	for _, threshold := range thresholds {
		var current model.AlarmThreshold
		if err := a.db.Where("type = ?", threshold.Type).Attrs(threshold).FirstOrCreate(&current).Error; err != nil {
			slog.Error("alarm threshold exist", "error", err)
		}
	}
}

func (a *Prepare) InitCasbinRules() {
	var users []*model.User
	if err := a.db.Preload(clause.Associations).Preload("Roles").Find(&users).Error; err != nil {
		slog.Error("get all users error", "error", err)
		return
	}
	var roles []*model.Role
	if err := a.db.Preload("Resources").Find(&roles).Error; err != nil {
		slog.Error("get all role resources error", "error", err)
		return
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

	syncCasbinGroupingPolicies(a.enforcer, desiredGroupPolicies)
	syncCasbinPolicies(a.enforcer, desiredPolicies)
}

func casbinRuleKey(rule []string) string {
	return strings.Join(rule, "\x00")
}

func syncCasbinGroupingPolicies(enforcer *casbin.SyncedEnforcer, desired map[string][]string) {
	existing, err := enforcer.GetNamedGroupingPolicy("g")
	if err != nil {
		slog.Error("get grouping policies error", "error", err)
		return
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
			slog.Error("remove stale grouping policy error", "policy", policy, "error", err)
		}
	}
	for _, policy := range desired {
		if _, err := enforcer.AddNamedGroupingPolicy("g", policy[0], policy[1]); err != nil {
			slog.Error("add grouping policy error", "policy", policy, "error", err)
		}
	}
}

func syncCasbinPolicies(enforcer *casbin.SyncedEnforcer, desired map[string][]string) {
	existing, err := enforcer.GetNamedPolicy("p")
	if err != nil {
		slog.Error("get policies error", "error", err)
		return
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
			slog.Error("remove stale policy error", "policy", policy, "error", err)
		}
	}
	for _, policy := range desired {
		if _, err := enforcer.AddNamedPolicy("p", policy[0], policy[1], policy[2]); err != nil {
			slog.Error("add policy error", "policy", policy, "error", err)
		}
	}
}
