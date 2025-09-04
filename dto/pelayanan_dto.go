package dto

import (
	"time"

	"github.com/google/uuid"
)

type PelayananRequest struct {
	Pelayanan    string    `json:"pelayanan" binding:"required"`
	Description  string    `json:"description"`
	DepartmentID uuid.UUID `json:"department_id" binding:"required"`
}

type UpdatePelayananRequest struct {
	Pelayanan    string    `json:"pelayanan"`
	Description  string    `json:"description"`
	DepartmentID uuid.UUID `json:"department_id"`
}

type PelayananSimpleResponse struct {
	ID          uuid.UUID `json:"id"`
	Pelayanan   string    `json:"pelayanan"`
	Description string    `json:"description"`
	IsPic       bool      `json:"is_pic"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PelayananResponse struct {
	ID          uuid.UUID         `json:"id"`
	Pelayanan   string            `json:"pelayanan"`
	Description string            `json:"description"`
	IsPic       bool              `json:"is_pic"`
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
	ID             uuid.UUID      `json:"id"`
	PersonID       uuid.UUID      `json:"person_id"`
	PersonName     string         `json:"person_name"`
	PelayananID    uuid.UUID      `json:"pelayanan_id"`
	Pelayanan      string         `json:"pelayanan"`
	PelayananIsPic bool           `json:"pelayanan_is_pic"`
	ChurchID       uuid.UUID      `json:"church_id"`
	ChurchName     string         `json:"church_name"`
	DepartmentID   uuid.UUID      `json:"department_id"`
	DepartmentName string         `json:"department_name"`
	HasUserAccount bool           `json:"has_user_account"`
	IsUserActive   bool           `json:"is_user_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type AssignPelayananRequest struct {
	PersonID    uuid.UUID `json:"person_id" binding:"required"`
	PelayananID uuid.UUID `json:"pelayanan_id" binding:"required"`
	ChurchID    uuid.UUID `json:"church_id" binding:"required"`
}

type PelayananAssignmentPaginationResponse struct {
	Data []PelayananAssignmentResponse `json:"data"`
	PaginationResponse
}
