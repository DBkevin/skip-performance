package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
	"skin-performance/models"
)

// ListVisits 获取就诊列表
func ListVisits(c *gin.Context) {
	var visits []models.Visit
	query := config.GetDB().Model(&models.Visit{}).
		Preload("Customer").
		Preload("Consultant").
		Preload("Items").
		Preload("Items.Project")

	// 筛选条件
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if consultantID := c.Query("consultant_id"); consultantID != "" {
		query = query.Where("consultant_id = ?", consultantID)
	}
	if visitID := c.Query("visit_id"); visitID != "" {
		query = query.Where("visit_id LIKE ?", "%"+visitID+"%")
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		query = query.Where("visit_date >= ?", dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		query = query.Where("visit_date <= ?", dateTo+" 23:59:59")
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

	if err := query.Order("visit_date DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&visits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      visits,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetVisit 获取单个就诊
func GetVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var visit models.Visit
	if err := config.GetDB().Preload("Customer").Preload("Consultant").
		Preload("Items").Preload("Items.Project").
		Preload("Items.MainDoctor").Preload("Items.Nurse1").Preload("Items.Nurse2").
		First(&visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "就诊记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    visit,
	})
}

// CreateVisit 创建就诊
func CreateVisit(c *gin.Context) {
	var visit models.Visit
	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 验证必填字段
	if visit.VisitID == "" || visit.CustomerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "单据号和顾客ID为必填项"})
		return
	}

	// 检查单据号是否已存在
	var existingVisit models.Visit
	if err := config.GetDB().Where("visit_id = ?", visit.VisitID).First(&existingVisit).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "单据号已存在"})
		return
	}

	// 设置创建时间
	now := time.Now()
	visit.CreatedAt = &now
	visit.UpdatedAt = &now

	if err := config.GetDB().Create(&visit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	// 重新加载关联数据
	config.GetDB().Preload("Customer").Preload("Consultant").First(&visit, visit.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    visit,
	})
}

// UpdateVisit 更新就诊
func UpdateVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var visit models.Visit
	if err := config.GetDB().First(&visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "就诊记录不存在"})
		return
	}

	var input models.Visit
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 更新时间
	now := time.Now()
	input.UpdatedAt = &now

	if err := config.GetDB().Model(&visit).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	// 更新总金额（根据明细自动计算）
	var totalAmount float64
	config.GetDB().Model(&models.VisitItem{}).Where("visit_id = ?", id).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	config.GetDB().Model(&visit).Update("total_amount", totalAmount)

	// 重新加载
	config.GetDB().Preload("Customer").Preload("Consultant").Preload("Items").First(&visit, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    visit,
	})
}

// DeleteVisit 删除就诊（软删除）
func DeleteVisit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	// 使用事务删除就诊及其明细
	tx := config.GetDB().Begin()
	if err := tx.Where("visit_id = ?", id).Delete(&models.VisitItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除明细失败"})
		return
	}

	if err := tx.Delete(&models.Visit{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
