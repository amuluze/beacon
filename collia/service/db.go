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
	// 启用WAL模式
	if err := db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
		db.Close()
		return nil, err
	}
	if strings.EqualFold(dbConfig.DBType, "sqlite") || dbConfig.DBType == "" {
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
