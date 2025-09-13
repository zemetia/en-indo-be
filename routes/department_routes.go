package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func Department(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	departmentController := do.MustInvoke[controller.DepartmentController](injector)

	department := route.Group("/api/department")
	{
		department.POST("", middleware.Authenticate(jwtService, userService), departmentController.Create)
		department.GET("", middleware.Authenticate(jwtService, userService), departmentController.GetAll)
		department.GET("/:id", middleware.Authenticate(jwtService, userService), departmentController.GetByID)
		department.PUT("/:id", middleware.Authenticate(jwtService, userService), departmentController.Update)
		department.DELETE("/:id", middleware.Authenticate(jwtService, userService), departmentController.Delete)
	}
}
