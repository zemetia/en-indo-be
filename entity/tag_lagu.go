package entity

import "github.com/google/uuid"

type TagLagu struct {
	ID   uuid.UUID `gorm:"type:char(36);primary_key"`
	Nama string    `gorm:"type:varchar(255);not null"`

	Timestamp
}
