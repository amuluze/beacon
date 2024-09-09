// Package service
// Date: 2024/3/6 11:08
// Author: Amu
// Description:
package service

import (
	"strings"

	"amprobe/service/model"

	"github.com/amuluze/amutool/database"
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
