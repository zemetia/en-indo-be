package entity

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key" json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_has_permissions;"`

	Timestamp
}
