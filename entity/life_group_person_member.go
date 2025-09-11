package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonMemberPosition string

const (
	PersonMemberPositionLeader   PersonMemberPosition = "LEADER"
	PersonMemberPositionCoLeader PersonMemberPosition = "CO_LEADER"
	PersonMemberPositionMember   PersonMemberPosition = "MEMBER"
)

type LifeGroupPersonMember struct {
	ID          uuid.UUID            `gorm:"type:char(36);primary_key" json:"id"`
	LifeGroupID uuid.UUID            `gorm:"type:char(36);not null" json:"life_group_id"`
	LifeGroup   LifeGroup            `gorm:"foreignKey:LifeGroupID" json:"life_group"`
	PersonID    uuid.UUID            `gorm:"type:char(36);not null" json:"person_id"`
	Person      Person               `gorm:"foreignKey:PersonID" json:"person"`
	Position    PersonMemberPosition `gorm:"type:varchar(20);not null;default:'MEMBER'" json:"position"`
	IsActive    bool                 `gorm:"type:boolean;default:true" json:"is_active"`
	JoinedDate  time.Time            `gorm:"type:datetime;not null" json:"joined_date"`

	Timestamp
}

func (lgpm *LifeGroupPersonMember) BeforeCreate(tx *gorm.DB) error {
	if lgpm.ID == uuid.Nil {
		lgpm.ID = uuid.New()
	}
	if lgpm.JoinedDate.IsZero() {
		lgpm.JoinedDate = time.Now()
	}
	return nil
}
