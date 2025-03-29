package entity

import "github.com/google/uuid"

type Pelayanan struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Pelayanan   string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`

	Timestamp
}
