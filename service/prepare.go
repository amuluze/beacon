// Package service
// Date: 2024/3/27 17:04
// Author: Amu
// Description:
package service

import (
	"amprobe/pkg/hash"
	"amprobe/pkg/uuid"
	"log/slog"

	"amprobe/pkg/database"
	"amprobe/service/model"

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

var users = []*model.User{
	{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: hash.SHA1String("admin123"), // hash.SHA1String(args.OldPassword)
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
		Username: "amprobe",
		Password: hash.SHA1String("123456"),
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

	// 收集资源
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
			if router.Method == "GET" {
				notAdminResources = append(notAdminResources, resource)
			}
			if _, ok := notAdminResourcesMap[router.Name]; ok {
				notAdminResources = append(notAdminResources, resource)
			}
		}
	}

	if err := a.db.RunInTransaction(func(tx *gorm.DB) error {
		// 更新 resource
		for _, resource := range adminResources {
			// First try to find existing resource by name
			var existingResource model.Resource
			if err := tx.Model(&model.Resource{}).Where("name = ?", resource.Name).First(&existingResource).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create new resource if not exists
					if createErr := tx.Create(&resource).Error; createErr != nil {
						slog.Error("创建资源失败", "resource", resource.Name, "error", err)
						return err
					}
				} else {
					slog.Error("查询资源失败", "resource", resource.Name, "error", err)
					return err
				}
			} else {
				// Update existing resource
				resource.ID = existingResource.ID
				if err := tx.Model(&existingResource).Updates(resource).Error; err != nil {
					slog.Error("更新资源失败", "resource", resource.Name, "error", err)
					return err
				}
			}
		}

		// 更新 user role
		for _, u := range users {
			var existingUser model.User
			if err := tx.Model(&model.User{}).Where("username = ?", u.Username).First(&existingUser).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create new user if not exists
					if u.Username == "admin" {
						u.Roles[0].Resources = adminResources
					} else {
						u.Roles[0].Resources = notAdminResources
					}
					if createErr := tx.Create(&u).Error; createErr != nil {
						slog.Error("创建用户失败", "username", u.Username, "error", err)
						return err
					}
				} else {
					slog.Error("查询用户失败", "username", u.Username, "error", err)
					return err
				}
			} else {
				// Update existing user
				if u.Username == "admin" {
					u.Roles[0].Resources = adminResources
				} else {
					u.Roles[0].Resources = notAdminResources
				}
				u.ID = existingUser.ID
				if err := tx.Model(&existingUser).Updates(u).Error; err != nil {
					slog.Error("更新用户失败", "username", u.Username, "error", err)
					return err
				}
			}
		}
		return nil
	}); err != nil {
		slog.Error("初始化账户事务失败", "error", err)
	}
}

func (a *Prepare) InitAlarmThreshold() {
	if err := a.db.RunInTransaction(func(tx *gorm.DB) error {
		for _, threshold := range thresholds {
			if err := tx.Model(&model.AlarmThreshold{}).FirstOrCreate(&threshold).Error; err != nil {
				slog.Error("创建或更新告警阈值失败", "type", threshold.Type, "error", err)
				return err
			}
		}
		return nil
	}); err != nil {
		slog.Error("初始化告警阈值事务失败", "error", err)
	}
}

func (a *Prepare) InitCasbinRules() {
	var users []*model.User
	if err := a.db.Preload(clause.Associations).Preload("Roles").Find(&users).Error; err != nil {
		slog.Error("获取所有用户失败", "error", err)
		return
	}

	for _, user := range users {
		for _, role := range user.Roles {
			if _, err := a.enforcer.AddNamedGroupingPolicy("g", user.ID.String(), role.ID.String()); err != nil {
				slog.Error("添加用户角色关系失败", "userId", user.ID, "roleId", role.ID, "error", err)
				continue
			}

			var roleResources model.Role
			if err := a.db.Where("id = ?", role.ID).Preload("Resources").Find(&roleResources).Error; err != nil {
				slog.Error("获取角色资源失败", "roleId", role.ID, "error", err)
				continue
			}

			for _, resource := range roleResources.Resources {
				if _, err := a.enforcer.AddNamedPolicy("p", role.ID.String(), resource.Path, resource.Method); err != nil {
					slog.Error("添加权限策略失败", "roleId", role.ID, "resource", resource.Path, "error", err)
				}
			}
		}
	}
}
