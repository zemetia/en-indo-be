package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LifeGroupVisitorMember struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	LifeGroupID uuid.UUID `gorm:"type:char(36);not null" json:"life_group_id"`
	LifeGroup   LifeGroup `gorm:"foreignKey:LifeGroupID" json:"life_group"`
	VisitorID   uuid.UUID `gorm:"type:char(36);not null" json:"visitor_id"`
	Visitor     Visitor   `gorm:"foreignKey:VisitorID" json:"visitor"`
	IsActive    bool      `gorm:"type:boolean;default:true" json:"is_active"`
	JoinedDate  time.Time `gorm:"type:datetime;not null" json:"joined_date"`

	Timestamp
}

func (lgvm *LifeGroupVisitorMember) BeforeCreate(tx *gorm.DB) error {
	if lgvm.ID == uuid.Nil {
		lgvm.ID = uuid.New()
	}
	if lgvm.JoinedDate.IsZero() {
		lgvm.JoinedDate = time.Now()
	}
	return nil
}