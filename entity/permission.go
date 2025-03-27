package entity

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp
}
