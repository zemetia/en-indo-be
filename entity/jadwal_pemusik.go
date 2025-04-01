package entity

import "github.com/google/uuid"

type JadwalPemusik struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	PersonID  uuid.UUID `gorm:"type:char(36);not null"`
	Person    Person    `gorm:"foreignKey:PersonID"`
	EventId   uuid.UUID `gorm:"type:char(36);not null"`
	Event     Event     `gorm:"foreignKey:EventId"`
	AlatMusik string    `gorm:"type:varchar(255);not null"`

	Timestamp
}
