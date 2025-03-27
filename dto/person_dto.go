package dto

import (
	"time"

	"github.com/google/uuid"
)

type PersonRequest struct {
	Nama              string    `json:"nama" binding:"required"`
	NamaLain          string    `json:"nama_lain"`
	Gender            string    `json:"gender"`
	TempatLahir       string    `json:"tempat_lahir"`
	TanggalLahir      time.Time `json:"tanggal_lahir"`
	FaseHidup         string    `json:"fase_hidup"`
	StatusPerkawinan  string    `json:"status_perkawinan"`
	Pasangan          string    `json:"pasangan"`
	TanggalPerkawinan string    `json:"tanggal_perkawinan"`
	NomorTelepon      string    `json:"nomor_telepon"`
	Email             string    `json:"email"`
	Gereja            string    `json:"gereja"`
	Ayah              string    `json:"ayah"`
	Ibu               string    `json:"ibu"`
	Kerinduan         string    `json:"kerinduan"`
	KomitmenBerjemaat string    `json:"komitmen_berjemaat"`
	DateAdded         time.Time `json:"date_added"`
	Status            string    `json:"status"`
	TagListID         string    `json:"tag_list_id"`
	KodeJemaat        string    `json:"kode_jemaat"`
	ChurchID          uuid.UUID `json:"church_id" binding:"required"`
}

type PersonResponse struct {
	ID                uuid.UUID `json:"id"`
	PersonID          string    `json:"person_id"`
	Nama              string    `json:"nama"`
	NamaLain          string    `json:"nama_lain"`
	Gender            string    `json:"gender"`
	TempatLahir       string    `json:"tempat_lahir"`
	TanggalLahir      time.Time `json:"tanggal_lahir"`
	FaseHidup         string    `json:"fase_hidup"`
	StatusPerkawinan  string    `json:"status_perkawinan"`
	Pasangan          string    `json:"pasangan"`
	TanggalPerkawinan string    `json:"tanggal_perkawinan"`
	NomorTelepon      string    `json:"nomor_telepon"`
	Email             string    `json:"email"`
	Gereja            string    `json:"gereja"`
	Ayah              string    `json:"ayah"`
	Ibu               string    `json:"ibu"`
	Kerinduan         string    `json:"kerinduan"`
	KomitmenBerjemaat string    `json:"komitmen_berjemaat"`
	DateAdded         time.Time `json:"date_added"`
	Status            string    `json:"status"`
	TagListID         string    `json:"tag_list_id"`
	KodeJemaat        string    `json:"kode_jemaat"`
	ChurchID          uuid.UUID `json:"church_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
