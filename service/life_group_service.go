package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type LifeGroupService interface {
	Create(req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error)
	GetAll() ([]dto.LifeGroupResponse, error)
	GetByID(id uuid.UUID) (*dto.LifeGroupResponse, error)
	Update(id uuid.UUID, req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error)
	Delete(id uuid.UUID) error
	UpdateLeader(id uuid.UUID, req *dto.UpdateLeaderRequest) (*dto.LifeGroupResponse, error)
	Search(ctx context.Context, search *dto.PersonSearchDto) ([]dto.LifeGroupResponse, error)
	GetByChurchID(churchID uuid.UUID) ([]dto.LifeGroupResponse, error)
	GetByUserID(userID uuid.UUID) ([]dto.LifeGroupResponse, error)
	CheckUserPermission(userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	GetByMultipleChurchIDs(churchIDs []uuid.UUID) ([]dto.BatchChurchLifeGroupsResponse, error)
}

type lifeGroupService struct {
	lifeGroupRepo repository.LifeGroupRepository
}

func NewLifeGroupService(lifeGroupRepo repository.LifeGroupRepository) LifeGroupService {
	return &lifeGroupService{
		lifeGroupRepo: lifeGroupRepo,
	}
}

func (s *lifeGroupService) Create(req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error) {
	lifeGroup := &entity.LifeGroup{
		ID:           uuid.New(),
		Name:         req.Name,
		Location:     req.Location,
		WhatsAppLink: req.WhatsAppLink,
		ChurchID:     req.ChurchID,
		LeaderID:     req.LeaderID,
		CoLeaderID:   req.CoLeaderID,
	}

	if err := s.lifeGroupRepo.Create(lifeGroup); err != nil {
		return nil, err
	}

	// Member management is now handled through separate PersonMember and VisitorMember APIs

	// Get updated life group with all relations
	updatedLifeGroup, err := s.lifeGroupRepo.GetByID(lifeGroup.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updatedLifeGroup), nil
}

func (s *lifeGroupService) GetAll() ([]dto.LifeGroupResponse, error) {
	lifeGroups, err := s.lifeGroupRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range lifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

func (s *lifeGroupService) Search(ctx context.Context, search *dto.PersonSearchDto) ([]dto.LifeGroupResponse, error) {
	lifegroups, err := s.lifeGroupRepo.Search(ctx, search)
	if err != nil {
		return nil, err
	}

	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range lifegroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

func (s *lifeGroupService) GetByID(id uuid.UUID) (*dto.LifeGroupResponse, error) {
	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(lifeGroup), nil
}

func (s *lifeGroupService) Update(id uuid.UUID, req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error) {
	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	lifeGroup.Name = req.Name
	lifeGroup.Location = req.Location
	lifeGroup.WhatsAppLink = req.WhatsAppLink
	lifeGroup.ChurchID = req.ChurchID
	lifeGroup.LeaderID = req.LeaderID
	lifeGroup.CoLeaderID = req.CoLeaderID

	if err := s.lifeGroupRepo.Update(lifeGroup); err != nil {
		return nil, err
	}

	// Member management is now handled through separate PersonMember and VisitorMember APIs

	// Get updated life group with all relations
	updatedLifeGroup, err := s.lifeGroupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updatedLifeGroup), nil
}

func (s *lifeGroupService) Delete(id uuid.UUID) error {
	return s.lifeGroupRepo.Delete(id)
}

func (s *lifeGroupService) UpdateLeader(id uuid.UUID, req *dto.UpdateLeaderRequest) (*dto.LifeGroupResponse, error) {
	if err := s.lifeGroupRepo.UpdateLeader(id, req.LeaderID); err != nil {
		return nil, err
	}

	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(lifeGroup), nil
}


func (s *lifeGroupService) toResponse(lifeGroup *entity.LifeGroup) *dto.LifeGroupResponse {
	response := &dto.LifeGroupResponse{
		ID:           lifeGroup.ID,
		Name:         lifeGroup.Name,
		Location:     lifeGroup.Location,
		WhatsAppLink: lifeGroup.WhatsAppLink,
		ChurchID:     lifeGroup.ChurchID,
		Church:       lifeGroup.Church,
		LeaderID:     lifeGroup.LeaderID,
		Leader:       lifeGroup.Leader,
		PersonMembers:  lifeGroup.PersonMembers,
		VisitorMembers: lifeGroup.VisitorMembers,
		CreatedAt:    lifeGroup.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    lifeGroup.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	if lifeGroup.CoLeaderID != nil {
		response.CoLeaderID = lifeGroup.CoLeaderID
		response.CoLeader = lifeGroup.CoLeader
	}
	
	return response
}

func (s *lifeGroupService) GetByChurchID(churchID uuid.UUID) ([]dto.LifeGroupResponse, error) {
	lifeGroups, err := s.lifeGroupRepo.GetByChurchID(churchID)
	if err != nil {
		return nil, err
	}

	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range lifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

func (s *lifeGroupService) GetByUserID(userID uuid.UUID) ([]dto.LifeGroupResponse, error) {
	lifeGroups, err := s.lifeGroupRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range lifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

func (s *lifeGroupService) CheckUserPermission(userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// Check if user is leader or co-leader of the lifegroup
	lifeGroup, err := s.lifeGroupRepo.GetByID(lifeGroupID)
	if err != nil {
		return false, err
	}

	// Check if user is leader or co-leader
	if lifeGroup.LeaderID == userID {
		return true, nil
	}
	if lifeGroup.CoLeaderID != nil && *lifeGroup.CoLeaderID == userID {
		return true, nil
	}

	return false, nil
}

func (s *lifeGroupService) GetByMultipleChurchIDs(churchIDs []uuid.UUID) ([]dto.BatchChurchLifeGroupsResponse, error) {
	responses := make([]dto.BatchChurchLifeGroupsResponse, 0, len(churchIDs))

	for _, churchID := range churchIDs {
		response := dto.BatchChurchLifeGroupsResponse{
			ChurchID: churchID,
		}

		// Get lifegroups for this church
		lifeGroups, err := s.lifeGroupRepo.GetByChurchID(churchID)
		if err != nil {
			errorMsg := err.Error()
			response.Error = &errorMsg
		} else {
			// Convert to response DTOs
			var lifegroupResponses []dto.LifeGroupResponse
			for _, lifeGroup := range lifeGroups {
				lifegroupResponses = append(lifegroupResponses, *s.toResponse(&lifeGroup))
				// Set church name from the first lifegroup's church info
				if response.ChurchName == "" && lifeGroup.Church.Name != "" {
					response.ChurchName = lifeGroup.Church.Name
				}
			}
			response.LifeGroups = lifegroupResponses
		}

		responses = append(responses, response)
	}

	return responses, nil
}
