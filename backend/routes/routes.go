package routes

import (
	"github.com/gin-gonic/gin"
	"skin-performance/controllers"
	"skin-performance/middleware"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/login", controllers.Login)
		public.POST("/register", controllers.Register)
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 当前用户信息
		auth.GET("/user/info", controllers.GetCurrentUser)

		// 顾客管理
		auth.GET("/customers", controllers.ListCustomers)
		auth.GET("/customers/:id", controllers.GetCustomer)
		auth.POST("/customers", controllers.CreateCustomer)
		auth.PUT("/customers/:id", controllers.UpdateCustomer)
		auth.DELETE("/customers/:id", controllers.DeleteCustomer)

		// 员工管理（仅管理员可修改）
		auth.GET("/employees", controllers.ListEmployees)
		auth.GET("/employees/:id", controllers.GetEmployee)
		auth.POST("/employees", middleware.AdminMiddleware(), controllers.CreateEmployee)
		auth.PUT("/employees/:id", middleware.AdminMiddleware(), controllers.UpdateEmployee)
		auth.DELETE("/employees/:id", middleware.AdminMiddleware(), controllers.DeleteEmployee)

		// 项目管理（仅管理员可修改）
		auth.GET("/projects", controllers.ListProjects)
		auth.GET("/projects/:id", controllers.GetProject)
		auth.POST("/projects", middleware.AdminMiddleware(), controllers.CreateProject)
		auth.PUT("/projects/:id", middleware.AdminMiddleware(), controllers.UpdateProject)
		auth.DELETE("/projects/:id", middleware.AdminMiddleware(), controllers.DeleteProject)

		// 就诊管理
		auth.GET("/visits", controllers.ListVisits)
		auth.GET("/visits/:id", controllers.GetVisit)
		auth.POST("/visits", controllers.CreateVisit)
		auth.PUT("/visits/:id", controllers.UpdateVisit)
		auth.DELETE("/visits/:id", controllers.DeleteVisit)

		// 就诊明细
		auth.GET("/visit-items", controllers.ListVisitItems)
		auth.GET("/visit-items/:id", controllers.GetVisitItem)
		auth.POST("/visit-items", controllers.CreateVisitItem)
		auth.PUT("/visit-items/:id", controllers.UpdateVisitItem)
		auth.DELETE("/visit-items/:id", controllers.DeleteVisitItem)

		// 回访记录
		auth.GET("/revisit-records", controllers.ListRevisitRecords)
		auth.GET("/revisit-records/:id", controllers.GetRevisitRecord)
		auth.POST("/revisit-records", controllers.CreateRevisitRecord)
		auth.PUT("/revisit-records/:id", controllers.UpdateRevisitRecord)
		auth.DELETE("/revisit-records/:id", controllers.DeleteRevisitRecord)

		// 报表统计
		auth.GET("/reports/performance", controllers.GetPerformanceReport)
		auth.GET("/reports/employee-performance", controllers.GetEmployeePerformance)
		auth.GET("/reports/project-performance", controllers.GetProjectPerformance)
	}
}
