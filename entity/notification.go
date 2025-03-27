package entity

import "github.com/google/uuid"

type Notification struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Title    string    `gorm:"type:varchar(200);not null" json:"title"`
	Message  string    `gorm:"type:text;not null" json:"message"`
	Type     string    `gorm:"type:varchar(50);not null" json:"type"` // contoh: "info", "warning", "success", "error"
	IsRead   bool      `gorm:"default:false" json:"is_read"`
	UserID   uuid.UUID `json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	ChurchID uuid.UUID `gorm:"not null"`
	Church   Church    `gorm:"foreignKey:ChurchID"`
	Timestamp
}
