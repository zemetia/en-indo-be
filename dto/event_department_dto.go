package dto

import "github.com/google/uuid"

// UpdateEventDepartmentsRequest is used to bulk update departments for an event
type UpdateEventDepartmentsRequest struct {
	DepartmentIDs []uuid.UUID `json:"department_ids" binding:"required"`
}

// EventDepartmentResponse represents a department associated with an event
type EventDepartmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
