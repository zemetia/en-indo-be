package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with optional recurrence rules
// @Tags events
// @Accept json
// @Produce json
// @Param event body dto.CreateEventRequest true "Event creation data"
// @Success 201 {object} dto.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events [post]
func (c *EventController) CreateEvent(ctx *gin.Context) {
	var req dto.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	event, err := c.eventService.CreateEvent(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create event",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusCreated, event)
}

// GetEvent godoc
// @Summary Get event by ID
// @Description Get a single event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /events/{id} [get]
func (c *EventController) GetEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	event, err := c.eventService.GetEvent(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Event not found",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, event)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an existing event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body dto.UpdateEventRequest true "Event update data"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [put]
func (c *EventController) UpdateEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	event, err := c.eventService.UpdateEvent(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update event",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event and all its recurrence data
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [delete]
func (c *EventController) DeleteEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	err = c.eventService.DeleteEvent(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to delete event: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to delete event",
			"details": err.Error(),
		})
		return
	}
	
	ctx.Status(http.StatusNoContent)
}

// ListEvents godoc
// @Summary List events with filters
// @Description Get a paginated list of events with optional filters
// @Tags events
// @Accept json
// @Produce json
// @Param type query string false "Event type (event, ibadah, spiritual_journey)"
// @Param isPublic query bool false "Filter by public/private events"
// @Param search query string false "Search in title, description, or location"
// @Param startDate query string false "Filter events from this date (YYYY-MM-DD)"
// @Param endDate query string false "Filter events until this date (YYYY-MM-DD)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} dto.EventListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events [get]
func (c *EventController) ListEvents(ctx *gin.Context) {
	var req dto.EventFilterRequest
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
	
	events, err := c.eventService.ListEvents(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list events",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, events)
}

// UpdateRecurringEvent godoc
// @Summary Update recurring event series or specific occurrence
// @Description Update either the entire series, single occurrence, or this and future occurrences
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param update body dto.UpdateRecurringEventRequest true "Recurring event update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/series [put]
func (c *EventController) UpdateRecurringEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.UpdateRecurringEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err = c.eventService.UpdateRecurringEvent(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "event is not recurring" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update recurring event",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recurring event updated successfully",
	})
}

// DeleteOccurrence godoc
// @Summary Delete specific occurrence or future occurrences
// @Description Delete a single occurrence or all future occurrences of a recurring event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param delete body dto.DeleteOccurrenceRequest true "Delete occurrence data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/occurrence [delete]
func (c *EventController) DeleteOccurrence(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.DeleteOccurrenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err = c.eventService.DeleteOccurrence(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "event is not recurring" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to delete occurrence",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Occurrence deleted successfully",
	})
}

// GetEventOccurrences godoc
// @Summary Get event occurrences in date range
// @Description Get all occurrences of a specific event within a date range
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param startDate query string true "Start date (YYYY-MM-DD)"
// @Param endDate query string true "End date (YYYY-MM-DD)"
// @Param timezone query string false "Timezone for results"
// @Success 200 {array} dto.EventOccurrenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/occurrences [get]
func (c *EventController) GetEventOccurrences(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.GetEventOccurrencesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}
	
	occurrences, err := c.eventService.GetEventOccurrences(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to get event occurrences",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, occurrences)
}

// GetOccurrencesInRange godoc
// @Summary Get all event occurrences in date range
// @Description Get all occurrences from all events within a specified date range
// @Tags events
// @Accept json
// @Produce json
// @Param startDate query string true "Start date (YYYY-MM-DD)"
// @Param endDate query string true "End date (YYYY-MM-DD)"
// @Param timezone query string false "Timezone for results"
// @Success 200 {array} dto.EventOccurrenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/occurrences [get]
func (c *EventController) GetOccurrencesInRange(ctx *gin.Context) {
	var req dto.GetEventOccurrencesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}
	
	occurrences, err := c.eventService.GetOccurrencesInRange(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get occurrences",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, occurrences)
}

// UpdateSingleOccurrence godoc
// @Summary Update a single occurrence of a recurring event
// @Description Update only one specific occurrence of a recurring event series
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param update body dto.UpdateOccurrenceRequest true "Single occurrence update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/occurrence [put]
func (c *EventController) UpdateSingleOccurrence(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.UpdateOccurrenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err = c.eventService.UpdateSingleOccurrence(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "event is not recurring" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update occurrence",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Single occurrence updated successfully",
	})
}

// UpdateFutureOccurrences godoc
// @Summary Update this and all future occurrences of a recurring event
// @Description Update the current occurrence and all future occurrences by creating a new series
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param update body dto.UpdateFutureOccurrencesRequest true "Future occurrences update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/future [put]
func (c *EventController) UpdateFutureOccurrences(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	var req dto.UpdateFutureOccurrencesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err = c.eventService.UpdateFutureOccurrences(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		} else if err.Error() == "event is not recurring" {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update future occurrences",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Future occurrences updated successfully",
	})
}

// ValidateRecurrenceRule godoc
// @Summary Validate a recurrence rule
// @Description Validate the format and logic of a recurrence rule
// @Tags events
// @Accept json
// @Produce json
// @Param rule body dto.CreateRecurrenceRuleRequest true "Recurrence rule to validate"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /events/validate-recurrence [post]
func (c *EventController) ValidateRecurrenceRule(ctx *gin.Context) {
	var req dto.CreateRecurrenceRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	
	err := c.eventService.ValidateRecurrenceRule(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid recurrence rule",
			"details": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recurrence rule is valid",
		"rule":    req,
	})
}

// GetNextOccurrence godoc
// @Summary Get next occurrence of an event after a specific date
// @Description Get the next occurrence of an event (recurring or single) after a given date
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param after query string true "Date after which to find next occurrence (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/next [get]
func (c *EventController) GetNextOccurrence(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}
	
	afterStr := ctx.Query("after")
	if afterStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'after' query parameter",
		})
		return
	}
	
	after, err := time.Parse("2006-01-02", afterStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format for 'after' parameter",
		})
		return
	}
	
	nextOccurrence, err := c.eventService.GetNextOccurrence(id, after)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "failed to get event: record not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to get next occurrence",
			"details": err.Error(),
		})
		return
	}
	
	if nextOccurrence == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message":        "No next occurrence found",
			"nextOccurrence": nil,
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"nextOccurrence": nextOccurrence.Format("2006-01-02T15:04:05Z07:00"),
	})
}