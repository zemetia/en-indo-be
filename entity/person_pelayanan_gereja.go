package entity

import "github.com/google/uuid"

type PersonPelayananGereja struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	PersonID    uuid.UUID `gorm:"type:char(36);not null"`
	Person      Person    `gorm:"foreignKey:PersonID"`
	PelayananID uuid.UUID `gorm:"type:char(36);not null"`
	Pelayanan   Pelayanan `gorm:"foreignKey:PelayananID"`
	ChurchID    uuid.UUID `gorm:"type:char(36);not null"`
	Church      Church    `gorm:"foreignKey:ChurchID"`

	Timestamp
}
