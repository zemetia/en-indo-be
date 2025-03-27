package entity

import (
	"github.com/google/uuid"
)

type Church struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Address    string    `gorm:"type:text;not null"`
	Phone      string    `gorm:"type:varchar(20)"`
	Email      string    `gorm:"type:varchar(255)"`
	Website    string    `gorm:"type:varchar(255)"`
	CityID     uuid.UUID `gorm:"type:uuid;not null"`
	ProvinceID uuid.UUID `gorm:"type:uuid;not null"`
	City       City      `gorm:"foreignKey:CityID"`
	Province   Province  `gorm:"foreignKey:ProvinceID"`

	Timestamp
}
