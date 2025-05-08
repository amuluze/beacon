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

	_ = a.db.RunInTransaction(func(tx *gorm.DB) error {
		// 更新 resource
		for _, resource := range adminResources {
			if err := tx.Model(&model.Resource{}).FirstOrCreate(&resource).Error; err != nil {
				slog.Error("search or create resource failed", "error", err)
			}
		}

		// 更新 user role
		for _, u := range users {
			if u.Username == "admin" {
				u.Roles[0].Resources = adminResources
			} else {
				u.Roles[0].Resources = notAdminResources
			}
			// 创建或更新
			if err := tx.Model(&model.User{}).FirstOrCreate(&u).Error; err != nil {
				slog.Error("search or create user failed", "error", err)
			}
		}
		return nil
	})
}

func (a *Prepare) InitAlarmThreshold() {
	for _, threshold := range thresholds {
		if err := a.db.Model(&model.AlarmThreshold{}).FirstOrCreate(&threshold).Error; err != nil {
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

	for _, user := range users {
		for _, role := range user.Roles {
			if _, err := a.enforcer.AddNamedGroupingPolicy("g", user.ID.String(), role.ID.String()); err != nil {
				slog.Error("add grouping policy error", "error", err)
			}
			var roleResources model.Role
			if err := a.db.Where("id = ?", role.ID).Preload("Resources").Find(&roleResources).Error; err != nil {
				slog.Error("get role resources error", "error", err)
			}
			for _, resource := range roleResources.Resources {
				if _, err := a.enforcer.AddNamedPolicy("p", role.ID.String(), resource.Path, resource.Method); err != nil {
					slog.Error("add policy error", "error", err)
				}
			}
		}
	}
}
