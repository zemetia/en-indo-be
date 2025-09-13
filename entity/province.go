package entity

import (
	"github.com/google/uuid"
)

type Province struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name string    `gorm:"type:varchar(255);not null"`

	Timestamp
}
