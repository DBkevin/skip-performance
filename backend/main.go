package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/routes"
)

func main() {
	// 初始化数据库
	if _, err := config.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 设置 gin 模式
	gin.SetMode(gin.DebugMode)

	// 创建路由
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
