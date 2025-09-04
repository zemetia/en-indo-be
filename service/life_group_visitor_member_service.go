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

type LifeGroupVisitorMemberService interface {
	AddVisitorMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.AddVisitorMemberRequest) (*dto.VisitorMemberResponse, error)
	GetVisitorMembers(ctx context.Context, lifeGroupID uuid.UUID) ([]dto.VisitorMemberResponse, error)
	RemoveVisitorMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.RemoveVisitorMemberRequest) error
	GetVisitorMemberByID(ctx context.Context, memberID uuid.UUID) (*dto.VisitorMemberResponse, error)
	GetVisitorLifeGroups(ctx context.Context, visitorID uuid.UUID) ([]dto.VisitorMemberResponse, error)
}

type lifeGroupVisitorMemberService struct {
	visitorMemberRepo repository.LifeGroupVisitorMemberRepository
	visitorRepo       repository.VisitorRepository
	lifeGroupRepo     repository.LifeGroupRepository
}

func NewLifeGroupVisitorMemberService(
	visitorMemberRepo repository.LifeGroupVisitorMemberRepository,
	visitorRepo repository.VisitorRepository,
	lifeGroupRepo repository.LifeGroupRepository,
) LifeGroupVisitorMemberService {
	return &lifeGroupVisitorMemberService{
		visitorMemberRepo: visitorMemberRepo,
		visitorRepo:       visitorRepo,
		lifeGroupRepo:     lifeGroupRepo,
	}
}

func (s *lifeGroupVisitorMemberService) AddVisitorMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.AddVisitorMemberRequest) (*dto.VisitorMemberResponse, error) {
	// Validate visitor exists
	_, err := s.visitorRepo.GetByID(ctx, req.VisitorID)
	if err != nil {
		return nil, fmt.Errorf("visitor not found: %w", err)
	}

	// Validate lifegroup exists
	_, err = s.lifeGroupRepo.GetByID(lifeGroupID)
	if err != nil {
		return nil, fmt.Errorf("lifegroup not found: %w", err)
	}

	// Check if visitor is already a member
	exists, err := s.visitorMemberRepo.ExistsByLifeGroupAndVisitorID(ctx, lifeGroupID, req.VisitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing membership: %w", err)
	}
	if exists {
		return nil, errors.New("visitor is already a member of this lifegroup")
	}

	// Create member
	member := &entity.LifeGroupVisitorMember{
		ID:          uuid.New(),
		LifeGroupID: lifeGroupID,
		VisitorID:   req.VisitorID,
		IsActive:    true,
	}

	if err := s.visitorMemberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to create visitor member: %w", err)
	}

	// Get created member with preloaded data
	createdMember, err := s.visitorMemberRepo.GetByID(ctx, member.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created member: %w", err)
	}

	return s.toVisitorMemberResponse(createdMember), nil
}

func (s *lifeGroupVisitorMemberService) GetVisitorMembers(ctx context.Context, lifeGroupID uuid.UUID) ([]dto.VisitorMemberResponse, error) {
	members, err := s.visitorMemberRepo.GetByLifeGroupID(ctx, lifeGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor members: %w", err)
	}

	var responses []dto.VisitorMemberResponse
	for _, member := range members {
		responses = append(responses, *s.toVisitorMemberResponse(&member))
	}

	return responses, nil
}

func (s *lifeGroupVisitorMemberService) RemoveVisitorMember(ctx context.Context, lifeGroupID uuid.UUID, req *dto.RemoveVisitorMemberRequest) error {
	member, err := s.visitorMemberRepo.GetByLifeGroupAndVisitorID(ctx, lifeGroupID, req.VisitorID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	return s.visitorMemberRepo.Delete(ctx, member.ID)
}

func (s *lifeGroupVisitorMemberService) GetVisitorMemberByID(ctx context.Context, memberID uuid.UUID) (*dto.VisitorMemberResponse, error) {
	member, err := s.visitorMemberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	return s.toVisitorMemberResponse(member), nil
}

func (s *lifeGroupVisitorMemberService) GetVisitorLifeGroups(ctx context.Context, visitorID uuid.UUID) ([]dto.VisitorMemberResponse, error) {
	members, err := s.visitorMemberRepo.GetByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor lifegroups: %w", err)
	}

	var responses []dto.VisitorMemberResponse
	for _, member := range members {
		responses = append(responses, *s.toVisitorMemberResponse(&member))
	}

	return responses, nil
}

func (s *lifeGroupVisitorMemberService) toVisitorMemberResponse(member *entity.LifeGroupVisitorMember) *dto.VisitorMemberResponse {
	return &dto.VisitorMemberResponse{
		ID:          member.ID,
		LifeGroupID: member.LifeGroupID,
		VisitorID:   member.VisitorID,
		Visitor:     member.Visitor,
		IsActive:    member.IsActive,
		JoinedDate:  member.JoinedDate,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
	}
}