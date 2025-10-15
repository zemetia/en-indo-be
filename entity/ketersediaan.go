package entity

import (
	"time"

	"github.com/google/uuid"
)

// KetersediaanStatus represents availability status
type KetersediaanStatus string

const (
	KetersediaanAvailable   KetersediaanStatus = "available"
	KetersediaanUnavailable KetersediaanStatus = "unavailable"
	KetersediaanTentative   KetersediaanStatus = "tentative"
)

// Ketersediaan represents a person's availability for a specific event occurrence
type Ketersediaan struct {
	ID             uuid.UUID          `gorm:"type:char(36);primary_key"`
	PersonID       uuid.UUID          `gorm:"type:char(36);not null;index:idx_ketersediaan_person_event"`
	Person         Person             `gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EventID        uuid.UUID          `gorm:"type:char(36);not null;index:idx_ketersediaan_person_event"`
	Event          Event              `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OccurrenceDate time.Time          `gorm:"type:date;not null;index:idx_ketersediaan_occurrence"`
	Status         KetersediaanStatus `gorm:"type:varchar(20);not null"`
	Notes          string             `gorm:"type:text"`

	Timestamp
}

// TableName specifies the table name for Ketersediaan
func (Ketersediaan) TableName() string {
	return "ketersediaan"
}
