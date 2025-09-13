package dto

import (
	"time"

	"github.com/google/uuid"
)

// EventPIC creation and management DTOs

type CreateEventPICRequest struct {
	PersonID    uuid.UUID `json:"personId" validate:"required"`
	Role        string    `json:"role" validate:"required,min=1,max=100"`
	Description string    `json:"description,omitempty"`
	IsPrimary   bool      `json:"isPrimary"`
	StartDate   string    `json:"startDate" validate:"required"` // YYYY-MM-DD format
	EndDate     *string   `json:"endDate,omitempty"`             // YYYY-MM-DD format, nullable
	
	// Permissions
	CanEdit      bool `json:"canEdit"`
	CanDelete    bool `json:"canDelete"`
	CanAssignPIC bool `json:"canAssignPIC"`
	
	// Notifications
	NotifyOnChanges   bool `json:"notifyOnChanges"`
	NotifyOnReminders bool `json:"notifyOnReminders"`
}

type UpdateEventPICRequest struct {
	Role        *string `json:"role,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	IsPrimary   *bool   `json:"isPrimary,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
	StartDate   *string `json:"startDate,omitempty"` // YYYY-MM-DD format
	EndDate     *string `json:"endDate,omitempty"`   // YYYY-MM-DD format, nullable
	
	// Permissions
	CanEdit      *bool `json:"canEdit,omitempty"`
	CanDelete    *bool `json:"canDelete,omitempty"`
	CanAssignPIC *bool `json:"canAssignPIC,omitempty"`
	
	// Notifications  
	NotifyOnChanges   *bool `json:"notifyOnChanges,omitempty"`
	NotifyOnReminders *bool `json:"notifyOnReminders,omitempty"`
}

type EventPICResponse struct {
	ID          uuid.UUID     `json:"id"`
	EventID     uuid.UUID     `json:"eventId"`
	PersonID    uuid.UUID     `json:"personId"`
	Person      PersonSummary `json:"person"`
	Role        string        `json:"role"`
	Description string        `json:"description"`
	IsActive    bool          `json:"isActive"`
	IsPrimary   bool          `json:"isPrimary"`
	StartDate   time.Time     `json:"startDate"`
	EndDate     *time.Time    `json:"endDate"`
	
	// Permissions
	CanEdit      bool `json:"canEdit"`
	CanDelete    bool `json:"canDelete"`
	CanAssignPIC bool `json:"canAssignPIC"`
	
	// Notifications
	NotifyOnChanges   bool `json:"notifyOnChanges"`
	NotifyOnReminders bool `json:"notifyOnReminders"`
	
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PersonSummary struct {
	ID           uuid.UUID `json:"id"`
	Nama         string    `json:"nama"`
	Email        string    `json:"email"`
	NomorTelepon string    `json:"nomorTelepon"`
	ChurchID     uuid.UUID `json:"churchId"`
}

// Bulk operations
type BulkAssignEventPICRequest struct {
	PICs []CreateEventPICRequest `json:"pics" validate:"required,min=1,dive"`
}

type BulkUpdateEventPICRequest struct {
	Updates []struct {
		PICID   uuid.UUID             `json:"picId" validate:"required"`
		Updates UpdateEventPICRequest `json:"updates" validate:"required"`
	} `json:"updates" validate:"required,min=1,dive"`
}

// Event PIC Role management
type CreateEventPICRoleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	
	// Default permissions
	DefaultCanEdit      bool `json:"defaultCanEdit"`
	DefaultCanDelete    bool `json:"defaultCanDelete"`
	DefaultCanAssignPIC bool `json:"defaultCanAssignPIC"`
}

type UpdateEventPICRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
	
	// Default permissions
	DefaultCanEdit      *bool `json:"defaultCanEdit,omitempty"`
	DefaultCanDelete    *bool `json:"defaultCanDelete,omitempty"`
	DefaultCanAssignPIC *bool `json:"defaultCanAssignPIC,omitempty"`
}

type EventPICRoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	
	// Default permissions
	DefaultCanEdit      bool `json:"defaultCanEdit"`
	DefaultCanDelete    bool `json:"defaultCanDelete"`
	DefaultCanAssignPIC bool `json:"defaultCanAssignPIC"`
	
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// History and audit
type EventPICHistoryResponse struct {
	ID          uuid.UUID     `json:"id"`
	EventID     uuid.UUID     `json:"eventId"`
	PersonID    uuid.UUID     `json:"personId"`
	Person      PersonSummary `json:"person"`
	Action      string        `json:"action"`
	OldRole     string        `json:"oldRole"`
	NewRole     string        `json:"newRole"`
	ChangedBy   uuid.UUID     `json:"changedBy"`
	ChangedByPerson PersonSummary `json:"changedByPerson"`
	Reason      string        `json:"reason"`
	ActionDate  time.Time     `json:"actionDate"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// Query filters
type EventPICFilterRequest struct {
	EventID   *uuid.UUID `form:"eventId,omitempty"`
	PersonID  *uuid.UUID `form:"personId,omitempty"`
	Role      string     `form:"role,omitempty"`
	IsActive  *bool      `form:"isActive,omitempty"`
	IsPrimary *bool      `form:"isPrimary,omitempty"`
	Search    string     `form:"search,omitempty"` // Search in person name, role
	Page      int        `form:"page,omitempty"`
	Limit     int        `form:"limit,omitempty"`
}

// Response wrappers
type EventPICListResponse struct {
	PICs       []EventPICResponse `json:"pics"`
	TotalCount int                `json:"totalCount"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

type EventPICRoleListResponse struct {
	Roles      []EventPICRoleResponse `json:"roles"`
	TotalCount int                    `json:"totalCount"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
}

type EventPICHistoryListResponse struct {
	History    []EventPICHistoryResponse `json:"history"`
	TotalCount int                       `json:"totalCount"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
}

// Transfer PIC ownership
type TransferEventPICRequest struct {
	FromPersonID uuid.UUID `json:"fromPersonId" validate:"required"`
	ToPersonID   uuid.UUID `json:"toPersonId" validate:"required"`
	TransferType string    `json:"transferType" validate:"required,oneof=replace add_as_secondary"`
	Reason       string    `json:"reason,omitempty"`
	EffectiveDate string   `json:"effectiveDate" validate:"required"` // YYYY-MM-DD format
}