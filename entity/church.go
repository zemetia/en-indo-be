package entity

import (
	"github.com/google/uuid"
)

type Church struct {
<<<<<<< HEAD
	ID          uuid.UUID `gorm:"type:char(36);primary_key;"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Address     string    `gorm:"type:text;not null"`
	Phone       string    `gorm:"type:varchar(20)"`
	Email       string    `gorm:"type:varchar(255)"`
	Website     string    `gorm:"type:varchar(255)"`
	KabupatenID uint      `gorm:"type:int;not null"`
	Kabupaten   Kabupaten `gorm:"foreignKey:KabupatenID"`
=======
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
>>>>>>> 5ed98ee3f618ce09f59ee7b2565092ff36194b0f

	Timestamp
}
