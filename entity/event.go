package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	EventName     string    `gorm:"type:varchar(255);not null"`
	EventDate     time.Time `gorm:"type:datetime;not null"`
	EventLocation string    `gorm:"type:varchar(255);not null"`
	Lagu          []Lagu    `gorm:"many2many:event_lagu;"`

	Timestamp
}
