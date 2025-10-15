package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zemetia/en-indo-be/constants"
	"github.com/zemetia/en-indo-be/controller"
	"github.com/zemetia/en-indo-be/middleware"
	"github.com/zemetia/en-indo-be/service"
)

func EventParticipantRoutes(router *gin.RouterGroup, injector *do.Injector) {
	// Get controller from DI container
	eventParticipantController := do.MustInvoke[*controller.EventParticipantController](injector)

	// Get services for authentication middleware
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	userService := do.MustInvokeNamed[service.UserService](injector, constants.UserService)

	// General participant routes
	participants := router.Group("/participants")
	{
		participants.GET("", eventParticipantController.ListParticipants)
	}

	// Event participant management routes (generic)
	eventParticipants := router.Group("/event-participants")
	{
		eventParticipants.POST("", eventParticipantController.RegisterParticipant)
		eventParticipants.POST("/bulk", eventParticipantController.BulkRegisterParticipants)
		eventParticipants.GET("/:id", eventParticipantController.GetParticipant)
		eventParticipants.PUT("/:id", eventParticipantController.UpdateParticipant)
		eventParticipants.DELETE("/:id", eventParticipantController.DeleteParticipant)
		eventParticipants.POST("/:id/checkin", eventParticipantController.CheckInParticipant)
		eventParticipants.POST("/checkin/bulk", eventParticipantController.BulkCheckIn)

		// QR scan check-in route (requires authentication to get person_id from JWT)
		eventParticipants.POST("/scan-qr", middleware.Authenticate(jwtService, userService), eventParticipantController.ScanQRCheckIn)
	}

	// Event-specific participant routes must be registered in EventRoutes
	// to avoid route conflicts with /events/:id
	// The following routes are handled in EventRoutes:
	// - GET /events/:id/participants
	// - GET /events/:id/attendance/report
}
