package dto

import "github.com/google/uuid"

type ChurchRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	KabupatenID uint   `json:"kabupaten_id" binding:"required"`
}

type ChurchResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	KabupatenID uint      `json:"kabupaten_id"`
	Kabupaten   string    `json:"kabupaten"`
	ProvinsiID  uint      `json:"provinsi_id"`
	Provinsi    string    `json:"provinsi"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
