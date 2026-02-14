package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
	"skin-performance/utils"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=32"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Phone    *string `json:"phone"`
}

// Register 用户注册（初始管理员创建）
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.GetDB().Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	// 创建员工记录
	employee := models.Employee{
		Name:      req.Name,
		Role:      req.Role,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: func() *time.Time { t := time.Now(); return &t }(),
		UpdatedAt: func() *time.Time { t := time.Now(); return &t }(),
	}

	if err := config.GetDB().Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建员工失败"})
		return
	}

	// 创建用户记录
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	user := models.User{
		Username:   req.Username,
		Password:   hashedPassword,
		EmployeeID: &employee.ID,
		Role:       req.Role,
		IsActive:   true,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建用户失败"})
		return
	}

	// 不返回密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"user":     user,
			"employee": employee,
		},
	})
}
