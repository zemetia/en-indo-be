package entity

import (
	"github.com/google/uuid"
)

type Church struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Address     string    `gorm:"type:text;not null"`
	Phone       string    `gorm:"type:varchar(20)"`
	Email       string    `gorm:"type:varchar(255)"`
	Website     string    `gorm:"type:varchar(255)"`
	KabupatenID uint      `gorm:"type:int;not null"`
	Kabupaten   Kabupaten `gorm:"foreignKey:KabupatenID"`

	Timestamp
}
