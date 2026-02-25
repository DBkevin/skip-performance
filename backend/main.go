package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/routes"
	"skin-performance/utils"
	"skin-performance/models"
)

func ptr(s string) *string {
	return &s
}

func main() {
	// 加载配置
	config.LoadConfig()

	// 设置 gin 模式
	gin.SetMode(config.AppConfig.GinMode)

	// 初始化数据库
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	log.Println("开始数据库迁移...")
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	// 初始化管理员用户（如果不存在）
	db = config.GetDB()
	var existingUser models.User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err != nil {
		// 创建管理员员工记录
		employee := models.Employee{
			Name:     "系统管理员",
			Role:     models.RoleAdmin,
			Phone:    ptr("13800138000"),
			JobNumber: ptr("ADMIN001"),
			IsActive: true,
		}
		if err := db.Create(&employee).Error; err == nil {
			// 加密密码
			hashedPassword, err := utils.HashPassword("admin123")
			if err == nil {
				// 创建管理员用户
				user := models.User{
					Username:   "admin",
					Password:   hashedPassword,
					EmployeeID: &employee.ID,
					Role:       "admin",
					IsActive:   true,
				}
				if err := db.Create(&user).Error; err == nil {
					log.Println("✅ 初始管理员用户创建成功")
					log.Println("用户名: admin")
					log.Println("密码: admin123")
				}
			}
		}
	}

	// 创建路由
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := strings.Split(config.AppConfig.CORSAllowOrigins, ",")
		
		// 检查是否允许该来源
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == origin || strings.TrimSpace(allowedOrigin) == "*" {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := config.AppConfig.ServerPort
	log.Printf("服务器启动在 :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
