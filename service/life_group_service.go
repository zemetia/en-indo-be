package service

// import (
// 	"github.com/google/uuid"
// 	"github.com/zemetia/en-indo-be/dto"
// 	"github.com/zemetia/en-indo-be/entity"
// 	"github.com/zemetia/en-indo-be/repository"
// )

// type LifeGroupService interface {
// 	Create(req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error)
// 	GetAll() ([]dto.LifeGroupResponse, error)
// 	GetByID(id uuid.UUID) (*dto.LifeGroupResponse, error)
// 	Update(id uuid.UUID, req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error)
// 	Delete(id uuid.UUID) error
// 	UpdateLeader(id uuid.UUID, req *dto.UpdateLeaderRequest) (*dto.LifeGroupResponse, error)
// 	UpdateMembers(id uuid.UUID, req *dto.UpdateMembersRequest) (*dto.LifeGroupResponse, error)
// 	GetByChurchID(churchID uuid.UUID) ([]dto.LifeGroupResponse, error)

// AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
// RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
// }

// type lifeGroupService struct {
// 	lifeGroupRepo repository.LifeGroupRepository
// }

// func NewLifeGroupService(lifeGroupRepo repository.LifeGroupRepository) LifeGroupService {
// 	return &lifeGroupService{
// 		lifeGroupRepo: lifeGroupRepo,
// 	}
// }

// func (s *lifeGroupService) Create(req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error) {
// 	lifeGroup := &entity.LifeGroup{
// 		ID:           uuid.New(),
// 		Name:         req.Name,
// 		Location:     req.Location,
// 		WhatsAppLink: req.WhatsAppLink,
// 		ChurchID:     req.ChurchID,
// 		LeaderID:     req.LeaderID,
// 	}

// 	if err := s.lifeGroupRepo.Create(lifeGroup); err != nil {
// 		return nil, err
// 	}

// 	// Update members and persons if provided
// 	if len(req.MemberIDs) > 0 {
// 		if err := s.lifeGroupRepo.UpdateMembers(lifeGroup.ID, req.MemberIDs); err != nil {
// 			return nil, err
// 		}
// 	}

// 	// if len(req.PersonIDs) > 0 {
// 	// 	if err := s.lifeGroupRepo.UpdatePersons(lifeGroup.ID, req.PersonIDs); err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// }

// 	// Get updated life group with all relations
// 	updatedLifeGroup, err := s.lifeGroupRepo.GetByID(lifeGroup.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.toResponse(updatedLifeGroup), nil
// }

// func (s *lifeGroupService) GetAll() ([]dto.LifeGroupResponse, error) {
// 	lifeGroups, err := s.lifeGroupRepo.GetAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []dto.LifeGroupResponse
// 	for _, lifeGroup := range lifeGroups {
// 		responses = append(responses, *s.toResponse(&lifeGroup))
// 	}

// 	return responses, nil
// }

// func (s *lifeGroupService) GetByID(id uuid.UUID) (*dto.LifeGroupResponse, error) {
// 	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.toResponse(lifeGroup), nil
// }

// func (s *lifeGroupService) Update(id uuid.UUID, req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error) {
// 	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lifeGroup.Name = req.Name
// 	lifeGroup.Location = req.Location
// 	lifeGroup.WhatsAppLink = req.WhatsAppLink
// 	lifeGroup.ChurchID = req.ChurchID
// 	lifeGroup.LeaderID = req.LeaderID

// 	if err := s.lifeGroupRepo.Update(lifeGroup); err != nil {
// 		return nil, err
// 	}

// 	// Update members and persons if provided
// 	if len(req.MemberIDs) > 0 {
// 		if err := s.lifeGroupRepo.UpdateMembers(id, req.MemberIDs); err != nil {
// 			return nil, err
// 		}
// 	}

// 	// if len(req.PersonIDs) > 0 {
// 	// 	if err := s.lifeGroupRepo.UpdatePersons(id, req.PersonIDs); err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// }

// 	// Get updated life group with all relations
// 	updatedLifeGroup, err := s.lifeGroupRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.toResponse(updatedLifeGroup), nil
// }

// func (s *lifeGroupService) Delete(id uuid.UUID) error {
// 	return s.lifeGroupRepo.Delete(id)
// }

// func (s *lifeGroupService) UpdateLeader(id uuid.UUID, req *dto.UpdateLeaderRequest) (*dto.LifeGroupResponse, error) {
// 	if err := s.lifeGroupRepo.UpdateLeader(id, req.LeaderID); err != nil {
// 		return nil, err
// 	}

// 	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.toResponse(lifeGroup), nil
// }

// func (s *lifeGroupService) UpdateMembers(id uuid.UUID, req *dto.UpdateMembersRequest) (*dto.LifeGroupResponse, error) {
// 	if err := s.lifeGroupRepo.UpdateMembers(id, req.MemberIDs); err != nil {
// 		return nil, err
// 	}

// 	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.toResponse(lifeGroup), nil
// }

// // func (s *lifeGroupService) UpdatePersons(id uuid.UUID, req *dto.UpdatePersonsRequest) (*dto.LifeGroupResponse, error) {
// // 	if err := s.lifeGroupRepo.UpdatePersons(id, req.PersonIDs); err != nil {
// // 		return nil, err
// // 	}

// // 	lifeGroup, err := s.lifeGroupRepo.GetByID(id)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return s.toResponse(lifeGroup), nil
// // }

// func (s *lifeGroupService) GetByChurchID(churchID uuid.UUID) ([]dto.LifeGroupResponse, error) {
// 	lifeGroups, err := s.lifeGroupRepo.GetByChurchID(churchID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responses []dto.LifeGroupResponse
// 	for _, lifeGroup := range lifeGroups {
// 		responses = append(responses, *s.toResponse(&lifeGroup))
// 	}

// 	return responses, nil
// }

// func (s *lifeGroupService) AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
// 	return s.personRepository.AddToLifeGroup(ctx, personID, lifeGroupID)
// }

// func (s *lifeGroupService) RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
// 	return s.personRepository.RemoveFromLifeGroup(ctx, personID, lifeGroupID)
// }

// func (s *lifeGroupService) toResponse(lifeGroup *entity.LifeGroup) *dto.LifeGroupResponse {
// 	return &dto.LifeGroupResponse{
// 		ID:           lifeGroup.ID,
// 		Name:         lifeGroup.Name,
// 		Location:     lifeGroup.Location,
// 		WhatsAppLink: lifeGroup.WhatsAppLink,
// 		ChurchID:     lifeGroup.ChurchID,
// 		Church:       lifeGroup.Church,
// 		LeaderID:     lifeGroup.LeaderID,
// 		Leader:       lifeGroup.Leader,
// 		Persons:      lifeGroup.Persons,
// 		CreatedAt:    lifeGroup.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt:    lifeGroup.UpdatedAt.Format("2006-01-02 15:04:05"),
// 	}
// }
