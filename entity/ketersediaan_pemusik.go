package entity

import "github.com/google/uuid"

type KetersediaanPemusik struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key"`
	PersonID     uuid.UUID `gorm:"type:char(36);not null"`
	Person       Person    `gorm:"foreignKey:PersonID"`
	Ketersediaan string    `gorm:"type:varchar(255);not null"`
	IsTerjadwal  bool      `gorm:"type:boolean;default:false"`
	EventId      uuid.UUID `gorm:"type:char(36);not null"`
	Event        Event     `gorm:"foreignKey:EventId"`

	Timestamp
}
