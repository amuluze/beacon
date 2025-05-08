// Package database
// Date       : 2024/8/22 14:24
// Author     : Amu
// Description:
package database

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(opts ...Option) (*DB, error) {
	opt := &option{
		Debug:        false,
		MaxIdleConns: 50,
		MaxOpenConns: 100,
		MaxLifetime:  7200,
	}
	for _, o := range opts {
		o(opt)
	}
	dial := dial(opt)

	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 开启调试模式
	if opt.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opt.MaxIdleConns)
	sqlDB.SetMaxOpenConns(opt.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(opt.MaxLifetime))

	return &DB{db}, nil
}

func dial(opt *option) gorm.Dialector {
	var dsn string
	var dialector gorm.Dialector
	switch opt.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.Port,
			opt.DBName,
		)
		dialector = mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,   // load size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		})
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
			opt.Host,
			opt.Port,
			opt.UserName,
			opt.DBName,
			opt.Password,
		)
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		})
	case "clickhouse":
		dsn = fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=200ms&max_execution_time=60",
			opt.UserName,
			opt.Password,
			opt.Host,
			opt.Port,
			opt.DBName,
		)
		dialector = clickhouse.Open(dsn)
	default:
		dsn = fmt.Sprintf("%s.db", opt.DBName)
		dialector = sqlite.Open(dsn)
	}

	return dialector
}

func (db *DB) Close() {
	if db != nil {
		conn, err := db.DB.DB()
		if err != nil {
			return
		}
		err = conn.Close()
		if err != nil {
			return
		}
	}
}

func (db *DB) RunInTransaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

func (db *DB) GetModel(model interface{}) *DB {
	return &DB{db.DB.Model(model)}
}
