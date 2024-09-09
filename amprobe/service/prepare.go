// Package service
// Date: 2024/3/27 17:04
// Author: Amu
// Description:
package service

import (
	"log/slog"
	"os"

	"amprobe/pkg/utils/hash"
	"amprobe/pkg/utils/uuid"
	"amprobe/service/model"

	"github.com/amuluze/amutool/database"
	"github.com/google/wire"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

var PrepareSet = wire.NewSet(wire.Struct(new(Prepare), "*"))

type Prepare struct {
	db *database.DB
}

func (a *Prepare) Init(config *Config) {
	if !config.InitData.Enable {
		return
	}
	file, err := os.ReadFile(config.InitData.InitConfigFile)
	if err != nil {
		slog.Error("read init config file failed", "error", err)
		return
	}
	var prepareData PrepareData
	err = yaml.Unmarshal(file, &prepareData)
	if err != nil {
		slog.Error("init config data unmarshal failed", "error", err)
		return
	}

	err = a.db.RunInTransaction(func(tx *gorm.DB) error {
		for _, u := range prepareData.Users {
			var ou model.User
			// 不存在则创建
			if err := a.db.Model(&model.User{}).Where("username = ?", u.Username).Take(&ou).Error; err != nil {
				a.db.Create(&model.User{
					ID:       uuid.MustUUID(),
					Username: u.Username,
					Password: hash.SHA1String(u.Password),
					Status:   u.Status,
					IsAdmin:  u.IsAdmin,
				})
			}
			// 存在则更新
			a.db.Model(&model.User{}).Where("username = ?", u.Username).Updates(model.User{Password: hash.SHA1String(u.Password), Status: u.Status, IsAdmin: u.IsAdmin})
		}
		return nil
	})
}

type PrepareData struct {
	Users Users `yaml:"users"`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Remark   string `yaml:"remark"`
	IsAdmin  string `yaml:"is_admin"`
	Status   int    `yaml:"status"`
}

type Users []*User
