package config

import (
	"log"
	"os"
	"strconv"

	"skin-performance/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// 加载配置
	LoadConfig()

	dsn := AppConfig.DBUser + ":" + AppConfig.DBPassword + "@tcp(" + 
		AppConfig.DBHost + ":" + AppConfig.DBPort + ")/" + 
		AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	log.Printf("连接数据库: %s", AppConfig.DBHost)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
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

// LoadConfig 从环境变量加载配置
func LoadConfig() {
	// 尝试加载.env文件（支持多种路径）
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	AppConfig = &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", "claw"),
		DBPassword:       getEnv("DB_PASSWORD", "thisopenclaw"),
		DBName:           getEnv("DB_NAME", "skin_performance"),
		JWTSecret:        getEnv("JWT_SECRET", "skin-performance-secret-key-2024"),
		JWTExpireHours:   getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		GinMode:          getEnv("GIN_MODE", gin.DebugMode),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
	}

	log.Println("配置加载完成")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	JWTExpireHours    int
	ServerPort        string
	GinMode           string
	CORSAllowOrigins  string
}

var AppConfig *Config
