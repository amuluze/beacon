// Package repository
// Date       : 2024/9/4 14:56
// Author     : Amu
// Description:
package repository

import (
	"amprobe/pkg/database"
	"amprobe/pkg/hash"
	"amprobe/pkg/uuid"
	"amprobe/service/model"
	"amprobe/service/schema"
	"context"
	"log/slog"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var AccountRepositorySet = wire.NewSet(NewAccountRepository, wire.Bind(new(IAccountRepository), new(*AccountRepository)))

var _ IAccountRepository = (*AccountRepository)(nil)

type IAccountRepository interface {
	UserQuery(context.Context, schema.UserQueryArgs) (model.Users, error)
	UserCreate(context.Context, schema.UserCreateArgs) (model.User, error)
	UserUpdate(context.Context, schema.UserUpdateArgs) (model.User, error)
	UserDelete(context.Context, schema.UserDeleteArgs) error
	UserCount(context.Context) (int64, error)

	RoleQuery(context.Context, schema.RoleQueryArgs) (model.Roles, error)
	RoleCreate(context.Context, schema.RoleCreateArgs) (model.Role, error)
	RoleUpdate(context.Context, schema.RoleUpdateArgs) (model.Role, error)
	RoleDelete(context.Context, schema.RoleDeleteArgs) error
	RoleCount(context.Context) (int64, error)

	ResourceQuery(context.Context, schema.ResourceQueryArgs) (model.Resources, error)
	ResourceCount(context.Context) (int64, error)
}

type AccountRepository struct {
	DB *database.DB
}

func NewAccountRepository(db *database.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (a *AccountRepository) UserQuery(ctx context.Context, args schema.UserQueryArgs) (model.Users, error) {
	db := a.DB.GetModel(&model.User{})
	db = database.OptionDB(
		db,
		database.WithId(args.ID),
		database.WithIds(args.IDs),
		database.WithUserName(args.Username),
		database.WithUsernames(args.Usernames),
		database.WithStatus(args.Status),
		database.OrderBy("created_at DESC"),
	)

	var users model.Users
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func (a *AccountRepository) UserCount(ctx context.Context) (int64, error) {
	var count int64
	if err := a.DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (a *AccountRepository) UserCreate(ctx context.Context, args schema.UserCreateArgs) (model.User, error) {
	u := model.User{
		ID:       uuid.MustUUID(),
		Username: args.Username,
		Password: hash.SHA1String(args.Password),
		Status:   args.Status,
		Remark:   args.Remark,
	}
	if len(args.RoleIDs) > 0 {
		var roles model.Roles
		if err := a.DB.Where("id IN (?)", args.RoleIDs).Find(&roles).Error; err != nil {
			return model.User{}, err
		}
		for _, r := range roles {
			if r.Name == "admin" {
				u.IsAdmin = 1
			}
		}
		u.Roles = roles
	}
	if err := a.DB.Create(&u).Error; err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (a *AccountRepository) UserUpdate(ctx context.Context, args schema.UserUpdateArgs) (model.User, error) {
	params := make(map[string]interface{})
	if args.Status != 0 {
		params["status"] = args.Status
	}
	if args.Remark != "" {
		params["remark"] = args.Remark
	}
	if args.Username != "" {
		params["username"] = args.Username
	}
	var user model.User
	err := a.DB.RunInTransaction(func(tx *gorm.DB) error {
		var roles model.Roles
		if len(args.RoleIDs) > 0 {
			if err := tx.Where("id IN (?)", args.RoleIDs).Find(&roles).Error; err != nil {
				return err
			}
			for _, r := range roles {
				if r.Name == "admin" {
					params["is_admin"] = 1
				}
			}
		}
		if err := tx.Model(&model.User{}).Where("id = ?", args.ID).Updates(params).Error; err != nil {
			return err
		}
		// 替换关联
		if len(roles) == 0 {
			if err := a.DB.Model(&model.User{}).Where("id = ?", args.ID).First(&user).Error; err != nil {
				return err
			}
			if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
				slog.Error("Error replace roles", "error", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	if err = a.DB.Preload("Roles").Where("id = ?", args.ID).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (a *AccountRepository) UserDelete(ctx context.Context, args schema.UserDeleteArgs) error {
	db := a.DB.GetModel(&model.User{})
	db = database.OptionDB(
		db,
		database.WithIds(args.IDs),
	)

	if err := db.Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *AccountRepository) RoleQuery(ctx context.Context, args schema.RoleQueryArgs) (model.Roles, error) {
	db := a.DB.GetModel(&model.Role{})
	db = database.OptionDB(
		db,
		database.WithId(args.ID),
		database.WithIds(args.IDs),
		database.WithName(args.Name),
		database.WithNames(args.Names),
		database.WithStatus(args.Status),
		database.OrderBy("created_at DESC"),
	)
	var roles model.Roles
	if err := db.Preload("Resources").Find(&roles).Error; err != nil {
		return roles, err
	}
	return roles, nil
}

func (a *AccountRepository) RoleCount(ctx context.Context) (int64, error) {
	var count int64
	if err := a.DB.Model(&model.Role{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (a *AccountRepository) RoleCreate(ctx context.Context, args schema.RoleCreateArgs) (model.Role, error) {
	r := model.Role{
		ID:     uuid.MustUUID(),
		Name:   args.Name,
		Status: args.Status,
		Remark: args.Remark,
	}
	if len(args.ResourceIDs) > 0 {
		var resources model.Resources
		if err := a.DB.Where("id IN (?)", args.ResourceIDs).Find(&resources).Error; err != nil {
			return model.Role{}, err
		}
		r.Resources = resources
	}
	if err := a.DB.Create(&r).Error; err != nil {
		return model.Role{}, err
	}
	return r, nil
}

func (a *AccountRepository) RoleUpdate(ctx context.Context, args schema.RoleUpdateArgs) (model.Role, error) {
	params := make(map[string]interface{})
	if args.Status != 0 {
		params["status"] = args.Status
	}
	if args.Remark != "" {
		params["remark"] = args.Remark
	}
	if args.Name != "" {
		params["name"] = args.Name
	}
	err := a.DB.RunInTransaction(func(tx *gorm.DB) error {
		if len(args.ResourceIDs) > 0 {
			var resource model.Resources
			if err := tx.Where("id IN (?)", args.ResourceIDs).Find(&resource).Error; err != nil {
				return err
			}
			err := tx.Model(&model.Role{}).Association("Resources").Replace(resource)
			if err != nil {
				return err
			}
		}
		if err := tx.Model(&model.Role{}).Where("id = ?", args.ID).Updates(params).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Role{}, err
	}
	var role model.Role
	if err := a.DB.Preload("Resources").Where("id = ?", args.ID).First(&role).Error; err != nil {
		return model.Role{}, err
	}
	return role, nil
}

func (a *AccountRepository) RoleDelete(ctx context.Context, args schema.RoleDeleteArgs) error {
	db := a.DB.GetModel(&model.Role{})
	db = database.OptionDB(
		db,
		database.WithIds(args.IDs),
	)
	if err := db.Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *AccountRepository) ResourceQuery(ctx context.Context, args schema.ResourceQueryArgs) (model.Resources, error) {
	db := a.DB.GetModel(&model.Resource{})
	db = database.OptionDB(
		db,
		database.WithId(args.ID),
	)
	var resources model.Resources
	if err := db.Order("created_at desc").Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func (a *AccountRepository) ResourceCount(ctx context.Context) (int64, error) {
	var count int64
	if err := a.DB.Model(&model.Resource{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
