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
		if err := tx.Where("username = ?", args.Username).First(&user).Error; err != nil {
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

		if err := tx.Where("username = ?", args.Username).Update("password", hash.SHA1String(args.OldPassword)).Error; err != nil {
			return err
		}
		return nil
	})
}
