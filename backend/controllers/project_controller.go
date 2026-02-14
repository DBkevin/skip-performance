package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListProjects 获取项目列表
func ListProjects(c *gin.Context) {
	var projects []models.Project
	query := config.GetDB().Model(&models.Project{})

	// 搜索条件
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	// 只查询启用的项目
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

	if err := query.Limit(pageSize).Offset((page - 1) * pageSize).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      projects,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProject 获取单个项目
func GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var project models.Project
	if err := config.GetDB().First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "项目不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    project,
	})
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 验证必填字段
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "项目名称为必填项"})
		return
	}

	// 设置创建时间
	now := time.Now()
	project.CreatedAt = &now
	project.UpdatedAt = &now
	if !project.IsActive {
		project.IsActive = true
	}

	if err := config.GetDB().Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    project,
	})
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var project models.Project
	if err := config.GetDB().First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "项目不存在"})
		return
	}

	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 更新时间
	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&project).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    project,
	})
}

// DeleteProject 删除项目（软删除）
func DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := config.GetDB().Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
