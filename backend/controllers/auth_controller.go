package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
	"skin-performance/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string       `json:"token"`
	User     models.User  `json:"user"`
	Employee *models.Employee `json:"employee,omitempty"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	var user models.User
	if err := config.GetDB().Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	// 加载员工信息
	var employee models.Employee
	if user.EmployeeID != nil {
		config.GetDB().First(&employee, *user.EmployeeID)
	}

	// 不返回密码
	user.Password = ""

	response := LoginResponse{
		Token: token,
		User:  user,
	}
	if user.EmployeeID != nil {
		response.Employee = &employee
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	username, _ := c.Get("username")
	employeeID, _ := c.Get("employeeID")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"user_id":     userID,
			"username":    username,
			"role":        role,
			"employee_id": employeeID,
		},
	})
}
