package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type KetersediaanController struct {
	ketersediaanService service.KetersediaanService
}

func NewKetersediaanController(ketersediaanService service.KetersediaanService) *KetersediaanController {
	return &KetersediaanController{
		ketersediaanService: ketersediaanService,
	}
}

// CreateKetersediaan godoc
// @Summary Create availability for an event occurrence
// @Description Create a person's availability status for a specific event occurrence
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param availability body dto.CreateKetersediaanRequest true "Availability data"
// @Success 201 {object} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan [post]
func (c *KetersediaanController) CreateKetersediaan(ctx *gin.Context) {
	var req dto.CreateKetersediaanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.CreateKetersediaan(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// UpsertKetersediaan godoc
// @Summary Create or update availability
// @Description Create or update (if exists) a person's availability status for a specific event occurrence
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param availability body dto.CreateKetersediaanRequest true "Availability data"
// @Success 200 {object} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan/upsert [post]
func (c *KetersediaanController) UpsertKetersediaan(ctx *gin.Context) {
	var req dto.CreateKetersediaanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.UpsertKetersediaan(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to upsert availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetKetersediaan godoc
// @Summary Get availability by ID
// @Description Get a specific availability record by ID
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Success 200 {object} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /ketersediaan/{id} [get]
func (c *KetersediaanController) GetKetersediaan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID format",
		})
		return
	}

	result, err := c.ketersediaanService.GetKetersediaan(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Availability not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// UpdateKetersediaan godoc
// @Summary Update availability
// @Description Update an existing availability record
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Param availability body dto.UpdateKetersediaanRequest true "Update data"
// @Success 200 {object} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan/{id} [put]
func (c *KetersediaanController) UpdateKetersediaan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID format",
		})
		return
	}

	var req dto.UpdateKetersediaanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.UpdateKetersediaan(id, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "availability not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to update availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteKetersediaan godoc
// @Summary Delete availability
// @Description Delete an availability record
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param id path string true "Availability ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan/{id} [delete]
func (c *KetersediaanController) DeleteKetersediaan(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid availability ID format",
		})
		return
	}

	err = c.ketersediaanService.DeleteKetersediaan(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "availability not found" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"error":   "Failed to delete availability",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetPersonAvailability godoc
// @Summary Get person's availability
// @Description Get availability records for a specific person with filters
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param personId query string true "Person ID"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Status filter (available, unavailable, tentative)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 50)"
// @Success 200 {object} dto.KetersediaanListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan [get]
func (c *KetersediaanController) GetPersonAvailability(ctx *gin.Context) {
	var req dto.GetKetersediaanRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.GetPersonAvailability(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetEventAvailability godoc
// @Summary Get event availability
// @Description Get all availability responses for a specific event occurrence
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param occurrenceDate query string false "Occurrence date (YYYY-MM-DD)"
// @Success 200 {array} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id}/availability [get]
func (c *KetersediaanController) GetEventAvailability(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	var occurrenceDate *time.Time
	occurrenceDateStr := ctx.Query("occurrenceDate")
	if occurrenceDateStr != "" {
		parsed, err := time.Parse("2006-01-02", occurrenceDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid occurrence date format",
			})
			return
		}
		occurrenceDate = &parsed
	}

	result, err := c.ketersediaanService.GetEventAvailability(eventID, occurrenceDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get event availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// BulkCreateKetersediaan godoc
// @Summary Bulk create availability
// @Description Create multiple availability records at once
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param availabilities body dto.BulkCreateKetersediaanRequest true "Bulk availability data"
// @Success 201 {array} dto.KetersediaanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan/bulk [post]
func (c *KetersediaanController) BulkCreateKetersediaan(ctx *gin.Context) {
	var req dto.BulkCreateKetersediaanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.BulkCreateKetersediaan(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to bulk create availability",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// GetEventAvailabilitySummary godoc
// @Summary Get event availability summary
// @Description Get summary statistics of availability for an event occurrence
// @Tags ketersediaan
// @Accept json
// @Produce json
// @Param eventId query string true "Event ID"
// @Param occurrenceDate query string true "Occurrence date (YYYY-MM-DD)"
// @Success 200 {object} dto.EventAvailabilitySummaryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ketersediaan/summary [get]
func (c *KetersediaanController) GetEventAvailabilitySummary(ctx *gin.Context) {
	var req dto.GetEventAvailabilitySummaryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := c.ketersediaanService.GetEventAvailabilitySummary(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get availability summary",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
