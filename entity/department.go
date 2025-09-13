package entity

import "github.com/google/uuid"

type Department struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;null"`

	Timestamp
}
