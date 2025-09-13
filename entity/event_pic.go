package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventPIC represents a Person in Charge (Penanggung Jawab) for an event
type EventPIC struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	EventID     uuid.UUID `gorm:"type:char(36);not null;index"`
	Event       Event     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:EventID"`
	PersonID    uuid.UUID `gorm:"type:char(36);not null;index"`
	Person      Person    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PersonID"`
	Role        string    `gorm:"type:varchar(100);not null"`        // e.g., "Primary PIC", "Co-PIC", "Technical PIC"
	Description string    `gorm:"type:text"`                         // Role description
	IsActive    bool      `gorm:"default:true;not null"`             // Active status
	IsPrimary   bool      `gorm:"default:false;not null;index"`      // Primary PIC flag
	StartDate   time.Time `gorm:"type:date;not null"`                // When this PIC assignment starts
	EndDate     *time.Time `gorm:"type:date"`                        // When this PIC assignment ends (null for ongoing)
	
	// Responsibilities and permissions
	CanEdit     bool `gorm:"default:false;not null"` // Can edit event details
	CanDelete   bool `gorm:"default:false;not null"` // Can delete event
	CanAssignPIC bool `gorm:"default:false;not null"` // Can assign other PICs
	
	// Contact preferences for this specific event
	NotifyOnChanges bool `gorm:"default:true;not null"`  // Notify when event changes
	NotifyOnReminders bool `gorm:"default:true;not null"` // Notify for event reminders

	Timestamp
}

// EventPICRole represents predefined roles for Event PICs
type EventPICRole struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"type:varchar(100);unique;not null"` // e.g., "Event Coordinator", "Technical Lead"
	Description string    `gorm:"type:text"`
	
	// Default permissions for this role
	DefaultCanEdit      bool `gorm:"default:false;not null"`
	DefaultCanDelete    bool `gorm:"default:false;not null"`
	DefaultCanAssignPIC bool `gorm:"default:false;not null"`
	
	IsActive bool `gorm:"default:true;not null"`
	
	Timestamp
}

// EventPICHistory tracks changes in PIC assignments for audit purposes
type EventPICHistory struct {
	ID          uuid.UUID  `gorm:"type:char(36);primary_key"`
	EventID     uuid.UUID  `gorm:"type:char(36);not null;index"`
	Event       Event      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:EventID"`
	PersonID    uuid.UUID  `gorm:"type:char(36);not null"`
	Person      Person     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PersonID"`
	Action      string     `gorm:"type:varchar(50);not null"` // "assigned", "removed", "role_changed", "permissions_changed"
	OldRole     string     `gorm:"type:varchar(100)"`
	NewRole     string     `gorm:"type:varchar(100)"`
	ChangedBy   uuid.UUID  `gorm:"type:char(36);not null"`
	ChangedByPerson Person `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ChangedBy"`
	Reason      string     `gorm:"type:text"`
	ActionDate  time.Time  `gorm:"type:timestamp;not null"`

	Timestamp
}

func (ep *EventPIC) BeforeCreate(tx *gorm.DB) error {
	if ep.ID == uuid.Nil {
		ep.ID = uuid.New()
	}
	return nil
}

func (epr *EventPICRole) BeforeCreate(tx *gorm.DB) error {
	if epr.ID == uuid.Nil {
		epr.ID = uuid.New()
	}
	return nil
}

func (eph *EventPICHistory) BeforeCreate(tx *gorm.DB) error {
	if eph.ID == uuid.Nil {
		eph.ID = uuid.New()
	}
	if eph.ActionDate.IsZero() {
		eph.ActionDate = time.Now()
	}
	return nil
}

// Business logic constants
const (
	EventPICActionAssigned           = "assigned"
	EventPICActionRemoved            = "removed"
	EventPICActionRoleChanged        = "role_changed"
	EventPICActionPermissionsChanged = "permissions_changed"
	
	// Default roles
	EventPICRolePrimary    = "Primary PIC"
	EventPICRoleSecondary  = "Secondary PIC"
	EventPICRoleTechnical  = "Technical PIC"
	EventPICRoleLogistics  = "Logistics PIC"
	EventPICRoleRegistration = "Registration PIC"
)