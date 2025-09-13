package entity

import (
	"github.com/google/uuid"
)

type DiscipleshipJourney struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`

	Timestamp
}
