package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func Person(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	personController := do.MustInvoke[controller.PersonController](injector)

	routes := route.Group("/api/person")
	{
		// Semua route person memerlukan autentikasi
		routes.POST("", middleware.Authenticate(jwtService, userService), personController.Create)
		routes.GET("", middleware.Authenticate(jwtService, userService), personController.GetAll)
		routes.GET("/by-pic-lifegroup-churches", middleware.Authenticate(jwtService, userService), personController.GetByPICLifegroupChurches)
		routes.GET("/:id", middleware.Authenticate(jwtService, userService), personController.GetByID)
		routes.GET("/user/:user_id", middleware.Authenticate(jwtService, userService), personController.GetByUserID)
		routes.PUT("/:id", middleware.Authenticate(jwtService, userService), personController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService, userService), personController.Delete)

	}
}
