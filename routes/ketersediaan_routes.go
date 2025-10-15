package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/repository"
	"github.com/zemetia/en-indo-be/service"
	"gorm.io/gorm"
)

func KetersediaanRoutes(router *gin.RouterGroup, injector *do.Injector) {
	// Get dependencies from injector
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)

	// Create repositories and services
	ketersediaanRepo := repository.NewKetersediaanRepository(db)
	personRepo := repository.NewPersonRepository(db)
	eventRepo := repository.NewEventRepository(db)

	ketersediaanService := service.NewKetersediaanService(ketersediaanRepo, personRepo, eventRepo)

	// Create controller
	ketersediaanController := controller.NewKetersediaanController(ketersediaanService)

	// Ketersediaan routes - specific paths first
	router.POST("/ketersediaan/upsert", ketersediaanController.UpsertKetersediaan)
	router.POST("/ketersediaan/bulk", ketersediaanController.BulkCreateKetersediaan)
	router.GET("/ketersediaan/summary", ketersediaanController.GetEventAvailabilitySummary)

	// Basic CRUD routes
	router.POST("/ketersediaan", ketersediaanController.CreateKetersediaan)
	router.GET("/ketersediaan", ketersediaanController.GetPersonAvailability)
	router.GET("/ketersediaan/:id", ketersediaanController.GetKetersediaan)
	router.PUT("/ketersediaan/:id", ketersediaanController.UpdateKetersediaan)
	router.DELETE("/ketersediaan/:id", ketersediaanController.DeleteKetersediaan)

	// Event-specific availability routes
	// Note: /events/:id/availability route conflicts with existing event routes
	// Use /ketersediaan?eventId=xxx instead to get event availability
	// router.GET("/events/:id/availability", ketersediaanController.GetEventAvailability)
}
