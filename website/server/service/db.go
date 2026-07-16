// Package service
// Date: 2025/02/12 15:04:23
// Author: Amu
// Description:
package service

import (
	"log/slog"
	"server/pkg/database"
	"server/service/model"
	"strings"
)

// NewDB 构造数据库连接并返回清理函数，供依赖注入的 cleanup 链在退出时关闭连接。
func NewDB(config *Config, models *model.Models) (*database.DB, func(), error) {
	if config.Gorm.GenDoc {
		return nil, func() {}, nil
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
		database.WithSSLMode(dbConfig.SSLMode),
		database.WithMaxLifetime(gormConfig.MaxLifetime),
		database.WithMaxOpenConns(gormConfig.MaxOpenConns),
		database.WithMaxIdleConns(gormConfig.MaxIdleConns),
	)
	if err != nil {
		return nil, nil, err
	}
	// PRAGMA 仅 SQLite 支持；DBType 为空时默认走 sqlite 驱动
	if gormConfig.DBType == "" || strings.EqualFold(gormConfig.DBType, "sqlite") {
		db.Exec("PRAGMA journal_mode=WAL;")
	}

	if gormConfig.EnableAutoMigrate {
		if dbType := gormConfig.DBType; strings.ToLower(dbType) == "mysql" {
			db.Set("gorm:table_options", "ENGINE=InnoDB")
		}
		err := db.AutoMigrate(models.GetAllModels()...)
		if err != nil {
			return nil, nil, err
		}
	}
	cleanup := func() {
		slog.Info("closing database connection")
		if err := db.Close(); err != nil {
			slog.Warn("close db error", "err", err)
		}
	}
	return db, cleanup, nil
}
