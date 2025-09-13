package dto

import "github.com/google/uuid"

type DepartmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type DepartmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
