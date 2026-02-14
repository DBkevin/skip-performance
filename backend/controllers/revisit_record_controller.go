package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListRevisitRecords 获取回访记录列表
func ListRevisitRecords(c *gin.Context) {
	var records []models.RevisitRecord
	query := config.GetDB().Model(&models.RevisitRecord{}).Preload("Nurse")

	// 筛选条件
	if nurseID := c.Query("nurse_id"); nurseID != "" {
		query = query.Where("nurse_id = ?", nurseID)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("date >= ?", dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("date <= ?", dateTo)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	var total int64
	query.Count(&total)

	if err := query.Order("date DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRevisitRecord 获取单个回访记录
func GetRevisitRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var record models.RevisitRecord
	if err := config.GetDB().Preload("Nurse").First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    record,
	})
}

// CreateRevisitRecord 创建回访记录
func CreateRevisitRecord(c *gin.Context) {
	var record models.RevisitRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if record.NurseID == 0 || record.Date.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "护士ID和日期为必填项"})
		return
	}

	now := time.Now()
	record.CreatedAt = &now
	record.UpdatedAt = &now

	if err := config.GetDB().Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    record,
	})
}

// UpdateRevisitRecord 更新回访记录
func UpdateRevisitRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var record models.RevisitRecord
	if err := config.GetDB().First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	var input models.RevisitRecord
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&record).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    record,
	})
}

// DeleteRevisitRecord 删除回访记录
func DeleteRevisitRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := config.GetDB().Delete(&models.RevisitRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
