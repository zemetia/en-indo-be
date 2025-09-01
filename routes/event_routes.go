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

func EventRoutes(router *gin.RouterGroup, injector *do.Injector) {
	// Get dependencies from injector
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	
	// Create repositories and services
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventController := controller.NewEventController(eventService)
	
	// Event CRUD routes
	router.POST("/events", eventController.CreateEvent)
	router.GET("/events", eventController.ListEvents)
	router.GET("/events/:id", eventController.GetEvent)
	router.PUT("/events/:id", eventController.UpdateEvent)
	router.DELETE("/events/:id", eventController.DeleteEvent)
	
	// Recurring event management routes - three-tier modifications
	router.PUT("/events/:id/series", eventController.UpdateRecurringEvent)          // Update entire series
	router.PUT("/events/:id/occurrence", eventController.UpdateSingleOccurrence)   // Update single occurrence
	router.PUT("/events/:id/future", eventController.UpdateFutureOccurrences)      // Update this and future occurrences
	router.DELETE("/events/:id/occurrence", eventController.DeleteOccurrence)
	
	// Event occurrences routes
	router.GET("/events/:id/occurrences", eventController.GetEventOccurrences)
	router.GET("/events/:id/next", eventController.GetNextOccurrence)
	router.GET("/events/occurrences", eventController.GetOccurrencesInRange)
	
	// Validation and utility routes
	router.POST("/events/validate-recurrence", eventController.ValidateRecurrenceRule)
}