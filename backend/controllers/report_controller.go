package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"skin-performance/config"
)

// PerformanceReport 业绩报表响应
type PerformanceReport struct {
	EmployeeID       uint    `json:"employee_id"`
	EmployeeName     string  `json:"employee_name"`
	EmployeeRole     string  `json:"employee_role"`
	MainPerformance  float64 `json:"main_performance"`
	CoPerformance    float64 `json:"co_performance"`
	NursePerformance float64 `json:"nurse_performance"`
	TotalPerformance float64 `json:"total_performance"`
}

// GetPerformanceReport 获取业绩报表
func GetPerformanceReport(c *gin.Context) {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	db := config.GetDB()

	var reports []PerformanceReport
	sql := `
		SELECT 
			e.id as employee_id,
			e.name as employee_name,
			e.role as employee_role,
			COALESCE(SUM(CASE WHEN vi.main_doctor_id = e.id THEN vi.amount ELSE 0 END), 0) as main_performance,
			COALESCE(SUM(CASE WHEN vi.co_doctor1_id = e.id THEN vi.amount * vi.co_ratio1 
							  WHEN vi.co_doctor2_id = e.id THEN vi.amount * vi.co_ratio2 ELSE 0 END), 0) as co_performance,
			COALESCE(SUM(CASE WHEN vi.nurse1_id = e.id OR vi.nurse2_id = e.id THEN vi.amount * 0.05 ELSE 0 END), 0) as nurse_performance
		FROM employees e
		LEFT JOIN visit_items vi ON (e.id = vi.main_doctor_id OR e.id = vi.co_doctor1_id OR e.id = vi.co_doctor2_id 
									 OR e.id = vi.nurse1_id OR e.id = vi.nurse2_id)
		LEFT JOIN visits v ON vi.visit_id = v.id AND v.visit_date >= ? AND v.visit_date <= ?
		WHERE e.is_active = ?
		GROUP BY e.id, e.name, e.role
		HAVING main_performance > 0 OR co_performance > 0 OR nurse_performance > 0
		ORDER BY (main_performance + co_performance + nurse_performance) DESC
	`

	if err := db.Raw(sql, dateFrom, dateTo+" 23:59:59", true).Scan(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var totalAmount float64
	db.Table("visit_items").Joins("JOIN visits ON visits.id = visit_items.visit_id").
		Where("visits.visit_date >= ? AND visits.visit_date <= ?", dateFrom, dateTo+" 23:59:59").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	for i := range reports {
		reports[i].TotalPerformance = reports[i].MainPerformance + reports[i].CoPerformance + reports[i].NursePerformance
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"date_from":    dateFrom,
			"date_to":      dateTo,
			"total_amount": totalAmount,
			"reports":      reports,
		},
	})
}

// GetEmployeePerformance 获取员工业绩
func GetEmployeePerformance(c *gin.Context) {
	employeeID, err := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的员工ID"})
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"employee_id": employeeID,
			"date_from":   dateFrom,
			"date_to":     dateTo,
		},
	})
}

// GetProjectPerformance 获取项目业绩
func GetProjectPerformance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    gin.H{},
	})
}
