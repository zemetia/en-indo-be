package dto

import "github.com/google/uuid"

type DepartmentRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	ChurchID    uuid.UUID `json:"church_id" binding:"required"`
}

type DepartmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ChurchID    uuid.UUID `json:"church_id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
