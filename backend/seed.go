package main

import (
	"log"
	"skin-performance/config"
	"skin-performance/models"
	"skin-performance/utils"
	"time"
)

// SeedData 初始化数据库数据
func SeedData() {
	db := config.GetDB()

	// 检查是否已有数据
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("数据已存在，跳过初始化")
		return
	}

	log.Println("开始初始化数据...")

	// 1. 创建管理员员工
	adminEmployee := models.Employee{
		Name:       "管理员",
		Role:       models.RoleAdmin,
		Department: func() *string { s := "管理"; return &s }(),
		JobNumber:  func() *string { s := "A001"; return &s }(),
		Phone:      func() *string { s := "13800138000"; return &s }(),
		IsActive:   true,
		CreatedAt:  func() *time.Time { t := time.Now(); return &t }(),
		UpdatedAt:  func() *time.Time { t := time.Now(); return &t }(),
	}
	db.Create(&adminEmployee)

	// 2. 创建医生
	doctors := []models.Employee{
		{Name: "张医生", Role: models.RoleDoctor, Department: func() *string { s := "皮肤科"; return &s }(), JobNumber: func() *string { s := "D001"; return &s }(), Phone: func() *string { s := "13800138001"; return &s }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "李医生", Role: models.RoleDoctor, Department: func() *string { s := "美容科"; return &s }(), JobNumber: func() *string { s := "D002"; return &s }(), Phone: func() *string { s := "13800138002"; return &s }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
	}
	db.Create(&doctors)

	// 3. 创建护士
	nurses := []models.Employee{
		{Name: "王护士", Role: models.RoleNurse, Department: func() *string { s := "皮肤科"; return &s }(), JobNumber: func() *string { s := "N001"; return &s }(), Phone: func() *string { s := "13800138003"; return &s }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "刘护士", Role: models.RoleNurse, Department: func() *string { s := "美容科"; return &s }(), JobNumber: func() *string { s := "N002"; return &s }(), Phone: func() *string { s := "13800138004"; return &s }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
	}
	db.Create(&nurses)

	// 4. 创建咨询师
	consultant := models.Employee{
		Name:       "陈咨询师",
		Role:       models.RoleConsultant,
		Department: func() *string { s := "前台"; return &s }(),
		JobNumber:  func() *string { s := "C001"; return &s }(),
		Phone:      func() *string { s := "13800138005"; return &s }(),
		IsActive:   true,
		CreatedAt:  func() *time.Time { t := time.Now(); return &t }(),
		UpdatedAt:  func() *time.Time { t := time.Now(); return &t }(),
	}
	db.Create(&consultant)

	// 5. 创建用户账号
	hashedPassword, _ := utils.HashPassword("admin123")
	adminUser := models.User{
		Username:   "admin",
		Password:   hashedPassword,
		EmployeeID: &adminEmployee.ID,
		Role:       models.RoleAdmin,
		IsActive:   true,
		CreatedAt:  func() *time.Time { t := time.Now(); return &t }(),
		UpdatedAt:  func() *time.Time { t := time.Now(); return &t }(),
	}
	db.Create(&adminUser)

	// 6. 创建项目
	projects := []models.Project{
		{Name: "肉毒素注射", Category: func() *string { s := "注射类"; return &s }(), StandardPrice: func() *float64 { p := 2800.0; return &p }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "光子嫩肤", Category: func() *string { s := "激光类"; return &s }(), StandardPrice: func() *float64 { p := 1200.0; return &p }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "玻尿酸填充", Category: func() *string { s := "注射类"; return &s }(), StandardPrice: func() *float64 { p := 3500.0; return &p }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "水光针", Category: func() *string { s := "注射类"; return &s }(), StandardPrice: func() *float64 { p := 800.0; return &p }(), IsActive: true, CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
	}
	db.Create(&projects)

	// 7. 创建顾客
	customers := []models.Customer{
		{Name: "王小姐", Phone: "13900139001", CustomerType: func() *string { s := "初诊"; return &s }(), FirstVisitDate: func() *time.Time { t := time.Now().AddDate(0, 0, -30); return &t }(), CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "李先生", Phone: "13900139002", CustomerType: func() *string { s := "复诊"; return &s }(), FirstVisitDate: func() *time.Time { t := time.Now().AddDate(0, 0, -60); return &t }(), CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
		{Name: "张女士", Phone: "13900139003", CustomerType: func() *string { s := "再消费"; return &s }(), FirstVisitDate: func() *time.Time { t := time.Now().AddDate(0, 0, -90); return &t }(), CreatedAt: func() *time.Time { t := time.Now(); return &t }(), UpdatedAt: func() *time.Time { t := time.Now(); return &t }()},
	}
	db.Create(&customers)

	log.Println("数据初始化完成！")
	log.Println("───────────────")
	log.Println("默认登录账号：")
	log.Println("用户名：admin")
	log.Println("密码：admin123")
	log.Println("───────────────")
}
