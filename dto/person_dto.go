package dto

import (
	"time"

	"github.com/google/uuid"
)

type PersonRequest struct {
	Nama              string      `json:"nama" binding:"required"`
	NamaLain          string      `json:"nama_lain"`
	Gender            string      `json:"gender"`
	TempatLahir       string      `json:"tempat_lahir"`
	TanggalLahir      time.Time   `json:"tanggal_lahir"`
	FaseHidup         string      `json:"fase_hidup"`
	StatusPerkawinan  string      `json:"status_perkawinan"`
	NamaPasangan      string      `json:"nama_pasangan"`
	PasanganID        *uuid.UUID  `json:"pasangan_id"`
	TanggalPerkawinan time.Time   `json:"tanggal_perkawinan"`
	Alamat            string      `json:"alamat"`
	NomorTelepon      string      `json:"nomor_telepon"`
	Email             string      `json:"email"`
	Ayah              string      `json:"ayah"`
	Ibu               string      `json:"ibu"`
	Kerinduan         string      `json:"kerinduan"`
	KomitmenBerjemaat string      `json:"komitmen_berjemaat"`
	Status            string      `json:"status"`
	KodeJemaat        string      `json:"kode_jemaat"`
	ChurchID          uuid.UUID   `json:"church_id" binding:"required"`
	UserID            *uuid.UUID  `json:"user_id"`
	KabupatenID       uint        `json:"kabupaten_id" binding:"required"`
	LifeGroupIDs      []uuid.UUID `json:"life_group_ids"`
}

type PersonResponse struct {
	ID                uuid.UUID                 `json:"id"`
	Nama              string                    `json:"nama"`
	NamaLain          string                    `json:"nama_lain"`
	Gender            string                    `json:"gender"`
	TempatLahir       string                    `json:"tempat_lahir"`
	TanggalLahir      string                    `json:"tanggal_lahir"`
	FaseHidup         string                    `json:"fase_hidup"`
	StatusPerkawinan  string                    `json:"status_perkawinan"`
	NamaPasangan      string                    `json:"nama_pasangan"`
	PasanganID        *uuid.UUID                `json:"pasangan_id"`
	TanggalPerkawinan string                    `json:"tanggal_perkawinan"`
	Alamat            string                    `json:"alamat"`
	NomorTelepon      string                    `json:"nomor_telepon"`
	Email             string                    `json:"email"`
	Ayah              string                    `json:"ayah"`
	Ibu               string                    `json:"ibu"`
	Kerinduan         string                    `json:"kerinduan"`
	KomitmenBerjemaat string                    `json:"komitmen_berjemaat"`
	Status            string                    `json:"status"`
	KodeJemaat        string                    `json:"kode_jemaat"`
	ChurchID          uuid.UUID                 `json:"church_id"`
	Church            string                    `json:"church"`
	UserID            *uuid.UUID                `json:"user_id"`
	KabupatenID       uint                      `json:"kabupaten_id"`
	Kabupaten         string                    `json:"kabupaten"`
	LifeGroups        []LifeGroupSimpleResponse `json:"life_groups"`
	CreatedAt         string                    `json:"created_at"`
	UpdatedAt         string                    `json:"updated_at"`
}

type PersonSearchDto struct {
	Name        *string    `json:"name" form:"name"`
	ChurchID    *uuid.UUID `json:"church_id" form:"church_id"`
	KabupatenID *uint      `json:"kabupaten_id" form:"kabupaten_id"`
	UserID      *uuid.UUID `json:"user_id" form:"user_id"`
}
