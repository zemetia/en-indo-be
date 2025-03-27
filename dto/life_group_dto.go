package dto

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
)

type LifeGroupRequest struct {
	Name         string      `json:"name" binding:"required"`
	Location     string      `json:"location" binding:"required"`
	WhatsAppLink string      `json:"whatsapp_link"`
	ChurchID     uuid.UUID   `json:"church_id" binding:"required"`
	LeaderID     uuid.UUID   `json:"leader_id" binding:"required"`
	MemberIDs    []uuid.UUID `json:"member_ids"`
	PersonIDs    []uuid.UUID `json:"person_ids"`
}

type LifeGroupResponse struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Location     string          `json:"location"`
	WhatsAppLink string          `json:"whatsapp_link"`
	ChurchID     uuid.UUID       `json:"church_id"`
	Church       entity.Church   `json:"church"`
	LeaderID     uuid.UUID       `json:"leader_id"`
	Leader       entity.User     `json:"leader"`
	Members      []entity.User   `json:"members"`
	Persons      []entity.Person `json:"persons"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type UpdateLeaderRequest struct {
	LeaderID uuid.UUID `json:"leader_id" binding:"required"`
}

type UpdateMembersRequest struct {
	MemberIDs []uuid.UUID `json:"member_ids" binding:"required"`
}

type UpdatePersonsRequest struct {
	PersonIDs []uuid.UUID `json:"person_ids" binding:"required"`
}
