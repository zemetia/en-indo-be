package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
	"gorm.io/gorm"
)

type PelayananService interface {
	// Pelayanan entity CRUD
	CreatePelayanan(ctx context.Context, request dto.PelayananRequest) (dto.PelayananResponse, error)
	UpdatePelayanan(ctx context.Context, id uuid.UUID, request dto.UpdatePelayananRequest) (dto.PelayananResponse, error)
	DeletePelayanan(ctx context.Context, id uuid.UUID) error
	GetPelayananByID(ctx context.Context, id uuid.UUID) (dto.PelayananResponse, error)
	GetAllPelayanan(ctx context.Context, departmentID string) ([]dto.PelayananResponse, error)

	// Assignment operations
	GetMyPelayanan(ctx context.Context, personID uuid.UUID) ([]dto.PersonHasPelayananResponse, error)
	GetAllAssignments(ctx context.Context, req dto.PaginationRequest) (dto.PelayananAssignmentPaginationResponse, error)
	AssignPelayanan(ctx context.Context, request dto.AssignPelayananRequest) error
	UnassignPelayanan(ctx context.Context, assignmentID uuid.UUID) error
	GetAssignmentByID(ctx context.Context, assignmentID uuid.UUID) (dto.PelayananAssignmentResponse, error)
}

type pelayananService struct {
	pelayananRepo  repository.PelayananRepository
	personRepo     repository.PersonRepository
	churchRepo     repository.ChurchRepository
	departmentRepo repository.DepartmentRepository
	userService    UserService
}

func NewPelayananService(
	pelayananRepo repository.PelayananRepository,
	personRepo repository.PersonRepository,
	churchRepo repository.ChurchRepository,
	departmentRepo repository.DepartmentRepository,
	userService UserService,
) PelayananService {
	return &pelayananService{
		pelayananRepo:  pelayananRepo,
		personRepo:     personRepo,
		churchRepo:     churchRepo,
		departmentRepo: departmentRepo,
		userService:    userService,
	}
}

// Pelayanan entity CRUD methods
func (s *pelayananService) CreatePelayanan(ctx context.Context, request dto.PelayananRequest) (dto.PelayananResponse, error) {
	// Validate department exists
	_, err := s.departmentRepo.GetByID(request.DepartmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.PelayananResponse{}, fmt.Errorf("department with ID %s not found", request.DepartmentID)
		}
		return dto.PelayananResponse{}, fmt.Errorf("failed to validate department: %w", err)
	}

	pelayanan := &entity.Pelayanan{
		ID:           uuid.New(),
		Pelayanan:    request.Pelayanan,
		Description:  request.Description,
		DepartmentID: request.DepartmentID,
	}

	if err := s.pelayananRepo.CreatePelayanan(ctx, pelayanan); err != nil {
		return dto.PelayananResponse{}, fmt.Errorf("failed to create pelayanan: %w", err)
	}

	// Get the created pelayanan with preloaded department
	createdPelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, pelayanan.ID)
	if err != nil {
		return dto.PelayananResponse{}, fmt.Errorf("failed to get created pelayanan: %w", err)
	}

	return dto.PelayananResponse{
		ID:          createdPelayanan.ID,
		Pelayanan:   createdPelayanan.Pelayanan,
		Description: createdPelayanan.Description,
		Department: dto.DepartmentResponse{
			ID:          createdPelayanan.Department.ID,
			Name:        createdPelayanan.Department.Name,
			Description: createdPelayanan.Department.Description,
			CreatedAt:   createdPelayanan.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   createdPelayanan.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		CreatedAt: createdPelayanan.CreatedAt,
		UpdatedAt: createdPelayanan.UpdatedAt,
	}, nil
}

func (s *pelayananService) UpdatePelayanan(ctx context.Context, id uuid.UUID, request dto.UpdatePelayananRequest) (dto.PelayananResponse, error) {
	// Get existing pelayanan
	pelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.PelayananResponse{}, fmt.Errorf("pelayanan with ID %s not found", id)
		}
		return dto.PelayananResponse{}, fmt.Errorf("failed to get pelayanan: %w", err)
	}

	// Update fields if provided
	if request.Pelayanan != "" {
		pelayanan.Pelayanan = request.Pelayanan
	}
	if request.Description != "" {
		pelayanan.Description = request.Description
	}
	if request.DepartmentID != uuid.Nil {
		// Validate department exists
		_, err := s.departmentRepo.GetByID(request.DepartmentID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.PelayananResponse{}, fmt.Errorf("department with ID %s not found", request.DepartmentID)
			}
			return dto.PelayananResponse{}, fmt.Errorf("failed to validate department: %w", err)
		}
		pelayanan.DepartmentID = request.DepartmentID
	}

	if err := s.pelayananRepo.UpdatePelayanan(ctx, pelayanan); err != nil {
		return dto.PelayananResponse{}, fmt.Errorf("failed to update pelayanan: %w", err)
	}

	// Get updated pelayanan with preloaded department
	updatedPelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, id)
	if err != nil {
		return dto.PelayananResponse{}, fmt.Errorf("failed to get updated pelayanan: %w", err)
	}

	return dto.PelayananResponse{
		ID:          updatedPelayanan.ID,
		Pelayanan:   updatedPelayanan.Pelayanan,
		Description: updatedPelayanan.Description,
		Department: dto.DepartmentResponse{
			ID:          updatedPelayanan.Department.ID,
			Name:        updatedPelayanan.Department.Name,
			Description: updatedPelayanan.Department.Description,
			CreatedAt:   updatedPelayanan.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   updatedPelayanan.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		CreatedAt: updatedPelayanan.CreatedAt,
		UpdatedAt: updatedPelayanan.UpdatedAt,
	}, nil
}

func (s *pelayananService) DeletePelayanan(ctx context.Context, id uuid.UUID) error {
	// Check if pelayanan exists
	pelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("pelayanan with ID %s not found", id)
		}
		return fmt.Errorf("failed to get pelayanan: %w", err)
	}

	// Prevent deletion of PIC pelayanan
	if pelayanan.IsPic {
		return fmt.Errorf("cannot delete PIC pelayanan. PIC pelayanan can only be deleted when the department is deleted")
	}

	if err := s.pelayananRepo.DeletePelayanan(ctx, id); err != nil {
		return fmt.Errorf("failed to delete pelayanan: %w", err)
	}

	return nil
}

func (s *pelayananService) GetPelayananByID(ctx context.Context, id uuid.UUID) (dto.PelayananResponse, error) {
	pelayanan, err := s.pelayananRepo.GetPelayananByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.PelayananResponse{}, fmt.Errorf("pelayanan with ID %s not found", id)
		}
		return dto.PelayananResponse{}, fmt.Errorf("failed to get pelayanan: %w", err)
	}

	return dto.PelayananResponse{
		ID:          pelayanan.ID,
		Pelayanan:   pelayanan.Pelayanan,
		Description: pelayanan.Description,
		Department: dto.DepartmentResponse{
			ID:          pelayanan.Department.ID,
			Name:        pelayanan.Department.Name,
			Description: pelayanan.Department.Description,
			CreatedAt:   pelayanan.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   pelayanan.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		CreatedAt: pelayanan.CreatedAt,
		UpdatedAt: pelayanan.UpdatedAt,
	}, nil
}

// Assignment operations
func (s *pelayananService) GetMyPelayanan(ctx context.Context, personID uuid.UUID) ([]dto.PersonHasPelayananResponse, error) {
	assignments, err := s.pelayananRepo.GetPelayananByPersonID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pelayanan assignments: %w", err)
	}

	var responses []dto.PersonHasPelayananResponse
	for _, assignment := range assignments {
		responses = append(responses, dto.PersonHasPelayananResponse{
			PelayananID: assignment.PelayananID,
			Pelayanan:   assignment.Pelayanan.Pelayanan,
			ChurchID:    assignment.ChurchID,
			ChurchName:  assignment.Church.Name,
			IsPic:       assignment.Pelayanan.IsPic,
		})
	}

	return responses, nil
}

func (s *pelayananService) GetAllAssignments(ctx context.Context, req dto.PaginationRequest) (dto.PelayananAssignmentPaginationResponse, error) {
	assignments, pagination, err := s.pelayananRepo.GetAllPelayananAssignments(ctx, req)
	if err != nil {
		return dto.PelayananAssignmentPaginationResponse{}, fmt.Errorf("failed to get all assignments: %w", err)
	}

	var responses []dto.PelayananAssignmentResponse
	for _, assignment := range assignments {
		// Check if user account exists and get its status
		hasUserAccount := false
		isUserActive := false
		if user, err := s.userService.GetUserByPersonID(ctx, assignment.PersonID); err == nil {
			hasUserAccount = true
			isUserActive = user.IsActive
		}

		responses = append(responses, dto.PelayananAssignmentResponse{
			ID:             assignment.ID,
			PersonID:       assignment.PersonID,
			PersonName:     assignment.Person.Nama,
			PelayananID:    assignment.PelayananID,
			Pelayanan:      assignment.Pelayanan.Pelayanan,
			ChurchID:       assignment.ChurchID,
			ChurchName:     assignment.Church.Name,
			DepartmentID:   assignment.Pelayanan.Department.ID,
			DepartmentName: assignment.Pelayanan.Department.Name,
			PelayananIsPic: assignment.Pelayanan.IsPic,
			HasUserAccount: hasUserAccount,
			IsUserActive:   isUserActive,
			CreatedAt:      assignment.CreatedAt,
			UpdatedAt:      assignment.UpdatedAt,
		})
	}

	return dto.PelayananAssignmentPaginationResponse{
		Data:               responses,
		PaginationResponse: *pagination,
	}, nil
}

func (s *pelayananService) GetAllPelayanan(ctx context.Context, departmentID string) ([]dto.PelayananResponse, error) {
	var pelayanan []entity.Pelayanan
	var err error

	if departmentID != "" {
		// Parse department ID
		deptID, parseErr := uuid.Parse(departmentID)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid department ID format: %w", parseErr)
		}

		// Get pelayanan filtered by department
		pelayanan, err = s.pelayananRepo.GetAllPelayananByDepartment(ctx, deptID)
	} else {
		// Get all pelayanan
		pelayanan, err = s.pelayananRepo.GetAllPelayanan(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get all pelayanan: %w", err)
	}

	var responses []dto.PelayananResponse
	for _, p := range pelayanan {
		responses = append(responses, dto.PelayananResponse{
			ID:          p.ID,
			Pelayanan:   p.Pelayanan,
			Description: p.Description,
			Department: dto.DepartmentResponse{
				ID:          p.Department.ID,
				Name:        p.Department.Name,
				Description: p.Department.Description,
				CreatedAt:   p.Department.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:   p.Department.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *pelayananService) AssignPelayanan(ctx context.Context, request dto.AssignPelayananRequest) error {
	// Validate that person exists
	person, err := s.personRepo.GetByID(ctx, request.PersonID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("person with ID %s not found", request.PersonID)
		}
		return fmt.Errorf("failed to validate person: %w", err)
	}

	// Validate that pelayanan exists
	_, err = s.pelayananRepo.GetPelayananByID(ctx, request.PelayananID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("pelayanan with ID %s not found", request.PelayananID)
		}
		return fmt.Errorf("failed to validate pelayanan: %w", err)
	}

	// Validate that church exists
	_, err = s.churchRepo.GetByID(request.ChurchID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("church with ID %s not found", request.ChurchID)
		}
		return fmt.Errorf("failed to validate church: %w", err)
	}

	// Check if assignment already exists to prevent duplicates
	existingAssignment, err := s.pelayananRepo.GetAssignmentByPersonPelayananChurch(ctx, request.PersonID, request.PelayananID, request.ChurchID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing assignment: %w", err)
	}
	if existingAssignment != nil {
		return fmt.Errorf("person is already assigned to this pelayanan in this church")
	}

	// Auto-create user if person has email and user doesn't exist
	if person.Email != "" {
		_, err := s.userService.CreateUserFromPerson(ctx, person)
		if err != nil {
			// Log error but don't fail the assignment
			// User creation is optional, pelayanan assignment is the main operation
			fmt.Printf("Warning: Failed to create user for person %s: %v\n", person.ID, err)
		}
	}

	// Create the assignment
	assignment := &entity.PersonPelayananGereja{
		ID:          uuid.New(),
		PersonID:    request.PersonID,
		PelayananID: request.PelayananID,
		ChurchID:    request.ChurchID,
	}

	if err := s.pelayananRepo.CreatePelayananAssignment(ctx, assignment); err != nil {
		return fmt.Errorf("failed to assign pelayanan: %w", err)
	}

	// Update user activation status (activate if user exists)
	if err := s.userService.UpdateUserActivationStatus(ctx, request.PersonID); err != nil {
		// Log error but don't fail the assignment
		fmt.Printf("Warning: Failed to update user activation status for person %s: %v\n", person.ID, err)
	}

	return nil
}

func (s *pelayananService) UnassignPelayanan(ctx context.Context, assignmentID uuid.UUID) error {
	// Validate that assignment exists
	assignment, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("assignment with ID %s not found", assignmentID)
		}
		return fmt.Errorf("failed to validate assignment: %w", err)
	}

	// Store person ID before deleting assignment
	personID := assignment.PersonID

	if err := s.pelayananRepo.DeletePelayananAssignment(ctx, assignmentID); err != nil {
		return fmt.Errorf("failed to unassign pelayanan: %w", err)
	}

	// Update user activation status (may deactivate if no more assignments)
	if err := s.userService.UpdateUserActivationStatus(ctx, personID); err != nil {
		// Log error but don't fail the unassignment
		fmt.Printf("Warning: Failed to update user activation status for person %s: %v\n", personID, err)
	}

	return nil
}

func (s *pelayananService) GetAssignmentByID(ctx context.Context, assignmentID uuid.UUID) (dto.PelayananAssignmentResponse, error) {
	assignment, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.PelayananAssignmentResponse{}, fmt.Errorf("assignment with ID %s not found", assignmentID)
		}
		return dto.PelayananAssignmentResponse{}, fmt.Errorf("failed to get assignment: %w", err)
	}

	// Check if user account exists and get its status
	hasUserAccount := false
	isUserActive := false
	if user, err := s.userService.GetUserByPersonID(ctx, assignment.PersonID); err == nil {
		hasUserAccount = true
		isUserActive = user.IsActive
	}

	return dto.PelayananAssignmentResponse{
		ID:             assignment.ID,
		PersonID:       assignment.PersonID,
		PersonName:     assignment.Person.Nama,
		PelayananID:    assignment.PelayananID,
		Pelayanan:      assignment.Pelayanan.Pelayanan,
		ChurchID:       assignment.ChurchID,
		ChurchName:     assignment.Church.Name,
		DepartmentID:   assignment.Pelayanan.Department.ID,
		DepartmentName: assignment.Pelayanan.Department.Name,
		PelayananIsPic: assignment.Pelayanan.IsPic,
		HasUserAccount: hasUserAccount,
		IsUserActive:   isUserActive,
		CreatedAt:      assignment.CreatedAt,
		UpdatedAt:      assignment.UpdatedAt,
	}, nil
}
