package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/service"
)

type EventParticipantController struct {
	participantService service.EventParticipantService
}

func NewEventParticipantController(participantService service.EventParticipantService) *EventParticipantController {
	return &EventParticipantController{
		participantService: participantService,
	}
}

// RegisterParticipant godoc
// @Summary Register a participant for an event
// @Description Register a person or visitor for a specific event occurrence
// @Tags event-participants
// @Accept json
// @Produce json
// @Param participant body dto.CreateEventParticipantRequest true "Participant registration data"
// @Success 201 {object} dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants [post]
func (c *EventParticipantController) RegisterParticipant(ctx *gin.Context) {
	var req dto.CreateEventParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	participant, err := c.participantService.RegisterParticipant(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "participant already registered for this event occurrence" {
			statusCode = http.StatusConflict
		}
		ctx.JSON(statusCode, gin.H{
			"error":   "Failed to register participant",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, participant)
}

// BulkRegisterParticipants godoc
// @Summary Bulk register participants for an event
// @Description Register multiple persons and/or visitors for a specific event occurrence
// @Tags event-participants
// @Accept json
// @Produce json
// @Param participants body dto.BulkRegisterParticipantsRequest true "Bulk registration data"
// @Success 201 {object} []dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants/bulk [post]
func (c *EventParticipantController) BulkRegisterParticipants(ctx *gin.Context) {
	var req dto.BulkRegisterParticipantsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	participants, err := c.participantService.BulkRegisterParticipants(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to bulk register participants",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, participants)
}

// UpdateParticipant godoc
// @Summary Update a participant
// @Description Update participant registration or attendance status
// @Tags event-participants
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Param participant body dto.UpdateEventParticipantRequest true "Update data"
// @Success 200 {object} dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants/{id} [put]
func (c *EventParticipantController) UpdateParticipant(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid participant ID format",
		})
		return
	}

	var req dto.UpdateEventParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	participant, err := c.participantService.UpdateParticipant(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update participant",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participant)
}

// DeleteParticipant godoc
// @Summary Delete a participant
// @Description Remove a participant from an event
// @Tags event-participants
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants/{id} [delete]
func (c *EventParticipantController) DeleteParticipant(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid participant ID format",
		})
		return
	}

	if err := c.participantService.DeleteParticipant(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete participant",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetParticipant godoc
// @Summary Get participant by ID
// @Description Get a single participant's details
// @Tags event-participants
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Success 200 {object} dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/events/participants/{id} [get]
func (c *EventParticipantController) GetParticipant(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid participant ID format",
		})
		return
	}

	participant, err := c.participantService.GetParticipant(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Participant not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participant)
}

// GetEventParticipants godoc
// @Summary Get all participants for an event occurrence
// @Description Get list of participants for a specific event occurrence
// @Tags event-participants
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param date query string true "Occurrence date (YYYY-MM-DD)"
// @Param participantType query string false "Filter by participant type (person/visitor)"
// @Param registrationStatus query string false "Filter by registration status"
// @Param attendanceStatus query string false "Filter by attendance status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.EventParticipantListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{eventId}/participants [get]
func (c *EventParticipantController) GetEventParticipants(ctx *gin.Context) {
	eventIDStr := ctx.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Occurrence date is required",
		})
		return
	}

	occurrenceDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format (use YYYY-MM-DD)",
		})
		return
	}

	// Build filter request
	filterReq := &dto.EventParticipantFilterRequest{
		EventID:        &eventID,
		OccurrenceDate: &dateStr,
	}

	if participantType := ctx.Query("participantType"); participantType != "" {
		filterReq.ParticipantType = &participantType
	}
	if regStatus := ctx.Query("registrationStatus"); regStatus != "" {
		filterReq.RegistrationStatus = &regStatus
	}
	if attStatus := ctx.Query("attendanceStatus"); attStatus != "" {
		filterReq.AttendanceStatus = &attStatus
	}

	// Parse pagination
	var page, limit int
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filterReq.Page = page
	filterReq.Limit = limit

	participants, err := c.participantService.GetEventParticipants(eventID, occurrenceDate, filterReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get participants",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participants)
}

// CheckInParticipant godoc
// @Summary Check in a participant
// @Description Mark a participant as present
// @Tags event-participants
// @Accept json
// @Produce json
// @Param id path string true "Participant ID"
// @Param checkIn body dto.CheckInParticipantRequest true "Check-in data"
// @Success 200 {object} dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants/{id}/checkin [post]
func (c *EventParticipantController) CheckInParticipant(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid participant ID format",
		})
		return
	}

	var req dto.CheckInParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	participant, err := c.participantService.CheckInParticipant(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check in participant",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participant)
}

// BulkCheckIn godoc
// @Summary Bulk check in participants
// @Description Mark multiple participants as present
// @Tags event-participants
// @Accept json
// @Produce json
// @Param checkIns body dto.BulkCheckInRequest true "Bulk check-in data"
// @Success 200 {object} []dto.EventParticipantResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/participants/checkin/bulk [post]
func (c *EventParticipantController) BulkCheckIn(ctx *gin.Context) {
	var req dto.BulkCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	participants, err := c.participantService.BulkCheckIn(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to bulk check in participants",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participants)
}

// GetAttendanceReport godoc
// @Summary Get attendance report
// @Description Get detailed attendance report for an event occurrence
// @Tags event-participants
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param date query string true "Occurrence date (YYYY-MM-DD)"
// @Success 200 {object} dto.AttendanceReportResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/events/{eventId}/attendance/report [get]
func (c *EventParticipantController) GetAttendanceReport(ctx *gin.Context) {
	eventIDStr := ctx.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event ID format",
		})
		return
	}

	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Occurrence date is required",
		})
		return
	}

	occurrenceDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format (use YYYY-MM-DD)",
		})
		return
	}

	report, err := c.participantService.GetAttendanceReport(eventID, occurrenceDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate attendance report",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// ListParticipants godoc
// @Summary List all participants
// @Description Get list of all participants with filters
// @Tags event-participants
// @Accept json
// @Produce json
// @Param eventId query string false "Filter by event ID"
// @Param date query string false "Filter by occurrence date"
// @Param participantType query string false "Filter by participant type"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.EventParticipantListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/participants [get]
func (c *EventParticipantController) ListParticipants(ctx *gin.Context) {
	filterReq := &dto.EventParticipantFilterRequest{}

	if eventIDStr := ctx.Query("eventId"); eventIDStr != "" {
		eventID, err := uuid.Parse(eventIDStr)
		if err == nil {
			filterReq.EventID = &eventID
		}
	}

	if dateStr := ctx.Query("date"); dateStr != "" {
		filterReq.OccurrenceDate = &dateStr
	}

	if participantType := ctx.Query("participantType"); participantType != "" {
		filterReq.ParticipantType = &participantType
	}

	// Parse pagination
	var page, limit int
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filterReq.Page = page
	filterReq.Limit = limit

	participants, err := c.participantService.ListParticipants(filterReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list participants",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, participants)
}

// ScanQRCheckIn godoc
// @Summary Check in via QR scan
// @Description Check in a participant by scanning QR code (uses person_id from JWT token)
// @Tags event-participants
// @Accept json
// @Produce json
// @Param request body dto.QRScanCheckInRequest true "QR scan check-in request"
// @Success 200 {object} dto.QRScanCheckInResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/event-participants/scan-qr [post]
// @Security BearerAuth
func (c *EventParticipantController) ScanQRCheckIn(ctx *gin.Context) {
	// Get person_id from JWT token (set by authentication middleware)
	personIDInterface, exists := ctx.Get("person_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Person ID not found in token",
		})
		return
	}

	personIDStr, ok := personIDInterface.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid person ID format in token",
		})
		return
	}

	personID, err := uuid.Parse(personIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid person ID format",
		})
		return
	}

	// Parse request body
	var req dto.QRScanCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call service to perform check-in
	response, err := c.participantService.ScanQRCheckIn(personID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process QR scan check-in",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
