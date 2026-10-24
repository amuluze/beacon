// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"strings"

	"beacon/service/model"
	"common/database"
)

func NewDB(config *Config, models *model.Models) (*database.DB, error) {
	if config.Gorm.GenDoc {
		return nil, nil
	}
	gormConfig := config.Gorm
	dbConfig := config.DB
	db, err := database.NewDB(
		database.WithDebug(gormConfig.Debug),
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
	if err := db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
		db.Close()
		return nil, err
	}
	// 时序表持续 append + 定期 delete，WAL 与主库文件会只增不减。
	// auto_vacuum=INCREMENTAL 只能在库首次创建时生效（在 AutoMigrate 之前设置），
	// 配合 cleanup 任务的 incremental_vacuum 增量回收已删除页面，避免 RSS 随 mmap 缓慢上涨。
	// synchronous=NORMAL 在 WAL 模式下安全且显著降低 fsync 开销。
	// 这些 PRAGMA 对历史库（已存在表）无副作用，对新库立即生效。
	if err := db.Exec("PRAGMA auto_vacuum=INCREMENTAL;").Error; err != nil {
		db.Close()
		return nil, err
	}
	if err := db.Exec("PRAGMA synchronous=NORMAL;").Error; err != nil {
		db.Close()
		return nil, err
	}
	if strings.EqualFold(gormConfig.DBType, "sqlite") || gormConfig.DBType == "" {
		if err := database.SecureSQLiteFiles(dbConfig.DBName); err != nil {
			db.Close()
			return nil, err
		}
	}

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
