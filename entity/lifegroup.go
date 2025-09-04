package entity

import "github.com/google/uuid"

type LifeGroup struct {
	ID                uuid.UUID                    `gorm:"type:char(36);primary_key" json:"id"`
	Name              string                       `gorm:"type:varchar(255);not null" json:"name"`
	Location          string                       `gorm:"type:text;not null" json:"location"`
	WhatsAppLink      string                       `gorm:"type:text;not null" json:"whatsapp_link"`
	ChurchID          uuid.UUID                    `gorm:"type:char(36);not null" json:"church_id"`
	Church            Church                       `gorm:"foreignKey:ChurchID" json:"church"`
	LeaderID          uuid.UUID                    `gorm:"type:char(36);not null" json:"leader_id"`
	Leader            User                         `gorm:"foreignKey:LeaderID" json:"leader"`
	CoLeaderID        *uuid.UUID                   `gorm:"type:char(36);null" json:"co_leader_id"`
	CoLeader          *User                        `gorm:"foreignKey:CoLeaderID" json:"co_leader"`
	PersonMembers     []LifeGroupPersonMember      `gorm:"foreignKey:LifeGroupID" json:"person_members"`
	VisitorMembers    []LifeGroupVisitorMember     `gorm:"foreignKey:LifeGroupID" json:"visitor_members"`

	Timestamp
}
