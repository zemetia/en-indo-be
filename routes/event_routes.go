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
	eventPICRepo := repository.NewEventPICRepository(db)
	eventService := service.NewEventService(eventRepo, eventPICRepo)
	eventPICService := service.NewEventPICService(eventPICRepo, eventRepo)
	
	// Create controllers
	eventController := controller.NewEventController(eventService)
	eventPICController := controller.NewEventPICController(eventPICService)
	eventPICRoleController := controller.NewEventPICRoleController(eventPICService)

	// Event CRUD routes - keep simple ones here
	router.POST("/events", eventController.CreateEvent)
	router.GET("/events", eventController.ListEvents)
	
	// Validation and utility routes - no params, put before parameterized routes
	router.POST("/events/validate-recurrence", eventController.ValidateRecurrenceRule)
	router.GET("/events/occurrences", eventController.GetOccurrencesInRange)

	// Event PIC management routes - put more specific paths first
	router.POST("/events/:id/pics/bulk", eventPICController.BulkAssignEventPICs)
	router.POST("/events/:id/pics/transfer", eventPICController.TransferPICRole)
	router.GET("/events/:id/pics/active", eventPICController.GetActivePICsForEvent)
	router.GET("/events/:id/pics/primary", eventPICController.GetPrimaryPICForEvent)
	router.GET("/events/:id/pics/history", eventPICController.GetEventPICHistory)
	router.GET("/events/:id/pics/validate/:personId", eventPICController.ValidatePICPermissions)
	router.POST("/events/:id/pics", eventPICController.CreateEventPIC)
	router.GET("/events/:id/pics", eventPICController.GetEventPICs)

	// Event occurrences routes - specific paths first
	router.GET("/events/:id/occurrences", eventController.GetEventOccurrences)
	router.GET("/events/:id/next", eventController.GetNextOccurrence)
	
	// Recurring event management routes - three-tier modifications
	router.PUT("/events/:id/series", eventController.UpdateRecurringEvent)       // Update entire series
	router.PUT("/events/:id/occurrence", eventController.UpdateSingleOccurrence) // Update single occurrence
	router.PUT("/events/:id/future", eventController.UpdateFutureOccurrences)    // Update this and future occurrences
	router.DELETE("/events/:id/occurrence", eventController.DeleteOccurrence)
	
	// Basic CRUD routes with :id param - put at end to avoid conflicts
	router.GET("/events/:id", eventController.GetEvent)
	router.PUT("/events/:id", eventController.UpdateEvent)
	router.DELETE("/events/:id", eventController.DeleteEvent)
	
	// Individual EventPIC operations
	router.GET("/event-pics/:id", eventPICController.GetEventPIC)
	router.PUT("/event-pics/:id", eventPICController.UpdateEventPIC)
	router.DELETE("/event-pics/:id", eventPICController.DeleteEventPIC)
	router.GET("/event-pics", eventPICController.ListEventPICs)
	router.GET("/event-pics/expiring", eventPICController.GetExpiringPICs)
	
	// Person-centric PIC routes
	router.GET("/persons/:personId/event-pics", eventPICController.GetPersonPICs)
	router.GET("/persons/:personId/event-pics/active", eventPICController.GetActivePersonPICs)
	router.GET("/persons/:personId/event-pics/history", eventPICController.GetPersonPICHistory)
	
	// Event PIC Role management routes
	router.POST("/event-pic-roles", eventPICRoleController.CreateEventPICRole)
	router.GET("/event-pic-roles", eventPICRoleController.ListEventPICRoles)
	router.GET("/event-pic-roles/:id", eventPICRoleController.GetEventPICRole)
	router.PUT("/event-pic-roles/:id", eventPICRoleController.UpdateEventPICRole)
	router.DELETE("/event-pic-roles/:id", eventPICRoleController.DeleteEventPICRole)
}
