package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateKetersediaanRequest is the request to create availability
type CreateKetersediaanRequest struct {
	PersonID       uuid.UUID `json:"personId" validate:"required"`
	EventID        uuid.UUID `json:"eventId" validate:"required"`
	OccurrenceDate string    `json:"occurrenceDate" validate:"required"` // YYYY-MM-DD format
	Status         string    `json:"status" validate:"required,oneof=available unavailable tentative"`
	Notes          string    `json:"notes,omitempty"`
}

// UpdateKetersediaanRequest is the request to update availability
type UpdateKetersediaanRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=available unavailable tentative"`
	Notes  *string `json:"notes,omitempty"`
}

// BulkCreateKetersediaanRequest is the request to create multiple availabilities at once
type BulkCreateKetersediaanRequest struct {
	Availabilities []CreateKetersediaanRequest `json:"availabilities" validate:"required,min=1,dive"`
}

// KetersediaanResponse is the response for availability data
type KetersediaanResponse struct {
	ID             uuid.UUID            `json:"id"`
	PersonID       uuid.UUID            `json:"personId"`
	PersonName     string               `json:"personName"`
	EventID        uuid.UUID            `json:"eventId"`
	EventTitle     string               `json:"eventTitle"`
	OccurrenceDate time.Time            `json:"occurrenceDate"`
	Status         string               `json:"status"`
	Notes          string               `json:"notes,omitempty"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	Event          *EventResponse       `json:"event,omitempty"`
	Person         *PersonBasicResponse `json:"person,omitempty"`
}

// PersonBasicResponse is a simplified person response
type PersonBasicResponse struct {
	ID    uuid.UUID `json:"id"`
	Nama  string    `json:"nama"`
	Email string    `json:"email"`
}

// GetKetersediaanRequest is the request to get availability with filters
type GetKetersediaanRequest struct {
	PersonID  *uuid.UUID `form:"personId"`
	EventID   *uuid.UUID `form:"eventId"`
	StartDate string     `form:"startDate"` // YYYY-MM-DD format
	EndDate   string     `form:"endDate"`   // YYYY-MM-DD format
	Status    string     `form:"status,omitempty" validate:"omitempty,oneof=available unavailable tentative"`
	Page      int        `form:"page,omitempty"`
	Limit     int        `form:"limit,omitempty"`
}

// KetersediaanListResponse is the paginated list response
type KetersediaanListResponse struct {
	Availabilities []KetersediaanResponse `json:"availabilities"`
	TotalCount     int                    `json:"totalCount"`
	Page           int                    `json:"page"`
	Limit          int                    `json:"limit"`
}

// GetEventAvailabilitySummaryRequest gets summary of availability for an event occurrence
type GetEventAvailabilitySummaryRequest struct {
	EventID        uuid.UUID `form:"eventId" validate:"required"`
	OccurrenceDate string    `form:"occurrenceDate" validate:"required"` // YYYY-MM-DD format
}

// EventAvailabilitySummaryResponse provides statistics about availability for an event
type EventAvailabilitySummaryResponse struct {
	EventID        uuid.UUID `json:"eventId"`
	EventTitle     string    `json:"eventTitle"`
	OccurrenceDate time.Time `json:"occurrenceDate"`
	TotalResponses int       `json:"totalResponses"`
	Available      int       `json:"available"`
	Unavailable    int       `json:"unavailable"`
	Tentative      int       `json:"tentative"`
	Responses      []KetersediaanResponse `json:"responses,omitempty"`
}
