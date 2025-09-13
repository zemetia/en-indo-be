package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventPICController struct {
	eventPICService service.EventPICService
}

func NewEventPICController(eventPICService service.EventPICService) *EventPICController {
	return &EventPICController{
		eventPICService: eventPICService,
	}
}

// CreateEventPIC godoc
// @Summary Assign a PIC to an event
// @Description Assign a Person in Charge (PIC) to a specific event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param pic body dto.CreateEventPICRequest true "Event PIC data"
// @Success 201 {object} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics [post]
func (c *EventPICController) CreateEventPIC(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	var req dto.CreateEventPICRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Get actual user ID from authentication context
	createdBy := uuid.New() // Placeholder

	eventPIC, err := c.eventPICService.CreateEventPIC(eventID, &req, createdBy)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "event not found: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "person already assigned as PIC for this event" ||
			err.Error() == "event already has a primary PIC" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to create event PIC",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, eventPIC)
}

// GetEventPIC godoc
// @Summary Get event PIC by ID
// @Description Get a specific event PIC assignment by its ID
// @Tags event-pics
// @Accept json
// @Produce json
// @Param id path string true "Event PIC ID"
// @Success 200 {object} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /event-pics/{id} [get]
func (c *EventPICController) GetEventPIC(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC ID format",
		})
		return
	}

	eventPIC, err := c.eventPICService.GetEventPIC(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event PIC: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to get event PIC",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPIC)
}

// UpdateEventPIC godoc
// @Summary Update event PIC
// @Description Update an existing event PIC assignment
// @Tags event-pics
// @Accept json
// @Produce json
// @Param id path string true "Event PIC ID"
// @Param pic body dto.UpdateEventPICRequest true "Event PIC update data"
// @Success 200 {object} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pics/{id} [put]
func (c *EventPICController) UpdateEventPIC(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC ID format",
		})
		return
	}

	var req dto.UpdateEventPICRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Get actual user ID from authentication context
	updatedBy := uuid.New() // Placeholder

	eventPIC, err := c.eventPICService.UpdateEventPIC(id, &req, updatedBy)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event PIC: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "event already has a primary PIC" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update event PIC",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPIC)
}

// DeleteEventPIC godoc
// @Summary Remove event PIC
// @Description Remove a PIC assignment from an event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param id path string true "Event PIC ID"
// @Param reason query string false "Reason for removal"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pics/{id} [delete]
func (c *EventPICController) DeleteEventPIC(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event PIC ID format",
		})
		return
	}

	reason := ctx.Query("reason")
	if reason == "" {
		reason = "PIC removed"
	}

	// TODO: Get actual user ID from authentication context
	deletedBy := uuid.New() // Placeholder

	err = c.eventPICService.DeleteEventPIC(id, deletedBy, reason)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event PIC: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to delete event PIC",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetEventPICs godoc
// @Summary Get PICs for an event
// @Description Get all PIC assignments for a specific event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics [get]
func (c *EventPICController) GetEventPICs(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	eventPICs, err := c.eventPICService.GetPICsByEventID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPICs)
}

// GetActivePICsForEvent godoc
// @Summary Get active PICs for an event
// @Description Get all active PIC assignments for a specific event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/active [get]
func (c *EventPICController) GetActivePICsForEvent(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	eventPICs, err := c.eventPICService.GetActivePICsByEventID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get active event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPICs)
}

// GetPrimaryPICForEvent godoc
// @Summary Get primary PIC for an event
// @Description Get the primary PIC assignment for a specific event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/primary [get]
func (c *EventPICController) GetPrimaryPICForEvent(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	eventPIC, err := c.eventPICService.GetPrimaryPICByEventID(eventID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to get primary event PIC",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPIC)
}

// GetPersonPICs godoc
// @Summary Get PICs for a person
// @Description Get all PIC assignments for a specific person
// @Tags event-pics
// @Accept json
// @Produce json
// @Param personId path string true "Person ID"
// @Success 200 {array} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /persons/{personId}/event-pics [get]
func (c *EventPICController) GetPersonPICs(ctx *gin.Context) {
	personIDStr := ctx.Param("personId")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid person ID format",
		})
		return
	}

	eventPICs, err := c.eventPICService.GetPICsByPersonID(personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get person event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPICs)
}

// GetActivePersonPICs godoc
// @Summary Get active PICs for a person
// @Description Get all active PIC assignments for a specific person
// @Tags event-pics
// @Accept json
// @Produce json
// @Param personId path string true "Person ID"
// @Success 200 {array} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /persons/{personId}/event-pics/active [get]
func (c *EventPICController) GetActivePersonPICs(ctx *gin.Context) {
	personIDStr := ctx.Param("personId")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid person ID format",
		})
		return
	}

	eventPICs, err := c.eventPICService.GetActivePICsByPersonID(personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get active person event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPICs)
}

// ListEventPICs godoc
// @Summary List event PICs with filters
// @Description Get a paginated list of event PICs with optional filters
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId query string false "Filter by event ID"
// @Param personId query string false "Filter by person ID"
// @Param role query string false "Filter by role"
// @Param isActive query bool false "Filter by active status"
// @Param isPrimary query bool false "Filter by primary status"
// @Param search query string false "Search in person name, role"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} dto.EventPICListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pics [get]
func (c *EventPICController) ListEventPICs(ctx *gin.Context) {
	var req dto.EventPICFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	eventPICs, err := c.eventPICService.ListEventPICs(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, eventPICs)
}

// BulkAssignEventPICs godoc
// @Summary Bulk assign PICs to an event
// @Description Assign multiple PICs to an event in one operation
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param pics body dto.BulkAssignEventPICRequest true "Bulk PIC assignment data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/bulk [post]
func (c *EventPICController) BulkAssignEventPICs(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	var req dto.BulkAssignEventPICRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Get actual user ID from authentication context
	createdBy := uuid.New() // Placeholder

	err = c.eventPICService.AssignMultiplePICs(eventID, &req, createdBy)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "event not found: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "cannot assign multiple primary PICs" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to bulk assign event PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Event PICs assigned successfully",
	})
}

// TransferPICRole godoc
// @Summary Transfer PIC role
// @Description Transfer PIC responsibilities from one person to another
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param transfer body dto.TransferEventPICRequest true "Transfer data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/transfer [post]
func (c *EventPICController) TransferPICRole(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	var req dto.TransferEventPICRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Get actual user ID from authentication context
	changedBy := uuid.New() // Placeholder

	err = c.eventPICService.TransferPICRole(eventID, &req, changedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to transfer PIC role",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "PIC role transferred successfully",
	})
}

// GetEventPICHistory godoc
// @Summary Get PIC history for an event
// @Description Get the history of PIC changes for a specific event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} dto.EventPICHistoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/history [get]
func (c *EventPICController) GetEventPICHistory(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	history, err := c.eventPICService.GetEventPICHistory(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get event PIC history",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}

// GetPersonPICHistory godoc
// @Summary Get PIC history for a person
// @Description Get the history of PIC assignments for a specific person
// @Tags event-pics
// @Accept json
// @Produce json
// @Param personId path string true "Person ID"
// @Success 200 {array} dto.EventPICHistoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /persons/{personId}/event-pics/history [get]
func (c *EventPICController) GetPersonPICHistory(ctx *gin.Context) {
	personIDStr := ctx.Param("personId")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid person ID format",
		})
		return
	}

	history, err := c.eventPICService.GetPersonPICHistory(personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get person PIC history",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, history)
}

// GetExpiringPICs godoc
// @Summary Get expiring PICs
// @Description Get PICs that are expiring within a specified number of days
// @Tags event-pics
// @Accept json
// @Produce json
// @Param days query int false "Number of days to look ahead (default: 30)"
// @Success 200 {array} dto.EventPICResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /event-pics/expiring [get]
func (c *EventPICController) GetExpiringPICs(ctx *gin.Context) {
	daysStr := ctx.Query("days")
	days := 30 // Default to 30 days
	if daysStr != "" {
		parsed, err := strconv.Atoi(daysStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid days parameter",
			})
			return
		}
		days = parsed
	}

	pics, err := c.eventPICService.GetExpiringPICs(days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get expiring PICs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pics)
}

// ValidatePICPermissions godoc
// @Summary Validate PIC permissions
// @Description Check if a person has specific permissions for an event
// @Tags event-pics
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param personId path string true "Person ID"
// @Param action query string true "Action to validate (edit, delete, assign_pic)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{eventId}/pics/validate/{personId} [get]
func (c *EventPICController) ValidatePICPermissions(ctx *gin.Context) {
	eventIDStr := ctx.Param("eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	personIDStr := ctx.Param("personId")
	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid person ID format",
		})
		return
	}

	action := ctx.Query("action")
	if action == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Action parameter is required",
		})
		return
	}

	hasPermission, err := c.eventPICService.ValidateEventPICPermissions(eventID, personID, action)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to validate permissions",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"hasPermission": hasPermission,
		"action":        action,
	})
}