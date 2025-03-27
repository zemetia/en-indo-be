package entity

import "github.com/google/uuid"

type Department struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`
	ChurchID    uuid.UUID `gorm:"not null"`
	Church      Church    `gorm:"foreignKey:ChurchID"`
	Users       []User    `gorm:"many2many:user_departments;"`
	Timestamp
}
