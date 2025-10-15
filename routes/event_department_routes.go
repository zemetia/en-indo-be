package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func EventDepartment(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	eventDepartmentController := do.MustInvoke[controller.EventDepartmentController](injector)

	routes := route.Group("/api/events")
	{
		routes.GET("/:id/departments", middleware.Authenticate(jwtService, userService), eventDepartmentController.GetEventDepartments)
		routes.PUT("/:id/departments", middleware.Authenticate(jwtService, userService), eventDepartmentController.UpdateEventDepartments)
	}
}
