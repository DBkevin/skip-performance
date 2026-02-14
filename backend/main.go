package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/routes"
)

func main() {
	// 命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	log.Println("正在加载配置...")
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("配置加载成功")

	// 初始化数据库
	log.Println("开始数据库迁移...")
	db, err := config.InitDBWithConfig(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	// 种子数据（初始化默认数据）
	SeedData()

	// 设置 gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

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
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	log.Printf("服务器启动在 :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
