package dto

import "github.com/google/uuid"

type PelayananResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PersonHasPelayananResponse struct {
	PelayananID uuid.UUID `json:"pelayanan_id"`
	Pelayanan   string    `json:"pelayanan"`
	ChurchID    uuid.UUID `json:"church_id"`
	ChurchName  string    `json:"church_name"`
	IsPic       bool      `json:"is_pic"`
}
