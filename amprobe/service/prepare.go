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

var users = []*model.User{
	{
		ID:       uuid.MustUUID(),
		Username: "admin",
		Password: hash.SHA1String("admin123"), // hash.SHA1String(args.OldPassword)
		Remark:   "",
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
		Remark:   "",
		IsAdmin:  1,
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

	_ = a.db.RunInTransaction(func(tx *gorm.DB) error {
		for _, u := range users {
			if u.Username == "admin" {
				u.Roles[0].Resources = postResources
			} else {
				u.Roles[0].Resources = getResources
			}

			var ou model.User
			// 不存在则创建
			if err := a.db.Model(&model.User{}).Where("username = ?", u.Username).Take(&ou).Error; err != nil {
				a.db.Create(u)
			} else {
				// 存在则更新
				a.db.Model(&model.User{}).Where("username = ?", u.Username).Updates(model.User{Password: u.Password, Status: u.Status, IsAdmin: u.IsAdmin, Roles: u.Roles, Remark: u.Remark})
			}
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
