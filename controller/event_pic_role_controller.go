package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventPICRoleController struct {
	eventPICService service.EventPICService
}

func NewEventPICRoleController(eventPICService service.EventPICService) *EventPICRoleController {
	return &EventPICRoleController{
		eventPICService: eventPICService,
	}
}

// CreateEventPICRole godoc
// @Summary Create a new event PIC role
// @Description Create a new predefined role for event PICs
// @Tags event-pic-roles
// @Accept json
// @Produce json
// @Param role body dto.CreateEventPICRoleRequest true "Event PIC role data"
// @Success 201 {object} dto.EventPICRoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pic-roles [post]
func (c *EventPICRoleController) CreateEventPICRole(ctx *gin.Context) {
	var req dto.CreateEventPICRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	role, err := c.eventPICService.CreateEventPICRole(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "role name already exists" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to create event PIC role",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, role)
}

// GetEventPICRole godoc
// @Summary Get event PIC role by ID
// @Description Get a specific event PIC role by its ID
// @Tags event-pic-roles
// @Accept json
// @Produce json
// @Param id path string true "Event PIC role ID"
// @Success 200 {object} dto.EventPICRoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /event-pic-roles/{id} [get]
func (c *EventPICRoleController) GetEventPICRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC role ID format",
		})
		return
	}

	role, err := c.eventPICService.GetEventPICRole(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event PIC role: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to get event PIC role",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

// UpdateEventPICRole godoc
// @Summary Update event PIC role
// @Description Update an existing event PIC role
// @Tags event-pic-roles
// @Accept json
// @Produce json
// @Param id path string true "Event PIC role ID"
// @Param role body dto.UpdateEventPICRoleRequest true "Event PIC role update data"
// @Success 200 {object} dto.EventPICRoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pic-roles/{id} [put]
func (c *EventPICRoleController) UpdateEventPICRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC role ID format",
		})
		return
	}

	var req dto.UpdateEventPICRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	role, err := c.eventPICService.UpdateEventPICRole(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event PIC role: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "role name already exists" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update event PIC role",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

// DeleteEventPICRole godoc
// @Summary Delete event PIC role
// @Description Delete an event PIC role
// @Tags event-pic-roles
// @Accept json
// @Produce json
// @Param id path string true "Event PIC role ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pic-roles/{id} [delete]
func (c *EventPICRoleController) DeleteEventPICRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC role ID format",
		})
		return
	}

	err = c.eventPICService.DeleteEventPICRole(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to delete event PIC role: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to delete event PIC role",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListEventPICRoles godoc
// @Summary List event PIC roles
// @Description Get a paginated list of event PIC roles with optional search
// @Tags event-pic-roles
// @Accept json
// @Produce json
// @Param search query string false "Search in role name or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} dto.EventPICRoleListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pic-roles [get]
func (c *EventPICRoleController) ListEventPICRoles(ctx *gin.Context) {
	search := ctx.Query("search")
	
	pageStr := ctx.Query("page")
	page := 1
	if pageStr != "" {
		parsed, err := strconv.Atoi(pageStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid page parameter",
			})
			return
		}
		page = parsed
	}

	limitStr := ctx.Query("limit")
	limit := 20
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid limit parameter",
			})
			return
		}
		limit = parsed
	}

	roles, err := c.eventPICService.ListEventPICRoles(page, limit, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list event PIC roles",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}