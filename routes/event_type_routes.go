package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
)

func EventTypeRoutes(router *gin.RouterGroup, injector *do.Injector) {
	// Get controller from DI container
	eventTypeController := do.MustInvokeNamed[controller.EventTypeController](injector, constants.EventTypeController)

	// Event Type CRUD routes
	router.POST("/event-types", eventTypeController.Create)
	router.GET("/event-types", eventTypeController.GetAll)
	router.GET("/event-types/active", eventTypeController.GetActive)
	router.GET("/event-types/:id", eventTypeController.GetByID)
	router.PUT("/event-types/:id", eventTypeController.Update)
	router.DELETE("/event-types/:id", eventTypeController.Delete)
}
