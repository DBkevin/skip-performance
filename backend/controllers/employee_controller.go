package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListEmployees 获取员工列表
func ListEmployees(c *gin.Context) {
	var employees []models.Employee
	query := config.GetDB().Model(&models.Employee{})

	// 搜索条件
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if department := c.Query("department"); department != "" {
		query = query.Where("department = ?", department)
	}
	// 只查询在职员工
	if c.Query("active_only") == "true" {
		query = query.Where("is_active = ?", true)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	if err := query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      employees,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetEmployee 获取单个员工
func GetEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var employee models.Employee
	if err := config.GetDB().First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    employee,
	})
}

// CreateEmployee 创建员工
func CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 验证必填字段
	if employee.Name == "" || employee.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "姓名和角色为必填项"})
		return
	}

	// 设置创建时间
	now := time.Now()
	employee.CreatedAt = &now
	employee.UpdatedAt = &now

	if err := config.GetDB().Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    employee,
	})
}

// UpdateEmployee 更新员工
func UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var employee models.Employee
	if err := config.GetDB().First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 更新时间
	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&employee).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    employee,
	})
}

// DeleteEmployee 删除员工（软删除）
func DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := config.GetDB().Delete(&models.Employee{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
