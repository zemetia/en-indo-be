package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type LifeGroupPersonMemberService interface {
	AddPersonMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.AddPersonMemberRequest) (*dto.PersonMemberResponse, error)
	GetPersonMembers(ctx context.Context, lifeGroupID uuid.UUID) ([]dto.PersonMemberResponse, error)
	UpdatePersonMemberPosition(ctx context.Context, lifeGroupID uuid.UUID, req *dto.UpdatePersonMemberPositionRequest) (*dto.PersonMemberResponse, error)
	RemovePersonMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.RemovePersonMemberRequest) error
	GetPersonMemberByID(ctx context.Context, memberID uuid.UUID) (*dto.PersonMemberResponse, error)
	GetPersonLifeGroups(ctx context.Context, personID uuid.UUID) ([]dto.PersonMemberResponse, error)
	ValidatePositionChange(ctx context.Context, lifeGroupID uuid.UUID, position entity.PersonMemberPosition) error
	GetLeadershipStructure(ctx context.Context, lifeGroupID uuid.UUID) (*dto.LeadershipStructureResponse, error)
}

type lifeGroupPersonMemberService struct {
	personMemberRepo repository.LifeGroupPersonMemberRepository
	personRepo       repository.PersonRepository
	lifeGroupRepo    repository.LifeGroupRepository
}


func NewLifeGroupPersonMemberService(
	personMemberRepo repository.LifeGroupPersonMemberRepository,
	personRepo repository.PersonRepository,
	lifeGroupRepo repository.LifeGroupRepository,
) LifeGroupPersonMemberService {
	return &lifeGroupPersonMemberService{
		personMemberRepo: personMemberRepo,
		personRepo:       personRepo,
		lifeGroupRepo:    lifeGroupRepo,
	}
}

func (s *lifeGroupPersonMemberService) AddPersonMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.AddPersonMemberRequest) (*dto.PersonMemberResponse, error) {
	// Validate person exists
	_, err := s.personRepo.GetByID(ctx, req.PersonID)
	if err != nil {
		return nil, fmt.Errorf("person not found: %w", err)
	}

	// Validate lifegroup exists
	_, err = s.lifeGroupRepo.GetByID(lifeGroupID)
	if err != nil {
		return nil, fmt.Errorf("lifegroup not found: %w", err)
	}

	// Check if person is already a member
	existingMember, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, req.PersonID)
	if err == nil && existingMember != nil {
		return nil, errors.New("person is already a member of this lifegroup")
	}

	// Validate position
	if err := s.ValidatePositionChange(ctx, lifeGroupID, req.Position); err != nil {
		return nil, err
	}

	// Create member
	member := &entity.LifeGroupPersonMember{
		ID:          uuid.New(),
		LifeGroupID: lifeGroupID,
		PersonID:    req.PersonID,
		Position:    req.Position,
		IsActive:    true,
	}

	if err := s.personMemberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to create person member: %w", err)
	}

	// Get created member with preloaded data
	createdMember, err := s.personMemberRepo.GetByID(ctx, member.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created member: %w", err)
	}

	return s.toPersonMemberResponse(createdMember), nil
}

func (s *lifeGroupPersonMemberService) GetPersonMembers(ctx context.Context, lifeGroupID uuid.UUID) ([]dto.PersonMemberResponse, error) {
	members, err := s.personMemberRepo.GetByLifeGroupID(ctx, lifeGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person members: %w", err)
	}

	var responses []dto.PersonMemberResponse
	for _, member := range members {
		responses = append(responses, *s.toPersonMemberResponse(&member))
	}

	return responses, nil
}

func (s *lifeGroupPersonMemberService) UpdatePersonMemberPosition(ctx context.Context, lifeGroupID uuid.UUID, req *dto.UpdatePersonMemberPositionRequest) (*dto.PersonMemberResponse, error) {
	// Get existing member
	existingMember, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, req.PersonID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	// Validate new position
	if err := s.ValidatePositionChange(ctx, lifeGroupID, req.Position); err != nil {
		return nil, err
	}

	// Update position
	if err := s.personMemberRepo.UpdatePosition(ctx, lifeGroupID, req.PersonID, req.Position); err != nil {
		return nil, fmt.Errorf("failed to update position: %w", err)
	}

	// Get updated member
	updatedMember, err := s.personMemberRepo.GetByID(ctx, existingMember.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated member: %w", err)
	}

	return s.toPersonMemberResponse(updatedMember), nil
}

func (s *lifeGroupPersonMemberService) RemovePersonMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.RemovePersonMemberRequest) error {
	member, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, req.PersonID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	return s.personMemberRepo.Delete(ctx, member.ID)
}

func (s *lifeGroupPersonMemberService) GetPersonMemberByID(ctx context.Context, memberID uuid.UUID) (*dto.PersonMemberResponse, error) {
	member, err := s.personMemberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	return s.toPersonMemberResponse(member), nil
}

func (s *lifeGroupPersonMemberService) GetPersonLifeGroups(ctx context.Context, personID uuid.UUID) ([]dto.PersonMemberResponse, error) {
	members, err := s.personMemberRepo.GetByPersonID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person lifegroups: %w", err)
	}

	var responses []dto.PersonMemberResponse
	for _, member := range members {
		responses = append(responses, *s.toPersonMemberResponse(&member))
	}

	return responses, nil
}

func (s *lifeGroupPersonMemberService) ValidatePositionChange(ctx context.Context, lifeGroupID uuid.UUID, position entity.PersonMemberPosition) error {
	switch position {
	case entity.PersonMemberPositionLeader:
		// Only one leader allowed - repository will handle demotion
		return nil
	case entity.PersonMemberPositionCoLeader:
		// Unlimited co-leaders allowed
		return nil
	case entity.PersonMemberPositionMember:
		// Unlimited members allowed
		return nil
	default:
		return fmt.Errorf("invalid position: %s", position)
	}
}

func (s *lifeGroupPersonMemberService) GetLeadershipStructure(ctx context.Context, lifeGroupID uuid.UUID) (*dto.LeadershipStructureResponse, error) {
	allMembers, err := s.personMemberRepo.GetByLifeGroupID(ctx, lifeGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	structure := &dto.LeadershipStructureResponse{
		CoLeaders: make([]dto.PersonMemberResponse, 0),
		Members:   make([]dto.PersonMemberResponse, 0),
	}

	for _, member := range allMembers {
		response := s.toPersonMemberResponse(&member)
		
		switch member.Position {
		case entity.PersonMemberPositionLeader:
			structure.Leader = response
		case entity.PersonMemberPositionCoLeader:
			structure.CoLeaders = append(structure.CoLeaders, *response)
		case entity.PersonMemberPositionMember:
			structure.Members = append(structure.Members, *response)
		}
	}

	return structure, nil
}

func (s *lifeGroupPersonMemberService) toPersonMemberResponse(member *entity.LifeGroupPersonMember) *dto.PersonMemberResponse {
	return &dto.PersonMemberResponse{
		ID:          member.ID,
		LifeGroupID: member.LifeGroupID,
		PersonID:    member.PersonID,
		Person:      member.Person,
		Position:    member.Position,
		IsActive:    member.IsActive,
		JoinedDate:  member.JoinedDate,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}
}