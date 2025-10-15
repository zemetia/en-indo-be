package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateEventParticipantRequest represents a request to register a participant for an event
type CreateEventParticipantRequest struct {
	EventID         uuid.UUID `json:"eventId" validate:"required"`
	OccurrenceDate  string    `json:"occurrenceDate" validate:"required"` // Format: 2006-01-02
	ParticipantType string    `json:"participantType" validate:"required,oneof=person visitor"`
	ParticipantID   uuid.UUID `json:"participantId" validate:"required"`
	Notes           string    `json:"notes,omitempty"`
}

// BulkRegisterParticipantsRequest represents a request to register multiple participants
type BulkRegisterParticipantsRequest struct {
	EventID        uuid.UUID   `json:"eventId" validate:"required"`
	OccurrenceDate string      `json:"occurrenceDate" validate:"required"`
	PersonIDs      []uuid.UUID `json:"personIds,omitempty"`
	VisitorIDs     []uuid.UUID `json:"visitorIds,omitempty"`
}

// UpdateEventParticipantRequest represents a request to update a participant
type UpdateEventParticipantRequest struct {
	RegistrationStatus *string `json:"registrationStatus,omitempty" validate:"omitempty,oneof=registered attended cancelled no_show"`
	AttendanceStatus   *string `json:"attendanceStatus,omitempty" validate:"omitempty,oneof=present absent excused"`
	Notes              *string `json:"notes,omitempty"`
}

// CheckInParticipantRequest represents a request to check in a participant
type CheckInParticipantRequest struct {
	CheckInMethod   string     `json:"checkInMethod" validate:"required,oneof=qr_scan manual nfc self_checkin"`
	RecordedBy      *uuid.UUID `json:"recordedBy,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	QRScanTimestamp *time.Time `json:"qrScanTimestamp,omitempty"` // Timestamp from QR code
}

// BulkCheckInRequest represents a request to check in multiple participants
type BulkCheckInRequest struct {
	EventID        uuid.UUID   `json:"eventId" validate:"required"`
	OccurrenceDate string      `json:"occurrenceDate" validate:"required"`
	ParticipantIDs []uuid.UUID `json:"participantIds" validate:"required,min=1"`
	CheckInMethod  string      `json:"checkInMethod" validate:"required,oneof=qr_scan manual nfc self_checkin"`
	RecordedBy     *uuid.UUID  `json:"recordedBy,omitempty"`
}

// EventParticipantResponse represents a participant's information
type EventParticipantResponse struct {
	ID                 uuid.UUID              `json:"id"`
	EventID            uuid.UUID              `json:"eventId"`
	EventTitle         string                 `json:"eventTitle,omitempty"`
	OccurrenceDate     time.Time              `json:"occurrenceDate"`
	ParticipantType    string                 `json:"participantType"`
	ParticipantID      uuid.UUID              `json:"participantId"`
	ParticipantDetails ParticipantDetails     `json:"participantDetails"`
	RegistrationStatus string                 `json:"registrationStatus"`
	AttendanceStatus   string                 `json:"attendanceStatus"`
	CheckInTime        *time.Time             `json:"checkInTime,omitempty"`
	CheckOutTime       *time.Time             `json:"checkOutTime,omitempty"`
	CheckInMethod      *string                `json:"checkInMethod,omitempty"`
	QRScanTimestamp    *time.Time             `json:"qrScanTimestamp,omitempty"`
	Notes              string                 `json:"notes,omitempty"`
	RecordedBy         *uuid.UUID             `json:"recordedBy,omitempty"`
	RecordedByDetails  *PersonSummary         `json:"recordedByDetails,omitempty"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
}

// ParticipantDetails contains either Person or Visitor details
type ParticipantDetails struct {
	Type    string         `json:"type"` // "person" or "visitor"
	Person  *PersonSummary `json:"person,omitempty"`
	Visitor *VisitorSummary `json:"visitor,omitempty"`
}

// VisitorSummary contains summary information about a visitor
type VisitorSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber *string   `json:"phoneNumber,omitempty"`
	IGUsername  *string   `json:"igUsername,omitempty"`
}

// EventParticipantListResponse represents a list of participants
type EventParticipantListResponse struct {
	Participants []EventParticipantResponse `json:"participants"`
	TotalCount   int                        `json:"totalCount"`
	Page         int                        `json:"page"`
	Limit        int                        `json:"limit"`
	Stats        AttendanceStats            `json:"stats"`
}

// AttendanceStats represents attendance statistics
type AttendanceStats struct {
	TotalRegistered int `json:"totalRegistered"`
	TotalPresent    int `json:"totalPresent"`
	TotalAbsent     int `json:"totalAbsent"`
	TotalExcused    int `json:"totalExcused"`
	TotalCancelled  int `json:"totalCancelled"`
	TotalNoShow     int `json:"totalNoShow"`
	AttendanceRate  float64 `json:"attendanceRate"` // Percentage
}

// EventParticipantFilterRequest represents filters for querying participants
type EventParticipantFilterRequest struct {
	EventID            *uuid.UUID `json:"eventId,omitempty"`
	OccurrenceDate     *string    `json:"occurrenceDate,omitempty"`
	ParticipantType    *string    `json:"participantType,omitempty" validate:"omitempty,oneof=person visitor"`
	RegistrationStatus *string    `json:"registrationStatus,omitempty" validate:"omitempty,oneof=registered attended cancelled no_show"`
	AttendanceStatus   *string    `json:"attendanceStatus,omitempty" validate:"omitempty,oneof=present absent excused"`
	Page               int        `json:"page,omitempty"`
	Limit              int        `json:"limit,omitempty"`
}

// AttendanceReportResponse represents a detailed attendance report
type AttendanceReportResponse struct {
	EventID        uuid.UUID                  `json:"eventId"`
	EventTitle     string                     `json:"eventTitle"`
	OccurrenceDate time.Time                  `json:"occurrenceDate"`
	Stats          AttendanceStats            `json:"stats"`
	Participants   []EventParticipantResponse `json:"participants"`
	GeneratedAt    time.Time                  `json:"generatedAt"`
}

// QRScanCheckInRequest represents a request to check in via QR scan
// The person_id is extracted from JWT token, not included in request
type QRScanCheckInRequest struct {
	EventID        uuid.UUID `json:"eventId" validate:"required"`
	OccurrenceDate string    `json:"occurrenceDate" validate:"required"` // Format: 2006-01-02
}

// QRScanCheckInResponse represents the response from QR scan check-in
type QRScanCheckInResponse struct {
	Status      string                    `json:"status"` // not_participant | success | already_checked_in
	Message     string                    `json:"message"`
	EventTitle  string                    `json:"eventTitle,omitempty"`
	CheckInTime *time.Time                `json:"checkInTime,omitempty"`
	Participant *EventParticipantResponse `json:"participant,omitempty"`
}
