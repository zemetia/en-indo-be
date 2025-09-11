package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VisitorInformation struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;unique" json:"id"`
	VisitorID uuid.UUID `gorm:"type:char(36);not null" json:"visitor_id"`
	Visitor   Visitor   `gorm:"foreignKey:VisitorID" json:"visitor"`
	Label     string    `gorm:"type:varchar(100);not null" json:"label"`
	Value     string    `gorm:"type:text;not null" json:"value"`

	TimestampHardDelete
}

func (vi *VisitorInformation) BeforeCreate(tx *gorm.DB) error {
	if vi.ID == uuid.Nil {
		vi.ID = uuid.New()
	}
	return nil
}
