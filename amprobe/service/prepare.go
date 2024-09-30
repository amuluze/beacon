// Package service
// Date: 2024/3/27 17:04
// Author: Amu
// Description:
package service

import (
	"log/slog"

	"amprobe/pkg/utils/hash"
	"amprobe/pkg/utils/uuid"
	"amprobe/service/model"
	"common/database"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var PrepareSet = wire.NewSet(wire.Struct(new(Prepare), "*"))

var user = model.User{
	ID:       uuid.MustUUID(),
	Username: "admin",
	Password: hash.SHA1String("admin123"), // hash.SHA1String(args.OldPassword)
	Remark:   "",
	IsAdmin:  1,
	Status:   1,
}

var roles = []model.Role{
	{
		ID:     uuid.MustUUID(),
		Name:   "管理员",
		Status: 1,
	},
	{
		ID:     uuid.MustUUID(),
		Name:   "普通用户",
		Status: 1,
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

	// init casbin rules
	a.InitCasbinRules()
}

func (a *Prepare) InitAccount(app *fiber.App) {
	var getResources []*model.Resource
	var postResources []*model.Resource

	for _, routers := range app.Stack() {
		for _, router := range routers {
			resource := &model.Resource{
				ID:     uuid.MustUUID(),
				Name:   router.Name,
				Path:   router.Path,
				Method: router.Method,
				Status: 1,
			}
			postResources = append(postResources, resource)
			if router.Method == "GET" || router.Name == "登录" || router.Name == "登出" || router.Name == "更新密码" || router.Name == "更新 token" {
				getResources = append(getResources, resource)
			}
		}
	}

	// 情况历史数据
	var u model.User
	if err := a.db.Model(&model.User{}).First(&u).Error; err == nil {
		slog.Info("admin exists", "username", u.Username)
		return
	}

	_ = a.db.RunInTransaction(func(tx *gorm.DB) error {
		var uRoles model.Roles
		for _, r := range roles {
			if r.Name == "管理员" {
				r.Resources = postResources
				uRoles = append(uRoles, &r)
			} else {
				r.Resources = getResources
			}
		}
		if err := a.db.Create(&roles).Error; err != nil {
			slog.Error("create roles error", "error", err)
			return err
		}

		user.Roles = uRoles
		if err := a.db.Create(&user).Error; err != nil {
			slog.Error("create user failed", "error", err)
			return err
		}
		return nil
	})
}

func (a *Prepare) InitCasbinRules() {
	var users []*model.User
	if err := a.db.Preload(clause.Associations).Preload("Roles").Find(&users).Error; err != nil {
		slog.Error("get all users error", "error", err)
		return
	}

	if err := a.db.Exec("delete from casbin_rule").Error; err != nil {
		slog.Error("delete from casbin_rule error", "error", err)
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
