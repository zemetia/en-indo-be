package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
)

type LifeGroupRequest struct {
	Name         string `json:"name" binding:"required"`
	Location     string `json:"location" binding:"required"`
	WhatsAppLink string `json:"whatsapp_link"`
	ChurchID     string `json:"church_id" binding:"required"`
}

type LifeGroupResponse struct {
	ID             uuid.UUID                       `json:"id"`
	Name           string                          `json:"name"`
	Location       string                          `json:"location"`
	WhatsAppLink   string                          `json:"whatsapp_link"`
	ChurchID       uuid.UUID                       `json:"church_id"`
	Church         entity.Church                   `json:"church"`
	PersonMembers  []entity.LifeGroupPersonMember  `json:"person_members"`
	VisitorMembers []entity.LifeGroupVisitorMember `json:"visitor_members"`
	CreatedAt      string                          `json:"created_at"`
	UpdatedAt      string                          `json:"updated_at"`
}

type UpdateLeaderRequest struct {
	LeaderID uuid.UUID `json:"leader_id" binding:"required"`
}

type LifeGroupSimpleResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	LeaderName string    `json:"leader_name"`
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role" binding:"required"`
}

type RemoveMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type UpdateMemberRoleRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role" binding:"required"`
}

type BatchChurchLifeGroupsRequest struct {
	ChurchIDs []uuid.UUID `json:"church_ids" binding:"required"`
}

type BatchChurchLifeGroupsResponse struct {
	ChurchID   uuid.UUID           `json:"church_id"`
	ChurchName string              `json:"church_name"`
	LifeGroups []LifeGroupResponse `json:"lifegroups"`
	Error      *string             `json:"error,omitempty"`
}

type AddPersonMemberRequest struct {
	PersonID uuid.UUID                   `json:"person_id" binding:"required"`
	Position entity.PersonMemberPosition `json:"position" binding:"required"`
}

type UpdatePersonMemberPositionRequest struct {
	PersonID uuid.UUID                   `json:"person_id" binding:"required"`
	Position entity.PersonMemberPosition `json:"position" binding:"required"`
}

type RemovePersonMemberRequest struct {
	PersonID uuid.UUID `json:"person_id" binding:"required"`
}

type PersonMemberResponse struct {
	ID          uuid.UUID                   `json:"id"`
	LifeGroupID uuid.UUID                   `json:"life_group_id"`
	PersonID    uuid.UUID                   `json:"person_id"`
	Person      entity.Person               `json:"person"`
	Position    entity.PersonMemberPosition `json:"position"`
	IsActive    bool                        `json:"is_active"`
	JoinedDate  time.Time                   `json:"joined_date"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type AddVisitorMemberRequest struct {
	VisitorID uuid.UUID `json:"visitor_id" binding:"required"`
}

type RemoveVisitorMemberRequest struct {
	VisitorID uuid.UUID `json:"visitor_id" binding:"required"`
}

type VisitorMemberResponse struct {
	ID          uuid.UUID      `json:"id"`
	LifeGroupID uuid.UUID      `json:"life_group_id"`
	VisitorID   uuid.UUID      `json:"visitor_id"`
	Visitor     entity.Visitor `json:"visitor"`
	IsActive    bool           `json:"is_active"`
	JoinedDate  time.Time      `json:"joined_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type LeadershipStructureResponse struct {
	Leader    *PersonMemberResponse  `json:"leader"`
	CoLeaders []PersonMemberResponse `json:"co_leaders"`
	Members   []PersonMemberResponse `json:"members"`
}

type AddPersonMembersBatchRequest struct {
	PersonIDs []uuid.UUID `json:"person_ids" binding:"required,min=1"`
}

type AddVisitorMembersBatchRequest struct {
	VisitorIDs []uuid.UUID `json:"visitor_ids" binding:"required,min=1"`
}

type BatchOperationResult struct {
	TotalRequested int                    `json:"total_requested"`
	Successful     int                    `json:"successful"`
	Failed         int                    `json:"failed"`
	Errors         []BatchOperationError  `json:"errors,omitempty"`
}

type BatchOperationError struct {
	ID    uuid.UUID `json:"id"`
	Error string    `json:"error"`
}
