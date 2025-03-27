package entity

import "github.com/google/uuid"

type Church struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Timestamp
}
