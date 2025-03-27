package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
)

func Department(router *gin.Engine, departmentController *controller.DepartmentController) {
	department := router.Group("/api/department")
	{
		department.POST("", middleware.Authenticate(jwtService), departmentController.Create)
		department.GET("", middleware.Authenticate(jwtService), departmentController.GetAll)
		department.GET("/:id", middleware.Authenticate(jwtService), departmentController.GetByID)
		department.GET("/church/:id", middleware.Authenticate(jwtService), departmentController.GetByChurchID)
		department.PUT("/:id", middleware.Authenticate(jwtService), departmentController.Update)
		department.DELETE("/:id", middleware.Authenticate(jwtService), departmentController.Delete)
	}
}
