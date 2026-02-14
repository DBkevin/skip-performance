package config

import (
	"log"

	"skin-performance/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDBWithConfig 使用配置初始化数据库
func InitDBWithConfig(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()

	logLevel := logger.Silent
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// InitDB 保持向后兼容
func InitDB() (*gorm.DB, error) {
	if AppConfig == nil {
		// 尝试加载默认配置
		cfg, err := LoadConfig("")
		if err != nil {
			return nil, err
		}
		return InitDBWithConfig(cfg)
	}
	return InitDBWithConfig(AppConfig)
}

func AutoMigrate(db *gorm.DB) error {
	// 禁用外键约束检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	return db.AutoMigrate(
		&models.Customer{},
		&models.Project{},
		&models.Employee{},
		&models.User{},
		&models.Visit{},
		&models.VisitItem{},
		&models.RevisitRecord{},
		&models.ProductConsumption{},
	)
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("database not initialized")
	}
	return DB
}
