package main

import (
	"log"
	"skin-performance/config"
	"skin-performance/models"
	"skin-performance/utils"
)

func initAdminUser() {
	db := config.GetDB()

	// 检查是否已存在admin用户
	var user models.User
	if err := db.Where("username = ?", "admin").First(&user).Error; err == nil {
		log.Println("管理员用户已存在")
		return
	}

	// 创建管理员员工记录
	employee := models.Employee{
		Name:     "系统管理员",
		Role:     models.RoleAdmin,
		Phone:    stringPtr("13800138000"),
		JobNumber: stringPtr("ADMIN001"),
		IsActive: true,
	}
	if err := db.Create(&employee).Error; err != nil {
		log.Fatalf("创建员工记录失败: %v", err)
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建管理员用户
	user = models.User{
		Username:   "admin",
		Password:   hashedPassword,
		EmployeeID: &employee.ID,
		Role:       "admin",
		IsActive:   true,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}

	log.Println("✅ 初始管理员用户创建成功")
	log.Println("用户名: admin")
	log.Println("密码: admin123")
}

func stringPtr(s string) *string {
	return &s
}