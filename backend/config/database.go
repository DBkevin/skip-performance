package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := "claw:thisopenclaw@tcp(120.25.70.117:3306)/skin_performance?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("database not initialized")
	}
	return DB
}
