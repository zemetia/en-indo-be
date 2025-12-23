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

type MusikService interface {
	GetPelayanMusik(ctx context.Context, userID uuid.UUID) ([]dto.MusicianResponse, error)
	GetPelayanMusikByID(ctx context.Context, userID uuid.UUID, personID uuid.UUID) (*dto.MusicianDetailResponse, error)
	AssignPelayananToMusician(ctx context.Context, userID uuid.UUID, req *dto.AssignPelayananRequest) error
	RemovePelayananFromMusician(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID) error
	ToggleActiveStatus(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID, isActive bool) error
	GetAvailableMusikPelayanan(ctx context.Context) ([]dto.PelayananRoleResponse, error)
	GetAvailablePeople(ctx context.Context, userID uuid.UUID) ([]dto.AvailablePersonResponse, error)
}

type musikService struct {
	pelayananRepo  repository.PelayananRepository
	departmentRepo repository.DepartmentRepository
	userRepo       repository.UserRepository
	personRepo     repository.PersonRepository
}

func NewMusikService(
	pelayananRepo repository.PelayananRepository,
	departmentRepo repository.DepartmentRepository,
	userRepo repository.UserRepository,
	personRepo repository.PersonRepository,
) MusikService {
	return &musikService{
		pelayananRepo:  pelayananRepo,
		departmentRepo: departmentRepo,
		userRepo:       userRepo,
		personRepo:     personRepo,
	}
}

// GetPelayanMusik gets all musicians based on PIC permissions
// If user is PIC Musik: returns all musicians from their assigned churches
// If user is not PIC Musik: returns empty list (only PICs can manage musicians)
func (s *musikService) GetPelayanMusik(ctx context.Context, userID uuid.UUID) ([]dto.MusicianResponse, error) {
	// Get the music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return nil, fmt.Errorf("failed to get music department: %v", err)
	}

	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Get all pelayanan assignments for this person
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user assignments: %v", err)
	}

	// Check if user has PIC Musik role and collect their church IDs
	picChurchIDs := make(map[uuid.UUID]bool)
	for _, assignment := range assignments {
		// Check if pelayanan is in Music department and is PIC role
		if assignment.Pelayanan.IsPic &&
			assignment.Pelayanan.DepartmentID == musikDept.ID &&
			(assignment.Pelayanan.Pelayanan == "PIC Musik" ||
				(strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "pic") &&
					strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "musik"))) {
			picChurchIDs[assignment.ChurchID] = true
		}
	}

	// If user has no PIC Musik role, return empty list
	if len(picChurchIDs) == 0 {
		return []dto.MusicianResponse{}, nil
	}

	// Convert map keys to slice
	churchIDSlice := make([]uuid.UUID, 0, len(picChurchIDs))
	for churchID := range picChurchIDs {
		churchIDSlice = append(churchIDSlice, churchID)
	}

	// Get all music assignments from churches where user has PIC role
	musicAssignments, err := s.pelayananRepo.GetAssignmentsByDepartmentAndChurches(ctx, musikDept.ID, churchIDSlice)
	if err != nil {
		return nil, fmt.Errorf("failed to get music assignments: %v", err)
	}

	// Group assignments by person to create musician responses
	musicianMap := make(map[uuid.UUID]*dto.MusicianResponse)
	for _, assignment := range musicAssignments {
		personID := assignment.PersonID

		// Initialize musician response if not exists
		if _, exists := musicianMap[personID]; !exists {
			musicianMap[personID] = &dto.MusicianResponse{
				ID:          assignment.Person.ID.String(),
				Nama:        assignment.Person.Nama,
				Email:       assignment.Person.Email,
				Telepon:     assignment.Person.NomorTelepon,
				Instruments: []string{},
				Status:      "active", // Will be set to inactive if any assignment is inactive
				Avatar:      "https://placehold.co/100x100.png", // Default avatar
			}
		}

		// Add pelayanan to instruments list only if NOT a PIC role
		if !assignment.Pelayanan.IsPic {
			musicianMap[personID].Instruments = append(musicianMap[personID].Instruments, assignment.Pelayanan.Pelayanan)
		}

		// If any assignment is inactive, set status to inactive
		if !assignment.IsActive {
			musicianMap[personID].Status = "inactive"
		}
	}

	// Convert map to slice, excluding entries with no instruments (PIC only)
	musicians := make([]dto.MusicianResponse, 0, len(musicianMap))
	for _, musician := range musicianMap {
		if len(musician.Instruments) > 0 {
			musicians = append(musicians, *musician)
		}
	}

	return musicians, nil
}

// GetPelayanMusikByID gets detailed information about a specific musician
func (s *musikService) GetPelayanMusikByID(ctx context.Context, userID uuid.UUID, personID uuid.UUID) (*dto.MusicianDetailResponse, error) {
	// Get the music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return nil, fmt.Errorf("failed to get music department: %v", err)
	}

	// Verify user has PIC permission (similar to GetPelayanMusik)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user assignments: %v", err)
	}

	// Check PIC permission and collect church IDs
	picChurchIDs := make(map[uuid.UUID]bool)
	for _, assignment := range assignments {
		if assignment.Pelayanan.IsPic &&
			assignment.Pelayanan.DepartmentID == musikDept.ID {
			picChurchIDs[assignment.ChurchID] = true
		}
	}

	if len(picChurchIDs) == 0 {
		return nil, fmt.Errorf("user does not have PIC Musik permission")
	}

	// Get person details
	person, err := s.personRepo.GetByID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %v", err)
	}

	// Verify person belongs to a church where user has PIC permission
	if !picChurchIDs[person.ChurchID] {
		return nil, fmt.Errorf("user does not have permission for this person's church")
	}

	// Get all music assignments for this person
	musicAssignments, err := s.pelayananRepo.GetAssignmentsByPersonAndDepartment(ctx, personID, musikDept.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get music assignments: %v", err)
	}

	// Build response
	response := &dto.MusicianDetailResponse{
		ID:         person.ID.String(),
		Nama:       person.Nama,
		Email:      person.Email,
		Telepon:    person.NomorTelepon,
		Avatar:     "https://placehold.co/100x100.png",
		ChurchID:   person.ChurchID.String(),
		ChurchName: person.Church.Name,
		Pelayanan:  []dto.MusicianPelayananDetailResponse{},
	}

	for _, assignment := range musicAssignments {
		response.Pelayanan = append(response.Pelayanan, dto.MusicianPelayananDetailResponse{
			AssignmentID: assignment.ID.String(),
			PelayananID:  assignment.PelayananID.String(),
			Pelayanan:    assignment.Pelayanan.Pelayanan,
			ChurchID:     assignment.ChurchID.String(),
			ChurchName:   assignment.Church.Name,
			IsActive:     assignment.IsActive,
			IsPic:        assignment.Pelayanan.IsPic,
		})
	}

	return response, nil
}

// AssignPelayananToMusician assigns a pelayanan role to a musician
func (s *musikService) AssignPelayananToMusician(ctx context.Context, userID uuid.UUID, req *dto.AssignPelayananRequest) error {
	// Verify the pelayanan is in music department
	pelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, req.PelayananID)
	if err != nil {
		return fmt.Errorf("failed to get pelayanan: %v", err)
	}

	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return fmt.Errorf("failed to get music department: %v", err)
	}

	if pelayanan.DepartmentID != musikDept.ID {
		return fmt.Errorf("pelayanan is not in music department")
	}

	// Verify person exists and belongs to the specified church
	person, err := s.personRepo.GetByID(ctx, req.PersonID)
	if err != nil {
		return fmt.Errorf("failed to get person: %v", err)
	}

	if person.ChurchID != req.ChurchID {
		return fmt.Errorf("person does not belong to specified church")
	}

	// Verify user has PIC permission
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return fmt.Errorf("failed to get user assignments: %v", err)
	}

	isPICMusik := false
	hasChurchPermission := false
	for _, assignment := range assignments {
		if assignment.Pelayanan.IsPic &&
			assignment.Pelayanan.DepartmentID == musikDept.ID {
			isPICMusik = true
			if assignment.ChurchID == req.ChurchID {
				hasChurchPermission = true
			}
		}
	}

	if !isPICMusik {
		return fmt.Errorf("user does not have PIC Musik permission")
	}

	if !hasChurchPermission {
		return fmt.Errorf("user does not have permission for this church")
	}

	// Check if assignment already exists
	existingAssignment, err := s.pelayananRepo.GetAssignmentByPersonPelayananChurch(ctx, req.PersonID, req.PelayananID, req.ChurchID)
	if err == nil && existingAssignment != nil {
		return fmt.Errorf("person already has this pelayanan assignment")
	}

	// Create new assignment
	newAssignment := &entity.PersonPelayananGereja{
		ID:          uuid.New(),
		PersonID:    req.PersonID,
		PelayananID: req.PelayananID,
		ChurchID:    req.ChurchID,
		IsActive:    true,
	}

	if err := s.pelayananRepo.CreatePelayananAssignment(ctx, newAssignment); err != nil {
		return fmt.Errorf("failed to create assignment: %v", err)
	}

	return nil
}

// RemovePelayananFromMusician removes a pelayanan assignment from a musician
func (s *musikService) RemovePelayananFromMusician(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID) error {
	// Get the assignment
	assignment, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %v", err)
	}

	// Verify the assignment is in music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return fmt.Errorf("failed to get music department: %v", err)
	}

	if assignment.Pelayanan.DepartmentID != musikDept.ID {
		return fmt.Errorf("assignment is not in music department")
	}

	// Verify user has PIC permission for this church
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return fmt.Errorf("failed to get user assignments: %v", err)
	}

	isPICMusik := false
	hasChurchPermission := false
	for _, userAssignment := range assignments {
		if userAssignment.Pelayanan.IsPic &&
			userAssignment.Pelayanan.DepartmentID == musikDept.ID {
			isPICMusik = true
			if userAssignment.ChurchID == assignment.ChurchID {
				hasChurchPermission = true
			}
		}
	}

	if !isPICMusik {
		return fmt.Errorf("user does not have PIC Musik permission")
	}

	if !hasChurchPermission {
		return fmt.Errorf("user does not have permission for this church")
	}

	// Delete the assignment
	if err := s.pelayananRepo.DeletePelayananAssignment(ctx, assignmentID); err != nil {
		return fmt.Errorf("failed to delete assignment: %v", err)
	}

	return nil
}

// ToggleActiveStatus toggles the active status of a pelayanan assignment
func (s *musikService) ToggleActiveStatus(ctx context.Context, userID uuid.UUID, assignmentID uuid.UUID, isActive bool) error {
	// Get the assignment
	assignment, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %v", err)
	}

	// Verify the assignment is in music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return fmt.Errorf("failed to get music department: %v", err)
	}

	if assignment.Pelayanan.DepartmentID != musikDept.ID {
		return fmt.Errorf("assignment is not in music department")
	}

	// Verify user has PIC permission for this church
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return fmt.Errorf("failed to get user assignments: %v", err)
	}

	isPICMusik := false
	hasChurchPermission := false
	for _, userAssignment := range assignments {
		if userAssignment.Pelayanan.IsPic &&
			userAssignment.Pelayanan.DepartmentID == musikDept.ID {
			isPICMusik = true
			if userAssignment.ChurchID == assignment.ChurchID {
				hasChurchPermission = true
			}
		}
	}

	if !isPICMusik {
		return fmt.Errorf("user does not have PIC Musik permission")
	}

	if !hasChurchPermission {
		return fmt.Errorf("user does not have permission for this church")
	}

	// Update the assignment
	assignment.IsActive = isActive
	if err := s.pelayananRepo.UpdatePelayananAssignment(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %v", err)
	}

	return nil
}

// GetAvailableMusikPelayanan gets all available pelayanan roles in the music department
func (s *musikService) GetAvailableMusikPelayanan(ctx context.Context) ([]dto.PelayananRoleResponse, error) {
	// Get the music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return nil, fmt.Errorf("failed to get music department: %v", err)
	}

	// Get all pelayanan in music department
	pelayananList, err := s.pelayananRepo.GetAllPelayananByDepartment(ctx, musikDept.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pelayanan list: %v", err)
	}

	// Convert to response DTOs
	responses := make([]dto.PelayananRoleResponse, 0, len(pelayananList))
	for _, pelayanan := range pelayananList {
		responses = append(responses, dto.PelayananRoleResponse{
			ID:          pelayanan.ID.String(),
			Pelayanan:   pelayanan.Pelayanan,
			Description: pelayanan.Description,
			IsPic:       pelayanan.IsPic,
		})
	}

	return responses, nil
}

// GetAvailablePeople gets all people available to be added to music ministry
// Returns people from churches where user is PIC Musik who don't have any music assignments yet
func (s *musikService) GetAvailablePeople(ctx context.Context, userID uuid.UUID) ([]dto.AvailablePersonResponse, error) {
	// Get the music department
	musikDept, err := s.departmentRepo.GetByName("Musik")
	if err != nil {
		return nil, fmt.Errorf("failed to get music department: %v", err)
	}

	// Get user to find their PersonID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Get all pelayanan assignments for this person
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, user.PersonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user assignments: %v", err)
	}

	// Check if user has PIC Musik role and collect their church IDs
	picChurchIDs := make(map[uuid.UUID]bool)
	for _, assignment := range assignments {
		// Check if pelayanan is in Music department and is PIC role
		if assignment.Pelayanan.IsPic &&
			assignment.Pelayanan.DepartmentID == musikDept.ID &&
			(assignment.Pelayanan.Pelayanan == "PIC Musik" ||
				(strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "pic") &&
					strings.Contains(strings.ToLower(assignment.Pelayanan.Pelayanan), "musik"))) {
			picChurchIDs[assignment.ChurchID] = true
		}
	}

	// If user has no PIC Musik role, return empty list
	if len(picChurchIDs) == 0 {
		return []dto.AvailablePersonResponse{}, nil
	}

	// Convert map keys to slice
	churchIDSlice := make([]uuid.UUID, 0, len(picChurchIDs))
	for churchID := range picChurchIDs {
		churchIDSlice = append(churchIDSlice, churchID)
	}

	// Get people who don't have any music department assignments in these churches
	availablePeople, err := s.personRepo.GetPeopleWithoutDepartmentInChurches(ctx, musikDept.ID, churchIDSlice)
	if err != nil {
		return nil, fmt.Errorf("failed to get available people: %v", err)
	}

	// Convert to response DTOs
	responses := make([]dto.AvailablePersonResponse, 0, len(availablePeople))
	for _, person := range availablePeople {
		responses = append(responses, dto.AvailablePersonResponse{
			ID:         person.ID.String(),
			Nama:       person.Nama,
			Email:      person.Email,
			Telepon:    person.NomorTelepon,
			ChurchID:   person.ChurchID.String(),
			ChurchName: person.Church.Name,
			Avatar:     "https://placehold.co/100x100.png",
		})
	}

	return responses, nil
}
