package entity

import "github.com/google/uuid"

type LifeGroup struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	WhatsAppLink string    `json:"whatsapp_link"`
	ChurchID     uuid.UUID `json:"church_id"`
	Church       Church    `json:"church"`
	LeaderID     uuid.UUID `json:"leader_id"`
	Leader       User      `json:"leader"`
	Members      []User    `gorm:"many2many:life_group_members;" json:"members"`
	Persons      []Person  `gorm:"many2many:life_group_persons;" json:"persons"`
	Timestamp
}
