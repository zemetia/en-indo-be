package dto

import "github.com/google/uuid"

type EventTypeRequest struct {
	Name        string `json:"name"`                            // Auto-generated from DisplayName if not provided
	DisplayName string `json:"display_name" binding:"required"` // Required display name
	Description string `json:"description"`                     // Optional description
	Color       string `json:"color"`                           // Hex color code (default: #2980b9)
	IsActive    *bool  `json:"is_active"`                       // Pointer to distinguish between false and not provided
	SortOrder   int    `json:"sort_order"`                      // Sort order (default: 0)
}

type EventTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	IsActive    bool      `json:"is_active"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
