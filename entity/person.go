package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	PersonID          string      `gorm:"type:varchar(36);unique;index;default:uuid_generate_v4()"`
	Nama              string      `gorm:"type:varchar(100);not null"`
	NamaLain          string      `gorm:"type:varchar(100)"`
	Gender            string      `gorm:"type:varchar(20)"`
	TempatLahir       string      `gorm:"type:varchar(100)"`
	TanggalLahir      time.Time   `gorm:"type:date"`
	FaseHidup         string      `gorm:"type:varchar(50)"`
	StatusPerkawinan  string      `gorm:"type:varchar(50)"`
	Pasangan          string      `gorm:"type:varchar(100)"`
	TanggalPerkawinan string      `gorm:"type:varchar(50)"`
	NomorTelepon      string      `gorm:"type:varchar(20)"`
	Email             string      `gorm:"type:varchar(100)"`
	Gereja            string      `gorm:"type:varchar(100)"`
	Ayah              string      `gorm:"type:varchar(100)"`
	Ibu               string      `gorm:"type:varchar(100)"`
	Kerinduan         string      `gorm:"type:text"`
	KomitmenBerjemaat string      `gorm:"type:text"`
	DateAdded         time.Time   `gorm:"type:date"`
	Status            string      `gorm:"type:varchar(50)"`
	TagListID         string      `gorm:"type:varchar(100)"`
	KodeJemaat        string      `gorm:"type:varchar(50)"`
	ChurchID          uuid.UUID   `gorm:"not null" json:"church_id"`
	Church            Church      `gorm:"foreignKey:ChurchID" json:"church"`
	LifeGroups        []LifeGroup `gorm:"many2many:life_group_persons;"`
	UserID            *uuid.UUID  `json:"user_id"`
	User              *User       `json:"user"`
	Timestamp
}

func (p *Person) BeforeCreate(tx *gorm.DB) error {
	if p.PersonID == "" {
		p.PersonID = uuid.New().String()
	}
	return nil
}
