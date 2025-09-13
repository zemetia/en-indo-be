package dto

import (
	"github.com/google/uuid"
)

// Visitor DTOs
type VisitorRequest struct {
	Name        string  `json:"name" binding:"required"`
	IGUsername  *string `json:"ig_username"`
	PhoneNumber *string `json:"phone_number"`
	KabupatenID *uint   `json:"kabupaten_id"`
}

type VisitorResponse struct {
	ID          uuid.UUID                    `json:"id"`
	Name        string                       `json:"name"`
	IGUsername  *string                      `json:"ig_username"`
	PhoneNumber *string                      `json:"phone_number"`
	KabupatenID *uint                        `json:"kabupaten_id"`
	Kabupaten   string                       `json:"kabupaten"`
	ProvinsiID  *uint                        `json:"provinsi_id"`
	Provinsi    string                       `json:"provinsi"`
	Information []VisitorInformationResponse `json:"information"`
	CreatedAt   string                       `json:"created_at"`
	UpdatedAt   string                       `json:"updated_at"`
}

type VisitorSimpleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IGUsername  *string   `json:"ig_username"`
	PhoneNumber *string   `json:"phone_number"`
	KabupatenID *uint     `json:"kabupaten_id"`
	Kabupaten   string    `json:"kabupaten"`
	ProvinsiID  *uint     `json:"provinsi_id"`
	Provinsi    string    `json:"provinsi"`
}

type VisitorSearchDto struct {
	Name        *string `json:"name" form:"name"`
	IGUsername  *string `json:"ig_username" form:"ig_username"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	KabupatenID *uint   `json:"kabupaten_id" form:"kabupaten_id"`
}

// Visitor Information DTOs
type VisitorInformationRequest struct {
	VisitorID uuid.UUID `json:"visitor_id" binding:"required"`
	Label     string    `json:"label" binding:"required"`
	Value     string    `json:"value" binding:"required"`
}

type VisitorInformationResponse struct {
	ID        uuid.UUID             `json:"id"`
	VisitorID uuid.UUID             `json:"visitor_id"`
	Visitor   VisitorSimpleResponse `json:"visitor,omitempty"`
	Label     string                `json:"label"`
	Value     string                `json:"value"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

type VisitorInformationSimpleResponse struct {
	ID    uuid.UUID `json:"id"`
	Label string    `json:"label"`
	Value string    `json:"value"`
}
