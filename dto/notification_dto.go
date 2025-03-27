package dto

import "github.com/google/uuid"

type NotificationRequest struct {
	Title         string    `json:"title" binding:"required"`
	Message       string    `json:"message" binding:"required"`
	Type          string    `json:"type" binding:"required"` // "info", "success", "warning", "error"
	UserID        uuid.UUID `json:"user_id" binding:"required"`
	IsRead        bool      `json:"is_read"`
	ReferenceID   uuid.UUID `json:"reference_id"`   // ID dari entitas yang terkait (opsional)
	ReferenceType string    `json:"reference_type"` // Tipe entitas yang terkait (opsional)
}

type NotificationResponse struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Message       string    `json:"message"`
	Type          string    `json:"type"`
	UserID        uuid.UUID `json:"user_id"`
	IsRead        bool      `json:"is_read"`
	ReferenceID   uuid.UUID `json:"reference_id"`
	ReferenceType string    `json:"reference_type"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}
