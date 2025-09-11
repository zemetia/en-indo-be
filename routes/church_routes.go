package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func Church(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	churchController := do.MustInvoke[controller.ChurchController](injector)

	routes := route.Group("/api/church")
	{
		// Semua route church memerlukan autentikasi
		routes.POST("", middleware.Authenticate(jwtService, userService), churchController.Create)
		routes.GET("", churchController.GetAll)

		// Specific routes must come before parameterized routes
		routes.GET("/by-kabupaten/:id", churchController.GetByKabupatenID)
		routes.GET("/by-provinsi/:id", churchController.GetByProvinsiID)

		// Parameterized routes come last
		routes.GET("/:id", churchController.GetByID)
		routes.PUT("/:id", middleware.Authenticate(jwtService, userService), churchController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService, userService), churchController.Delete)
	}
}
