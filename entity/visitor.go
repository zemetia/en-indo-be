package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visitor struct {
	ID          uuid.UUID           `gorm:"type:char(36);primary_key;unique" json:"id"`
	Name        string              `gorm:"type:varchar(100);not null" json:"name"`
	IGUsername  *string             `gorm:"type:varchar(100);null" json:"ig_username"`
	PhoneNumber *string             `gorm:"type:varchar(20);null" json:"phone_number"`
	KabupatenID *uint               `gorm:"type:int;null" json:"kabupaten_id"`
	Kabupaten   Kabupaten           `gorm:"foreignKey:KabupatenID" json:"kabupaten"`
	Information []VisitorInformation `gorm:"foreignKey:VisitorID" json:"information"`

	Timestamp
}

func (v *Visitor) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}