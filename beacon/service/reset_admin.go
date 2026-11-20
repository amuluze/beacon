package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"beacon/pkg/auth"
	"beacon/pkg/utils/hash"
	"beacon/service/model"
	"common/database"

	"gorm.io/gorm"
)

// ErrAdminMissing 表示数据库中不存在 admin 账号。已初始化实例缺失管理员时
// 不允许通过本命令静默重建（IR024），必须从数据备份恢复。
var ErrAdminMissing = errors.New("administrator account (username=admin) not found; an initialized instance must not recreate it silently, restore from backup instead")

// ResetAdminFromConfig 按配置文件打开数据库并重置 admin 密码，供 `beacon reset-admin` 调用。
func ResetAdminFromConfig(configFile, password string) (string, bool, error) {
	config, err := NewConfig(configFile)
	if err != nil {
		return "", false, fmt.Errorf("load config %s: %w", configFile, err)
	}
	db, err := NewDB(config, &model.Models{})
	if err != nil {
		return "", false, fmt.Errorf("open database: %w", err)
	}
	if db == nil {
		return "", false, errors.New("database is disabled in this configuration (gorm.genDoc)")
	}
	defer db.Close()
	return ResetAdminPassword(db, password)
}

// ResetAdminPassword 重置 admin 密码：校验强度（显式密码）或生成临时密码、
// bcrypt 入库、强制下次登录改密，并写入一条系统审计。返回最终密码与是否由本命令生成。
func ResetAdminPassword(db *database.DB, password string) (string, bool, error) {
	generated := false
	if password == "" {
		var err error
		password, err = generateTemporaryPassword()
		if err != nil {
			return "", false, err
		}
		generated = true
	} else if err := auth.ValidatePassword(password); err != nil {
		return "", false, err
	}

	var admin model.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, ErrAdminMissing
		}
		return "", false, fmt.Errorf("query administrator: %w", err)
	}

	hashed, err := hash.BcryptHash(password)
	if err != nil {
		return "", false, fmt.Errorf("hash password: %w", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ?", admin.ID).Updates(map[string]interface{}{
			"password":             hashed,
			"must_change_password": true,
		}).Error; err != nil {
			return err
		}
		audit := model.Audit{
			Username: "system",
			Operate:  "reset admin password",
			Object:   admin.ID.String(),
		}
		return tx.Model(&model.Audit{}).Create(&audit).Error
	})
	if err != nil {
		return "", false, fmt.Errorf("update administrator password: %w", err)
	}
	return password, generated, nil
}

// generateTemporaryPassword 生成满足策略的临时密码（≥12 字符且含大小写与符号）。
func generateTemporaryPassword() (string, error) {
	raw := make([]byte, 18)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("read randomness: %w", err)
	}
	return "Beacon-" + base64.RawURLEncoding.EncodeToString(raw) + "!", nil
}
