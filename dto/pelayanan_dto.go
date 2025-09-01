package dto

import (
	"time"

	"github.com/google/uuid"
)

type PelayananResponse struct {
	ID          uuid.UUID         `json:"id"`
	Pelayanan   string            `json:"pelayanan"`
	Description string            `json:"description"`
	Department  DepartmentResponse `json:"department"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type PersonHasPelayananResponse struct {
	PelayananID uuid.UUID `json:"pelayanan_id"`
	Pelayanan   string    `json:"pelayanan"`
	ChurchID    uuid.UUID `json:"church_id"`
	ChurchName  string    `json:"church_name"`
	IsPic       bool      `json:"is_pic"`
}

type PelayananAssignmentResponse struct {
	ID          uuid.UUID      `json:"id"`
	PersonID    uuid.UUID      `json:"person_id"`
	PersonName  string         `json:"person_name"`
	PelayananID uuid.UUID      `json:"pelayanan_id"`
	Pelayanan   string         `json:"pelayanan"`
	ChurchID    uuid.UUID      `json:"church_id"`
	ChurchName  string         `json:"church_name"`
	IsPic       bool           `json:"is_pic"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type AssignPelayananRequest struct {
	PersonID    uuid.UUID `json:"person_id" binding:"required"`
	PelayananID uuid.UUID `json:"pelayanan_id" binding:"required"`
	ChurchID    uuid.UUID `json:"church_id" binding:"required"`
	IsPic       bool      `json:"is_pic"`
}

type UpdatePelayananAssignmentRequest struct {
	IsPic bool `json:"is_pic"`
}

type PelayananAssignmentPaginationResponse struct {
	Data []PelayananAssignmentResponse `json:"data"`
	PaginationResponse
}
