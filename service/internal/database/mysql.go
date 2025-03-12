package database

import (
	"fmt"
	"glossika/service/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func newMySQL(dbName string, db config.MySQL) *gorm.DB {
	dataSourceName :=
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&readTimeout=30s&writeTimeout=30s",
			db.Account, db.Password, db.Host, db.Port, db.Database)

	gormDB, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	if db.MaxIdle == 0 || db.MaxOpen == 0 {
		log.Fatalf("%s missing maxIdle or maxOpen", dbName)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(db.MaxIdle)
	sqlDB.SetMaxOpenConns(db.MaxOpen)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Error pinging database %s: %v", dbName, err)
	}
	log.Printf("Pinged successfully mysql database: %s", dbName)

	return gormDB
}
