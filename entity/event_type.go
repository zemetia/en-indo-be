package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventType represents a type of event (e.g., "event", "ibadah", "spiritual journey")
type EventType struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"type:varchar(255);unique;not null"` // e.g., "event", "ibadah", "spiritual_journey" (auto-generated from DisplayName)
	DisplayName string    `gorm:"type:varchar(255);not null"`        // Display name for UI
	Description string    `gorm:"type:text"`                         // Description of this event type
	Color       string    `gorm:"type:varchar(7);default:'#2980b9'"` // Hex color for UI (e.g., "#2980b9")
	IsActive    bool      `gorm:"default:true;not null;index"`       // Active status
	SortOrder   int       `gorm:"default:0;not null"`                // Display order

	Timestamp
}

func (et *EventType) BeforeCreate(tx *gorm.DB) error {
	if et.ID == uuid.Nil {
		et.ID = uuid.New()
	}
	return nil
}
