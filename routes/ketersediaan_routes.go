package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func KetersediaanRoutes(router *gin.RouterGroup, injector *do.Injector) {
	// Get authentication services
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)

	// Get database and create repositories/services
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	ketersediaanRepo := repository.NewKetersediaanRepository(db)
	personRepo := repository.NewPersonRepository(db)
	eventRepo := repository.NewEventRepository(db)
	ketersediaanService := service.NewKetersediaanService(ketersediaanRepo, personRepo, eventRepo)
	ketersediaanController := controller.NewKetersediaanController(ketersediaanService)

	// Create authenticated route group
	routes := router.Group("/ketersediaan")
	routes.Use(middleware.Authenticate(jwtService, userService))
	{
		// Specific paths first
		routes.POST("/upsert", ketersediaanController.UpsertKetersediaan)
		routes.POST("/bulk", ketersediaanController.BulkCreateKetersediaan)
		routes.GET("/summary", ketersediaanController.GetEventAvailabilitySummary)
		routes.DELETE("/cleanup", ketersediaanController.CleanupKetersediaan)

		// Basic CRUD routes
		routes.POST("", ketersediaanController.CreateKetersediaan)
		routes.GET("", ketersediaanController.GetPersonAvailability)
		routes.GET("/:id", ketersediaanController.GetKetersediaan)
		routes.PUT("/:id", ketersediaanController.UpdateKetersediaan)
		routes.DELETE("/:id", ketersediaanController.DeleteKetersediaan)
	}

	// Event-specific availability routes
	// Note: /events/:id/availability route conflicts with existing event routes
	// Use /ketersediaan?eventId=xxx instead to get event availability
	// router.GET("/events/:id/availability", ketersediaanController.GetEventAvailability)
}
