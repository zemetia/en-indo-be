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
	churchController := do.MustInvoke[controller.ChurchController](injector)

	routes := route.Group("/api/church")
	{
		// Semua route church memerlukan autentikasi
		routes.POST("", middleware.Authenticate(jwtService), churchController.Create)
		routes.GET("", middleware.Authenticate(jwtService), churchController.GetAll)
		routes.GET("/:id", middleware.Authenticate(jwtService), churchController.GetByID)
		routes.GET("/kabupaten/:id", middleware.Authenticate(jwtService), churchController.GetByKabupatenID)
		routes.GET("/provinsi/:id", middleware.Authenticate(jwtService), churchController.GetByProvinsiID)
		routes.PUT("/:id", middleware.Authenticate(jwtService), churchController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), churchController.Delete)
	}
}
