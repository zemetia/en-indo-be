package entity

import "github.com/google/uuid"

type UserPelayananGereja struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:char(36);not null"`
	User        User      `gorm:"foreignKey:UserID"`
	PelayananID uuid.UUID `gorm:"type:char(36);not null"`
	Pelayanan   Pelayanan `gorm:"foreignKey:PelayananID"`
	ChurchID    uuid.UUID `gorm:"type:char(36);not null"`
	Church      Church    `gorm:"foreignKey:ChurchID"`
	IsPic       bool      `gorm:"type:boolean;not null; default:false"`

	Timestamp
}
