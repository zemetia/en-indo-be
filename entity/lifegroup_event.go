package entity

import (
	"time"

	"github.com/google/uuid"
)

type LifeGroupEvent struct {
	ID                       uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	LifeGroupID              uuid.UUID `gorm:"type:char(36);not null" json:"lifegroup_id"`
	LifeGroup                LifeGroup `gorm:"foreignKey:LifeGroupID" json:"lifegroup"`
	Title                    string    `gorm:"type:varchar(255);not null" json:"title"`
	Description              string    `gorm:"type:text" json:"description"`
	EventType                string    `gorm:"type:varchar(50);default:'Lifegroup'" json:"event_type"`
	ChurchID                 uuid.UUID `gorm:"type:char(36);not null" json:"church_id"`
	Church                   Church    `gorm:"foreignKey:ChurchID" json:"church"`
	StartTime                time.Time `gorm:"not null" json:"start_time"`
	EndTime                  time.Time `gorm:"not null" json:"end_time"`
	RecurrenceRule           string    `gorm:"type:text" json:"recurrence_rule"` // RFC5545 string or custom JSON
	ExpectedParticipantCount int       `gorm:"default:0" json:"expected_participant_count"`

	Timestamp
}
