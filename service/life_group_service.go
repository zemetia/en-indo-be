package service

import (
	"context"
	"fmt"
	"strings"

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
	GetMyLifeGroup(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error)
	GetDaftarLifeGroup(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error)
	CheckUserPermission(userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	CheckUserPICPermission(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	CheckUserCanManageLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	CheckUserCanEditLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	CheckUserCanDeleteLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	CheckUserCanViewLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error)
	GetByMultipleChurchIDs(churchIDs []uuid.UUID) ([]dto.BatchChurchLifeGroupsResponse, error)
	GetLifeGroupsByPICRole(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error)
}

type lifeGroupService struct {
	lifeGroupRepo    repository.LifeGroupRepository
	pelayananRepo    repository.PelayananRepository
	userRepo         repository.UserRepository
	personRepo       repository.PersonRepository
	personMemberRepo repository.LifeGroupPersonMemberRepository
}

func NewLifeGroupService(lifeGroupRepo repository.LifeGroupRepository, pelayananRepo repository.PelayananRepository, userRepo repository.UserRepository, personRepo repository.PersonRepository, personMemberRepo repository.LifeGroupPersonMemberRepository) LifeGroupService {
	return &lifeGroupService{
		lifeGroupRepo:    lifeGroupRepo,
		pelayananRepo:    pelayananRepo,
		userRepo:         userRepo,
		personRepo:       personRepo,
		personMemberRepo: personMemberRepo,
	}
}

func (s *lifeGroupService) Create(req *dto.LifeGroupRequest) (*dto.LifeGroupResponse, error) {
	// Parse string UUIDs to uuid.UUID
	churchID, err := uuid.Parse(req.ChurchID)
	if err != nil {
		return nil, fmt.Errorf("invalid church_id format: %v", err)
	}

	// Create LifeGroup without any leader references
	lifeGroup := &entity.LifeGroup{
		ID:           uuid.New(),
		Name:         req.Name,
		Location:     req.Location,
		WhatsAppLink: req.WhatsAppLink,
		ChurchID:     churchID,
	}

	if err := s.lifeGroupRepo.Create(lifeGroup); err != nil {
		return nil, err
	}

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

	// Parse string UUIDs to uuid.UUID
	churchID, err := uuid.Parse(req.ChurchID)
	if err != nil {
		return nil, fmt.Errorf("invalid church_id format: %v", err)
	}

	lifeGroup.Name = req.Name
	lifeGroup.Location = req.Location
	lifeGroup.WhatsAppLink = req.WhatsAppLink
	lifeGroup.ChurchID = churchID

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
	// Leaders are now managed through PersonMember API
	// This method is deprecated - use PersonMember endpoints instead
	return nil, fmt.Errorf("leader management has been moved to PersonMember API endpoints")
}

func (s *lifeGroupService) toResponse(lifeGroup *entity.LifeGroup) *dto.LifeGroupResponse {
	response := &dto.LifeGroupResponse{
		ID:             lifeGroup.ID,
		Name:           lifeGroup.Name,
		Location:       lifeGroup.Location,
		WhatsAppLink:   lifeGroup.WhatsAppLink,
		ChurchID:       lifeGroup.ChurchID,
		Church:         lifeGroup.Church,
		PersonMembers:  lifeGroup.PersonMembers,
		VisitorMembers: lifeGroup.VisitorMembers,
		CreatedAt:      lifeGroup.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      lifeGroup.UpdatedAt.Format("2006-01-02 15:04:05"),
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
	ctx := context.Background()

	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all pelayanan assignments for this person to check if they're PIC Lifegroup
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return nil, err
	}

	var allLifeGroups []entity.LifeGroup
	isLifeGroupPIC := false
	picChurchIDs := make(map[uuid.UUID]bool)

	// Check if user is PIC Lifegroup and collect their church IDs
	for _, assignment := range assignments {
		if assignment.Pelayanan.IsPic &&
			(assignment.Pelayanan.Pelayanan == "PIC Lifegroup" ||
				(strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "pic") &&
					strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "lifegroup"))) {
			isLifeGroupPIC = true
			picChurchIDs[assignment.ChurchID] = true
		}
	}

	if isLifeGroupPIC {
		// If user is PIC Lifegroup, get all lifegroups from their assigned churches
		for churchID := range picChurchIDs {
			churchLifeGroups, err := s.lifeGroupRepo.GetByChurchID(churchID)
			if err != nil {
				continue // Skip this church if error, don't fail entire request
			}
			allLifeGroups = append(allLifeGroups, churchLifeGroups...)
		}
	} else {
		// If not PIC, get lifegroups where user is leader/co-leader
		leaderLifeGroups, err := s.lifeGroupRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		allLifeGroups = append(allLifeGroups, leaderLifeGroups...)

		// Also get lifegroups where user's person is a member
		memberLifeGroups, err := s.personMemberRepo.GetByPersonID(ctx, user.PersonID)
		if err != nil {
			// If error getting member lifegroups, just continue with leader ones
		} else {
			// For each membership, get the lifegroup if not already included
			existingIDs := make(map[uuid.UUID]bool)
			for _, lg := range allLifeGroups {
				existingIDs[lg.ID] = true
			}

			for _, member := range memberLifeGroups {
				if member.IsActive && !existingIDs[member.LifeGroupID] {
					lifeGroup, err := s.lifeGroupRepo.GetByID(member.LifeGroupID)
					if err == nil {
						allLifeGroups = append(allLifeGroups, *lifeGroup)
						existingIDs[member.LifeGroupID] = true
					}
				}
			}
		}
	}

	// Remove duplicates and convert to responses
	uniqueLifeGroups := make(map[uuid.UUID]entity.LifeGroup)
	for _, lifeGroup := range allLifeGroups {
		uniqueLifeGroups[lifeGroup.ID] = lifeGroup
	}

	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range uniqueLifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

// GetMyLifeGroup returns lifegroups where the user is an active member in LifeGroupPersonMember
// This only shows lifegroups based on explicit membership, ignoring PIC privileges
func (s *lifeGroupService) GetMyLifeGroup(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error) {
	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get lifegroups where user's person is a member only
	memberLifeGroups, err := s.personMemberRepo.GetByPersonID(ctx, user.PersonID)
	if err != nil {
		return []dto.LifeGroupResponse{}, nil // Return empty slice if error
	}

	var allLifeGroups []entity.LifeGroup
	for _, member := range memberLifeGroups {
		if member.IsActive {
			lifeGroup, err := s.lifeGroupRepo.GetByID(member.LifeGroupID)
			if err == nil {
				allLifeGroups = append(allLifeGroups, *lifeGroup)
			}
		}
	}

	// Convert to responses
	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range allLifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}

// GetDaftarLifeGroup returns all lifegroups for PIC Lifegroup users only
// - PIC Lifegroup: all lifegroups from their assigned churches
// - Non-PIC users: access denied (empty result)
func (s *lifeGroupService) GetDaftarLifeGroup(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error) {
	// Use the new PIC role logic to get lifegroups
	return s.GetLifeGroupsByPICRole(ctx, userID)
}

func (s *lifeGroupService) CheckUserPermission(userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	ctx := context.Background()
	
	// Check if user has a person record and if that person is a leader or co-leader
	// First get user's person ID
	person, err := s.personRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil // User doesn't have a person record
	}
	
	// Check if person is a leader or co-leader in PersonMembers
	member, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, person.ID)
	if err != nil {
		return false, nil // Person is not a member of this lifegroup
	}
	
	// Check if member has leadership role
	return member.Position == entity.PersonMemberPositionLeader || 
		   member.Position == entity.PersonMemberPositionCoLeader, nil
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

// CheckUserPICPermission checks if user is a PIC Lifegroup for the church containing this lifegroup
func (s *lifeGroupService) CheckUserPICPermission(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// Get lifegroup to find its church
	lifeGroup, err := s.lifeGroupRepo.GetByID(lifeGroupID)
	if err != nil {
		return false, err
	}

	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Get all pelayanan assignments for this person
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return false, err
	}

	// Check if user has PIC Lifegroup role in the same church as the lifegroup
	for _, assignment := range assignments {
		// Check if this is a PIC Lifegroup assignment for the same church
		if assignment.Pelayanan.IsPic &&
			(assignment.Pelayanan.Pelayanan == "PIC Lifegroup" ||
				(strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "pic") &&
					strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "lifegroup"))) &&
			assignment.ChurchID == lifeGroup.ChurchID {
			return true, nil
		}
	}

	return false, nil
}

// CheckUserCanManageLifeGroup checks if user can manage (edit/delete/add members) this lifegroup
// Returns true if user is PIC Lifegroup for the church OR leader/co-leader of this specific lifegroup
func (s *lifeGroupService) CheckUserCanManageLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// First check if user is PIC Lifegroup for this church
	isPIC, err := s.CheckUserPICPermission(ctx, userID, lifeGroupID)
	if err != nil {
		return false, err
	}
	if isPIC {
		return true, nil
	}

	// If not PIC, check if user is leader or co-leader of this specific lifegroup
	isLeaderOrCoLeader, err := s.CheckUserPermission(userID, lifeGroupID)
	if err != nil {
		return false, err
	}

	return isLeaderOrCoLeader, nil
}

// CheckUserCanEditLifeGroup checks if user can edit this lifegroup
// Returns true if user is PIC Lifegroup for the church OR leader OR co-leader of this specific lifegroup
func (s *lifeGroupService) CheckUserCanEditLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// First check if user is PIC Lifegroup for this church
	isPIC, err := s.CheckUserPICPermission(ctx, userID, lifeGroupID)
	if err != nil {
		return false, err
	}
	if isPIC {
		return true, nil
	}

	// If not PIC, check if user is leader or co-leader of this specific lifegroup
	isLeaderOrCoLeader, err := s.CheckUserPermission(userID, lifeGroupID)
	if err != nil {
		return false, err
	}

	return isLeaderOrCoLeader, nil
}

// CheckUserCanDeleteLifeGroup checks if user can delete this lifegroup
// Returns true if user is PIC Lifegroup for the church OR leader (NOT co-leader) of this specific lifegroup
func (s *lifeGroupService) CheckUserCanDeleteLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// First check if user is PIC Lifegroup for this church
	isPIC, err := s.CheckUserPICPermission(ctx, userID, lifeGroupID)
	if err != nil {
		return false, err
	}
	if isPIC {
		return true, nil
	}

	// If not PIC, check if user is leader (NOT co-leader) of this specific lifegroup
	// First get user's person ID
	person, err := s.personRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil // User doesn't have a person record
	}
	
	// Check if person is a leader (NOT co-leader) in PersonMembers
	member, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, person.ID)
	if err != nil {
		return false, nil // Person is not a member of this lifegroup
	}
	
	// Only allow leader, not co-leader
	return member.Position == entity.PersonMemberPositionLeader, nil
}

// CheckUserCanViewLifeGroup checks if user can view this lifegroup
// Returns true if user is a person member in LifeGroupPersonMember OR has PIC Lifegroup pelayanan in same church
func (s *lifeGroupService) CheckUserCanViewLifeGroup(ctx context.Context, userID uuid.UUID, lifeGroupID uuid.UUID) (bool, error) {
	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if user's person is an active member of this lifegroup
	member, err := s.personMemberRepo.GetByLifeGroupAndPersonID(ctx, lifeGroupID, user.PersonID)
	if err == nil && member != nil && member.IsActive {
		return true, nil
	}

	// Check if user has PIC Lifegroup pelayanan for the lifegroup's church
	isPIC, err := s.CheckUserPICPermission(ctx, userID, lifeGroupID)
	if err != nil {
		return false, err
	}
	if isPIC {
		return true, nil
	}

	// Access denied - user is neither a lifegroup member nor PIC Lifegroup for this church
	return false, nil
}

// GetLifeGroupsByPICRole returns all lifegroups from churches where the user has PIC Lifegroup role
func (s *lifeGroupService) GetLifeGroupsByPICRole(ctx context.Context, userID uuid.UUID) ([]dto.LifeGroupResponse, error) {
	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all pelayanan assignments for this person
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return nil, err
	}

	// Check if user has PIC Lifegroup role and collect their church IDs
	picChurchIDs := make(map[uuid.UUID]bool)
	for _, assignment := range assignments {
		if assignment.Pelayanan.IsPic &&
			(assignment.Pelayanan.Pelayanan == "PIC Lifegroup" ||
				(strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "pic") &&
					strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "lifegroup"))) {
			picChurchIDs[assignment.ChurchID] = true
		}
	}

	// If user has no PIC Lifegroup role, return empty slice
	if len(picChurchIDs) == 0 {
		return []dto.LifeGroupResponse{}, nil
	}

	// Get all lifegroups from churches where user has PIC Lifegroup role
	var allLifeGroups []entity.LifeGroup
	for churchID := range picChurchIDs {
		churchLifeGroups, err := s.lifeGroupRepo.GetByChurchID(churchID)
		if err != nil {
			continue // Skip this church if error, don't fail entire request
		}
		allLifeGroups = append(allLifeGroups, churchLifeGroups...)
	}

	// Convert to responses
	var responses []dto.LifeGroupResponse
	for _, lifeGroup := range allLifeGroups {
		responses = append(responses, *s.toResponse(&lifeGroup))
	}

	return responses, nil
}
