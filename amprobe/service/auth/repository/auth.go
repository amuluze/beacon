// Package repository
// Date: 2024/3/27 16:39
// Author: Amu
// Description:
package repository

import (
	"context"

	"amprobe/pkg/errors"
	"amprobe/pkg/utils/hash"
	"amprobe/service/model"
	"amprobe/service/schema"
	"common/database"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var AuthRepoSet = wire.NewSet(NewAuthRepo, wire.Bind(new(IAuthRepository), new(*AuthRepo)))

var _ IAuthRepository = (*AuthRepo)(nil)

type IAuthRepository interface {
	Login(ctx context.Context, args schema.LoginArgs) (model.User, error)
	PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error
	UserInfo(ctx context.Context, userID string) (model.User, error)
}

type AuthRepo struct {
	DB *database.DB
}

func NewAuthRepo(db *database.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

func (a *AuthRepo) Login(ctx context.Context, args schema.LoginArgs) (model.User, error) {
	var user model.User
	err := a.DB.RunInTransaction(func(tx *gorm.DB) error {
		if err := tx.Where("username = ?", args.Username).Where("status = ?", 1).First(&user).Error; err != nil {
			return err
		}
		if user.Password != hash.SHA1String(args.Password) {
			return errors.New("invalid password")
		}

		return nil
	})
	return user, err
}

func (a *AuthRepo) PassUpdate(ctx context.Context, args schema.PasswordUpdateArgs) error {
	return a.DB.RunInTransaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Where("username = ?", args.Username).First(&user).Error; err != nil {
			return err
		}

		if user.Password != hash.SHA1String(args.OldPassword) {
			return errors.New("invalid password")
		}

		if err := tx.Model(&model.User{}).Where("username = ?", args.Username).Update("password", hash.SHA1String(args.NewPassword)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (a *AuthRepo) UserInfo(ctx context.Context, userID string) (model.User, error) {
	var user model.User
	if err := a.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
