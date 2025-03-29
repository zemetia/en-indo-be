package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func User(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userController := do.MustInvoke[controller.UserController](injector)

	routes := route.Group("/api/user")
	{
		// Public routes
		routes.POST("/register", userController.Register)
		routes.POST("/login", userController.Login)

		// Protected routes
		protected := routes.Group("")
		protected.Use(middleware.Authenticate(jwtService))
		{
			protected.GET("", userController.GetAll)
			protected.GET("/:id", userController.GetByID)
			protected.GET("/email/:email", userController.GetByEmail)
			// protected.GET("/person/:person_id", userController.GetByPersonID)
			protected.PUT("/:id", userController.Update)
			protected.DELETE("/:id", userController.Delete)
			protected.POST("/:id/upload-profile-image", userController.UploadProfileImage)
		}
	}
}
