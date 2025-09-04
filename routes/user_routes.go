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
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)
	userController := do.MustInvoke[controller.UserController](injector)

	routes := route.Group("/api/user")
	{
		// Public routes
		routes.POST("/register", userController.Register)
		routes.POST("/login", userController.Login)

		// Auth-related routes
		authRoutes := routes.Group("/auth")
		authRoutes.Use(middleware.Authenticate(jwtService, userService))
		{
			authRoutes.POST("/setup-password", userController.SetupPassword)
		}

		// Protected routes
		protected := routes.Group("")
		protected.Use(middleware.Authenticate(jwtService, userService))
		{
			protected.GET("", userController.GetAll)
			protected.GET("/:id", userController.GetByID)
			protected.GET("/email/:email", userController.GetByEmail)
			// protected.GET("/person/:person_id", userController.GetByPersonID)
			protected.PUT("/:id", userController.Update)
			protected.DELETE("/:id", userController.Delete)
			protected.POST("/:id/upload-profile-image", userController.UploadProfileImage)
			protected.PUT("/person/:person_id/toggle-status", userController.ToggleActivationStatus)
		}
	}
}
