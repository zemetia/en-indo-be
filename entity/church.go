package entity

import (
	"github.com/google/uuid"
)

type Church struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key;"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Address     string    `gorm:"type:text;not null"`
	ChurchCode  string    `gorm:"type:varchar(10);uniqueIndex;null" json:"church_code"`
	Phone       string    `gorm:"type:varchar(20);null"`
	Email       string    `gorm:"type:varchar(255);null"`
	Website     string    `gorm:"type:varchar(255);null"`
	Latitude    float64   `gorm:"type:decimal(10,8);null" json:"latitude"`
	Longitude   float64   `gorm:"type:decimal(11,8);null" json:"longitude"`
	KabupatenID uint      `gorm:"type:int;not null" json:"kabupaten_id"`
	Kabupaten   Kabupaten `gorm:"foreignKey:KabupatenID"`

	Timestamp
}
