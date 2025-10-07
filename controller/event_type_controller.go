package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventTypeController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetActive(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type eventTypeController struct {
	eventTypeService *service.EventTypeService
}

func NewEventTypeController(eventTypeService *service.EventTypeService) EventTypeController {
	return &eventTypeController{
		eventTypeService: eventTypeService,
	}
}

func (c *eventTypeController) Create(ctx *gin.Context) {
	log.Printf("[INFO] EventType Controller: Received create event type request")

	var req dto.EventTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] EventType Controller: Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[INFO] EventType Controller: Creating event type with data:")
	log.Printf("  - DisplayName: %s", req.DisplayName)
	log.Printf("  - Description: %s", req.Description)
	log.Printf("  - Color: %s", req.Color)
	log.Printf("  - SortOrder: %d", req.SortOrder)

	// Validate required fields
	if req.DisplayName == "" {
		log.Printf("[ERROR] EventType Controller: DisplayName is required but empty")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Display name is required",
			"error":   "display_name field cannot be empty",
		})
		return
	}

	eventType, err := c.eventTypeService.Create(&req)
	if err != nil {
		log.Printf("[ERROR] EventType Controller: Service error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event type",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[INFO] EventType Controller: Successfully created event type: %s", eventType.Name)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success create event type",
		"data":    eventType,
	})
}

func (c *eventTypeController) GetAll(ctx *gin.Context) {
	eventTypes, err := c.eventTypeService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get event types",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get all event types",
		"data":    eventTypes,
		"count":   len(eventTypes),
	})
}

func (c *eventTypeController) GetActive(ctx *gin.Context) {
	eventTypes, err := c.eventTypeService.GetActive()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get active event types",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get active event types",
		"data":    eventTypes,
		"count":   len(eventTypes),
	})
}

func (c *eventTypeController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	eventType, err := c.eventTypeService.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Event type not found",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get event type",
		"data":    eventType,
	})
}

func (c *eventTypeController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.EventTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	eventType, err := c.eventTypeService.Update(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update event type",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success update event type",
		"data":    eventType,
	})
}

func (c *eventTypeController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.eventTypeService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete event type",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success delete event type",
	})
}
