// Package service
// Date: 2024/06/10 18:42:59
// Author: Amu
// Description:
package service

import (
	"strings"

	"github.com/amuluze/amutool/database"
	"github.com/amuluze/amvector/service/model"
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
