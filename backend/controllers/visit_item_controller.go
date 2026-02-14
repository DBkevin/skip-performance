package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListVisitItems 获取就诊明细列表
func ListVisitItems(c *gin.Context) {
	var items []models.VisitItem
	query := config.GetDB().Model(&models.VisitItem{}).
		Preload("Project").Preload("MainDoctor").
		Preload("Nurse1").Preload("Nurse2")

	if visitID := c.Query("visit_id"); visitID != "" {
		query = query.Where("visit_id = ?", visitID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      items,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetVisitItem 获取单个就诊明细
func GetVisitItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var item models.VisitItem
	if err := config.GetDB().Preload("Project").Preload("MainDoctor").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "明细记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    item,
	})
}

// CreateVisitItem 创建就诊明细
func CreateVisitItem(c *gin.Context) {
	var item models.VisitItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if item.VisitID == 0 || item.ProjectID == 0 || item.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 业绩分配计算
	calculatePerformance(&item)

	now := time.Now()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	tx := config.GetDB().Begin()
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	var totalAmount float64
	tx.Model(&models.VisitItem{}).Where("visit_id = ?", item.VisitID).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	tx.Model(&models.Visit{}).Where("id = ?", item.VisitID).Update("total_amount", totalAmount)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    item,
	})
}

// UpdateVisitItem 更新就诊明细
func UpdateVisitItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var item models.VisitItem
	if err := config.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "明细记录不存在"})
		return
	}

	var input models.VisitItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	calculatePerformance(&input)
	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&item).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    item,
	})
}

// DeleteVisitItem 删除就诊明细
func DeleteVisitItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var item models.VisitItem
	if err := config.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "明细记录不存在"})
		return
	}

	visitID := item.VisitID
	tx := config.GetDB().Begin()
	tx.Delete(&item)

	var totalAmount float64
	tx.Model(&models.VisitItem{}).Where("visit_id = ?", visitID).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	tx.Model(&models.Visit{}).Where("id = ?", visitID).Update("total_amount", totalAmount)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// calculatePerformance 计算业绩分配
func calculatePerformance(item *models.VisitItem) {
	coTotalRatio := item.CoRatio1 + item.CoRatio2
	if coTotalRatio > 1 {
		coTotalRatio = 1
	}
	item.MainDoctorPerformance = item.Amount * (1 - coTotalRatio)

	if item.CoDoctor1ID != nil {
		item.CoDoctor1Performance = item.Amount * item.CoRatio1
	}
	if item.CoDoctor2ID != nil {
		item.CoDoctor2Performance = item.Amount * item.CoRatio2
	}

	nurseRatio := 0.05
	if item.Nurse1ID != nil {
		item.Nurse1Performance = item.Amount * nurseRatio
	}
	if item.Nurse2ID != nil {
		item.Nurse2Performance = item.Amount * nurseRatio
	}
}
