package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	ID                uuid.UUID   `gorm:"type:char(36);primary_key;"`
	Nama              string      `gorm:"type:varchar(100);not null"`
	NamaLain          string      `gorm:"type:varchar(100)"`
	Gender            string      `gorm:"type:varchar(20)"`
	TempatLahir       string      `gorm:"type:varchar(100)"`
	TanggalLahir      time.Time   `gorm:"type:date"`
	FaseHidup         string      `gorm:"type:varchar(50)"`
	StatusPerkawinan  string      `gorm:"type:varchar(50)"`
	NamaPasangan      string      `gorm:"type:varchar(100);null"`
	PasanganID        *uuid.UUID  `gorm:"type:char(36);null"`
	Pasangan          *Person     `gorm:"foreignKey:PasanganID"`
	TanggalPerkawinan time.Time   `gorm:"type:date;null"`
	Alamat            string      `gorm:"type:text"`
	NomorTelepon      string      `gorm:"type:varchar(20)"`
	Email             string      `gorm:"type:varchar(100)"`
	Ayah              string      `gorm:"type:varchar(100)"`
	Ibu               string      `gorm:"type:varchar(100)"`
	Kerinduan         string      `gorm:"type:text"`
	KomitmenBerjemaat string      `gorm:"type:text"`
	Status            string      `gorm:"type:varchar(50)"`
	KodeJemaat        string      `gorm:"type:varchar(50)"`
	ChurchID          uuid.UUID   `gorm:"type:char(36);not null"`
	Church            Church      `gorm:"foreignKey:ChurchID"`
	LifeGroups        []LifeGroup `gorm:"many2many:life_group_persons;"`
	KabupatenID       uint        `gorm:"type:int;not null"`
	Kabupaten         Kabupaten   `gorm:"foreignKey:KabupatenID"`

	Timestamp
}

func (p *Person) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
