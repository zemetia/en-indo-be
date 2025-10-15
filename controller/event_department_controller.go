package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventDepartmentController interface {
	GetEventDepartments(ctx *gin.Context)
	UpdateEventDepartments(ctx *gin.Context)
}

type eventDepartmentController struct {
	eventDepartmentService service.EventDepartmentService
}

func NewEventDepartmentController(eventDepartmentService service.EventDepartmentService) EventDepartmentController {
	return &eventDepartmentController{
		eventDepartmentService: eventDepartmentService,
	}
}

func (c *eventDepartmentController) GetEventDepartments(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID format",
			"error":   err.Error(),
		})
		return
	}

	departments, err := c.eventDepartmentService.GetEventDepartments(ctx, eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get event departments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get event departments",
		"data":    departments,
	})
}

func (c *eventDepartmentController) UpdateEventDepartments(ctx *gin.Context) {
	eventID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid event ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.UpdateEventDepartmentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get data from request body",
			"error":   err.Error(),
		})
		return
	}

	if err := c.eventDepartmentService.UpdateEventDepartments(ctx, eventID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update event departments",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success update event departments",
	})
}
