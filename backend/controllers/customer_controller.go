package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListCustomers 获取顾客列表
func ListCustomers(c *gin.Context) {
	var customers []models.Customer
	query := config.GetDB().Model(&models.Customer{})

	// 搜索条件
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone := c.Query("phone"); phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if customerType := c.Query("customer_type"); customerType != "" {
		query = query.Where("customer_type = ?", customerType)
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

	if err := query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      customers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetCustomer 获取单个顾客
func GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var customer models.Customer
	if err := config.GetDB().First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "顾客不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    customer,
	})
}

// CreateCustomer 创建顾客
func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 设置创建时间
	now := time.Now()
	customer.CreatedAt = &now
	customer.UpdatedAt = &now

	if err := config.GetDB().Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    customer,
	})
}

// UpdateCustomer 更新顾客
func UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var customer models.Customer
	if err := config.GetDB().First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "顾客不存在"})
		return
	}

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 更新时间
	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&customer).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    customer,
	})
}

// DeleteCustomer 删除顾客（软删除）
func DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := config.GetDB().Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
