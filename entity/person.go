package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	ID                uuid.UUID   `gorm:"type:char(36);primary_key;"`
	Nama              string      `gorm:"type:varchar(100);not null" json:"nama"`
	NamaLain          string      `gorm:"type:varchar(100);null" json:"nama_lain"`
	Gender            string      `gorm:"type:char(1);not null" json:"gender"`
	TempatLahir       string      `gorm:"type:varchar(100)" json:"tempat_lahir"`
	TanggalLahir      time.Time   `gorm:"type:date" json:"tanggal_lahir"`
	FaseHidup         string      `gorm:"type:varchar(50)" json:"fase_hidup"`
	StatusPerkawinan  string      `gorm:"type:varchar(50)" json:"status_perkawinan"`
	NamaPasangan      string      `gorm:"type:varchar(100);null" json:"nama_pasangan"`
	PasanganID        *uuid.UUID  `gorm:"type:char(36);null" json:"pasangan_id"`
	Pasangan          *Person     `gorm:"foreignKey:PasanganID"`
	TanggalPerkawinan time.Time   `gorm:"type:date;null" json:"tanggal_perkawinan"`
	Alamat            string      `gorm:"type:text" json:"alamat"`
	NomorTelepon      string      `gorm:"type:varchar(20)" json:"nomor_telepon"`
	Email             string      `gorm:"type:varchar(100)" json:"email"`
	Ayah              string      `gorm:"type:varchar(100)" json:"ayah"`
	Ibu               string      `gorm:"type:varchar(100)" json:"ibu"`
	Kerinduan         string      `gorm:"type:text" json:"kerinduan"`
	KomitmenBerjemaat string      `gorm:"type:text" json:"komitmen_berjemaat"`
	Status            string      `gorm:"type:varchar(50)" json:"status"`
	IsAktif           bool        `gorm:"type:boolean;default:true"`
	KodeJemaat        string      `gorm:"type:varchar(50)" json:"kode_jemaat"`
	ChurchID          uuid.UUID   `gorm:"type:char(36);not null" json:"church_id"`
	Church            Church      `gorm:"foreignKey:ChurchID"`
	LifeGroups        []LifeGroup `gorm:"many2many:life_group_persons;"`
	KabupatenID       uint        `gorm:"type:int;not null" json:"kabupaten_id"`
	Kabupaten         Kabupaten   `gorm:"foreignKey:KabupatenID"`

	Timestamp
}

func (p *Person) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
