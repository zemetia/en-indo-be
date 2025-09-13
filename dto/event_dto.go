package dto

import (
	"time"

	"github.com/google/uuid"
)

// Event creation request
type CreateEventRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	BannerImage string `json:"bannerImage,omitempty"`
	Description string `json:"description,omitempty"`
	Capacity    int    `json:"capacity,omitempty"`
	Type        string `json:"type" validate:"required,oneof=event ibadah spiritual_journey"`

	EventDate     string `json:"eventDate" validate:"required"`
	EventLocation string `json:"eventLocation" validate:"required,max=255"`
	StartTime     string `json:"startTime" validate:"required"`
	EndTime       string `json:"endTime" validate:"required"`
	AllDay        bool   `json:"allDay"`
	Timezone      string `json:"timezone" validate:"required"`

	IsPublic              bool       `json:"isPublic"`
	DiscipleshipJourneyID *uuid.UUID `json:"discipleshipJourneyId,omitempty"`

	// Expected participant counts for planning
	ExpectedParticipants *int `json:"expectedParticipants,omitempty"`
	ExpectedAdults      *int `json:"expectedAdults,omitempty"`
	ExpectedYouth       *int `json:"expectedYouth,omitempty"`
	ExpectedKids        *int `json:"expectedKids,omitempty"`

	RecurrenceRule *CreateRecurrenceRuleRequest `json:"recurrenceRule,omitempty"`
	LaguIDs        []uuid.UUID                  `json:"laguIds,omitempty"`
	
	// Event PIC assignments during creation
	EventPICs      []CreateEventPICRequest      `json:"eventPics,omitempty"`
}

type CreateRecurrenceRuleRequest struct {
	Frequency  string   `json:"frequency" validate:"required,oneof=DAILY WEEKLY MONTHLY YEARLY"`
	Interval   int      `json:"interval,omitempty"`
	ByWeekday  []string `json:"byWeekday,omitempty"`
	ByMonthDay []int64  `json:"byMonthDay,omitempty"`
	ByMonth    []int64  `json:"byMonth,omitempty"`
	BySetPos   []int64  `json:"bySetPos,omitempty"`  // e.g. [1,-1] for first/last occurrence
	WeekStart  string   `json:"weekStart,omitempty"` // week start day (MO, SU, etc.)
	ByYearDay  []int64  `json:"byYearDay,omitempty"` // day of year (1-366)
	Count      *int     `json:"count,omitempty"`
	Until      *string  `json:"until,omitempty"`
}

// Event update request
type UpdateEventRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	BannerImage *string `json:"bannerImage,omitempty"`
	Description *string `json:"description,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=event ibadah spiritual_journey"`

	EventDate     *string `json:"eventDate,omitempty"`
	EventLocation *string `json:"eventLocation,omitempty" validate:"omitempty,max=255"`
	StartTime     *string `json:"startTime,omitempty"`
	EndTime       *string `json:"endTime,omitempty"`
	AllDay        *bool   `json:"allDay,omitempty"`
	Timezone      *string `json:"timezone,omitempty"`

	IsPublic              *bool      `json:"isPublic,omitempty"`
	DiscipleshipJourneyID *uuid.UUID `json:"discipleshipJourneyId,omitempty"`

	// Expected participant counts for planning
	ExpectedParticipants *int `json:"expectedParticipants,omitempty"`
	ExpectedAdults      *int `json:"expectedAdults,omitempty"`
	ExpectedYouth       *int `json:"expectedYouth,omitempty"`
	ExpectedKids        *int `json:"expectedKids,omitempty"`

	LaguIDs *[]uuid.UUID `json:"laguIds,omitempty"`
}

// Update type for recurring events
type RecurringUpdateType string

const (
	UpdateThisEvent    RecurringUpdateType = "single"
	UpdateAllEvents    RecurringUpdateType = "all"
	UpdateFutureEvents RecurringUpdateType = "future"
)

type UpdateRecurringEventRequest struct {
	UpdateType     RecurringUpdateType `json:"updateType" validate:"required,oneof=single all future"`
	OccurrenceDate string              `json:"occurrenceDate,omitempty"`
	StartTime      *string             `json:"startTime,omitempty"` // New start time for this modification
	EndTime        *string             `json:"endTime,omitempty"`   // New end time for this modification
	Event          UpdateEventRequest  `json:"event"`
}

// Event occurrence for individual recurring event instance
type EventOccurrenceResponse struct {
	EventID        uuid.UUID      `json:"eventId"`
	OccurrenceDate time.Time      `json:"occurrenceDate"`
	StartDatetime  time.Time      `json:"startDatetime"`
	EndDatetime    time.Time      `json:"endDatetime"`
	IsException    bool           `json:"isException"`
	IsSkipped      bool           `json:"isSkipped"`
	ExceptionNotes string         `json:"exceptionNotes,omitempty"`
	OriginalEvent  *EventResponse `json:"originalEvent,omitempty"`
}

// Main event response
type EventResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	BannerImage string    `json:"bannerImage,omitempty"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	Type        string    `json:"type"`

	EventDate     time.Time `json:"eventDate"`
	EventLocation string    `json:"eventLocation"`
	StartDatetime time.Time `json:"startDatetime"`
	EndDatetime   time.Time `json:"endDatetime"`
	AllDay        bool      `json:"allDay"`
	Timezone      string    `json:"timezone"`

	IsPublic              bool       `json:"isPublic"`
	DiscipleshipJourneyID *uuid.UUID `json:"discipleshipJourneyId,omitempty"`

	// Expected participant counts for planning
	ExpectedParticipants int `json:"expectedParticipants"`
	ExpectedAdults      int `json:"expectedAdults"`
	ExpectedYouth       int `json:"expectedYouth"`
	ExpectedKids        int `json:"expectedKids"`

	RecurrenceRule *RecurrenceRuleResponse `json:"recurrenceRule,omitempty"`
	Lagu           []LaguResponse          `json:"lagu,omitempty"`
	
	// Event PIC information
	EventPICs      []EventPICResponse      `json:"eventPics,omitempty"`
	PrimaryPIC     *EventPICResponse       `json:"primaryPic,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RecurrenceRuleResponse struct {
	ID         uuid.UUID  `json:"id"`
	Frequency  string     `json:"frequency"`
	Interval   int        `json:"interval"`
	ByWeekday  []string   `json:"byWeekday"`
	ByMonthDay []int64    `json:"byMonthDay"`
	ByMonth    []int64    `json:"byMonth"`
	BySetPos   []int64    `json:"bySetPos"`
	WeekStart  string     `json:"weekStart"`
	ByYearDay  []int64    `json:"byYearDay"`
	Count      *int       `json:"count"`
	Until      *time.Time `json:"until"`
}

type LaguResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

// Request for getting event occurrences
type GetEventOccurrencesRequest struct {
	StartDate string `form:"startDate" validate:"required"`
	EndDate   string `form:"endDate" validate:"required"`
	Timezone  string `form:"timezone,omitempty"`
}

// Request for updating single occurrence
type UpdateOccurrenceRequest struct {
	OccurrenceDate string             `json:"occurrenceDate" validate:"required"`
	StartTime      *string            `json:"startTime,omitempty"`
	EndTime        *string            `json:"endTime,omitempty"`
	Event          UpdateEventRequest `json:"event"`
}

// Request for updating future occurrences
type UpdateFutureOccurrencesRequest struct {
	FromDate       string                       `json:"fromDate" validate:"required"`
	StartTime      *string                      `json:"startTime,omitempty"`
	EndTime        *string                      `json:"endTime,omitempty"`
	Event          UpdateEventRequest           `json:"event"`
	RecurrenceRule *CreateRecurrenceRuleRequest `json:"recurrenceRule,omitempty"`
}

// Request for deleting specific occurrence
type DeleteOccurrenceRequest struct {
	OccurrenceDate string `json:"occurrenceDate" validate:"required"`
	DeleteType     string `json:"deleteType" validate:"required,oneof=single future"`
}

// Batch event response
type EventListResponse struct {
	Events     []EventResponse `json:"events"`
	TotalCount int             `json:"totalCount"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
}

// Event filter request
type EventFilterRequest struct {
	Type      string `form:"type,omitempty"`
	IsPublic  *bool  `form:"isPublic,omitempty"`
	StartDate string `form:"startDate,omitempty"`
	EndDate   string `form:"endDate,omitempty"`
	Search    string `form:"search,omitempty"`
	Page      int    `form:"page,omitempty"`
	Limit     int    `form:"limit,omitempty"`
	Timezone  string `form:"timezone,omitempty"`
}
