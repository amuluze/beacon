// Package service
// Date: 2024/06/10 18:42:59
// Author: Amu
// Description:
package service

import (
	"collia/service/model"
	"common/database"
	"strings"
)

func NewDB(config *Config, models *model.Models) (*database.DB, error) {
	dbConfig := config.DB
	db, err := database.NewDB(
		database.WithDebug(false),
		database.WithType(dbConfig.DBType),
		database.WithHost(dbConfig.Host),
		database.WithPort(dbConfig.Port),
		database.WithUsername(dbConfig.User),
		database.WithPassword(dbConfig.Password),
		database.WithDBName(dbConfig.DBName),
		database.WithMaxLifetime(7200),
		database.WithMaxOpenConns(150),
		database.WithMaxIdleConns(50),
	)
	if err != nil {
		return nil, err
	}
	// SQLite 专属 PRAGMA 与文件加固：非 SQLite 后端执行 WAL 语句会因
	// 语法错误而启动失败，必须限定在 sqlite/未指定类型分支内。
	if strings.EqualFold(dbConfig.DBType, "sqlite") || dbConfig.DBType == "" {
		if err := db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
			db.Close()
			return nil, err
		}
		if err := database.SecureSQLiteFiles(dbConfig.DBName); err != nil {
			db.Close()
			return nil, err
		}
	}

	err = db.AutoMigrate(models.GetAllModels()...)
	if err != nil {
		return nil, err
	}
	return db, nil
}
