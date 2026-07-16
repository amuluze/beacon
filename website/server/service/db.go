// Package service
// Date: 2025/02/12 15:04:23
// Author: Amu
// Description:
package service

import (
	"server/pkg/database"
	"server/service/model"
	"strings"
)

func NewDB(config *Config, models *model.Models) (*database.DB, error) {
	if config.Gorm.GenDoc {
		return nil, nil
	}
	gormConfig := config.Gorm
	dbConfig := config.DB
	// 生产环境强制关闭 SQL 调试输出，避免日志泄露与体积膨胀
	debug := gormConfig.Debug && !config.App.IsProduction()
	db, err := database.NewDB(
		database.WithDebug(debug),
		database.WithType(gormConfig.DBType),
		database.WithHost(dbConfig.Host),
		database.WithPort(dbConfig.Port),
		database.WithUsername(dbConfig.User),
		database.WithPassword(dbConfig.Password),
		database.WithDBName(dbConfig.DBName),
		database.WithMaxLifetime(gormConfig.MaxLifetime),
		database.WithMaxOpenConns(gormConfig.MaxOpenConns),
		database.WithMaxIdleConns(gormConfig.MaxIdleConns),
	)
	if err != nil {
		return nil, err
	}
	// SQLite 启用WAL模式
	db.Exec("PRAGMA journal_mode=WAL;")

	if gormConfig.EnableAutoMigrate {
		if dbType := gormConfig.DBType; strings.ToLower(dbType) == "mysql" {
			db.Set("gorm:table_options", "ENGINE=InnoDB")
		}
		err := db.AutoMigrate(models.GetAllModels()...)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
