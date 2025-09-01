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
	GetMyPelayanan(ctx context.Context, personID uuid.UUID) ([]dto.PersonHasPelayananResponse, error)
	GetAllAssignments(ctx context.Context, req dto.PaginationRequest) (dto.PelayananAssignmentPaginationResponse, error)
	GetAllPelayanan(ctx context.Context) ([]dto.PelayananResponse, error)
	AssignPelayanan(ctx context.Context, request dto.AssignPelayananRequest) error
	UnassignPelayanan(ctx context.Context, assignmentID uuid.UUID) error
	UpdatePelayananAssignment(ctx context.Context, assignmentID uuid.UUID, request dto.UpdatePelayananAssignmentRequest) error
	GetAssignmentByID(ctx context.Context, assignmentID uuid.UUID) (dto.PelayananAssignmentResponse, error)
}

type pelayananService struct {
	pelayananRepo repository.PelayananRepository
	personRepo    repository.PersonRepository
	churchRepo    repository.ChurchRepository
}

func NewPelayananService(
	pelayananRepo repository.PelayananRepository,
	personRepo repository.PersonRepository,
	churchRepo repository.ChurchRepository,
) PelayananService {
	return &pelayananService{
		pelayananRepo: pelayananRepo,
		personRepo:    personRepo,
		churchRepo:    churchRepo,
	}
}

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
			IsPic:       assignment.IsPic,
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
		responses = append(responses, dto.PelayananAssignmentResponse{
			ID:          assignment.ID,
			PersonID:    assignment.PersonID,
			PersonName:  assignment.Person.Nama,
			PelayananID: assignment.PelayananID,
			Pelayanan:   assignment.Pelayanan.Pelayanan,
			ChurchID:    assignment.ChurchID,
			ChurchName:  assignment.Church.Name,
			IsPic:       assignment.IsPic,
			CreatedAt:   assignment.CreatedAt,
			UpdatedAt:   assignment.UpdatedAt,
		})
	}

	return dto.PelayananAssignmentPaginationResponse{
		Data:               responses,
		PaginationResponse: *pagination,
	}, nil
}

func (s *pelayananService) GetAllPelayanan(ctx context.Context) ([]dto.PelayananResponse, error) {
	pelayanan, err := s.pelayananRepo.GetAllPelayanan(ctx)
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
	_, err := s.personRepo.GetByID(ctx, request.PersonID)
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

	// Create the assignment
	assignment := &entity.PersonPelayananGereja{
		ID:          uuid.New(),
		PersonID:    request.PersonID,
		PelayananID: request.PelayananID,
		ChurchID:    request.ChurchID,
		IsPic:       request.IsPic,
	}

	if err := s.pelayananRepo.CreatePelayananAssignment(ctx, assignment); err != nil {
		return fmt.Errorf("failed to assign pelayanan: %w", err)
	}

	return nil
}

func (s *pelayananService) UnassignPelayanan(ctx context.Context, assignmentID uuid.UUID) error {
	// Validate that assignment exists
	_, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("assignment with ID %s not found", assignmentID)
		}
		return fmt.Errorf("failed to validate assignment: %w", err)
	}

	if err := s.pelayananRepo.DeletePelayananAssignment(ctx, assignmentID); err != nil {
		return fmt.Errorf("failed to unassign pelayanan: %w", err)
	}

	return nil
}

func (s *pelayananService) UpdatePelayananAssignment(ctx context.Context, assignmentID uuid.UUID, request dto.UpdatePelayananAssignmentRequest) error {
	// Get existing assignment
	assignment, err := s.pelayananRepo.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("assignment with ID %s not found", assignmentID)
		}
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Update the assignment
	assignment.IsPic = request.IsPic

	if err := s.pelayananRepo.UpdatePelayananAssignment(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
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

	return dto.PelayananAssignmentResponse{
		ID:          assignment.ID,
		PersonID:    assignment.PersonID,
		PersonName:  assignment.Person.Nama,
		PelayananID: assignment.PelayananID,
		Pelayanan:   assignment.Pelayanan.Pelayanan,
		ChurchID:    assignment.ChurchID,
		ChurchName:  assignment.Church.Name,
		IsPic:       assignment.IsPic,
		CreatedAt:   assignment.CreatedAt,
		UpdatedAt:   assignment.UpdatedAt,
	}, nil
}