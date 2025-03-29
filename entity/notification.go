package entity

import "github.com/google/uuid"

type Notification struct {
	ID       uuid.UUID  `gorm:"type:char(36);primary_key" json:"id"`
	Title    string     `gorm:"type:varchar(200);not null" json:"title"`
	Message  string     `gorm:"type:text;not null" json:"message"`
	Type     string     `gorm:"type:varchar(50);not null" json:"type"` // contoh: "info", "warning", "success", "error"
	IsRead   bool       `gorm:"type:boolean;not null;default:false" json:"is_read"`
	UserID   uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	User     User       `gorm:"foreignKey:UserID" json:"user"`
	ChurchID *uuid.UUID `gorm:"type:char(36);null"`  // optional nullable
	Church   *Church    `gorm:"foreignKey:ChurchID"` // optional nullable
	Timestamp
}
