package dto

import "github.com/google/uuid"

type ChurchRequest struct {
	Name       string    `json:"name" binding:"required"`
	Address    string    `json:"address" binding:"required"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Website    string    `json:"website"`
	CityID     uuid.UUID `json:"city_id" binding:"required"`
	ProvinceID uuid.UUID `json:"province_id" binding:"required"`
}

type ChurchResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Website    string    `json:"website"`
	CityID     uuid.UUID `json:"city_id"`
	ProvinceID uuid.UUID `json:"province_id"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}
