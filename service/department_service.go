package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type DepartmentService interface {
	Create(req *dto.DepartmentRequest) (*dto.DepartmentResponse, error)
	GetAll() ([]dto.DepartmentResponse, error)
	GetByID(id uuid.UUID) (*dto.DepartmentResponse, error)
	Update(id uuid.UUID, req *dto.DepartmentRequest) (*dto.DepartmentResponse, error)
	Delete(id uuid.UUID) error
}

type departmentService struct {
	departmentRepository repository.DepartmentRepository
	pelayananRepository  repository.PelayananRepository
}

func NewDepartmentService(departmentRepository repository.DepartmentRepository, pelayananRepository repository.PelayananRepository) DepartmentService {
	return &departmentService{
		departmentRepository: departmentRepository,
		pelayananRepository:  pelayananRepository,
	}
}

func (s *departmentService) Create(req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department := &entity.Department{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.departmentRepository.Create(department); err != nil {
		return nil, err
	}

	// Create PIC pelayanan for this department
	picPelayanan := &entity.Pelayanan{
		ID:           uuid.New(),
		Pelayanan:    fmt.Sprintf("PIC %s", department.Name),
		Description:  fmt.Sprintf("Person in Charge untuk departemen %s", department.Name),
		DepartmentID: department.ID,
		IsPic:        true,
	}

	if err := s.pelayananRepository.CreatePelayanan(context.Background(), picPelayanan); err != nil {
		// If PIC creation fails, we should rollback the department creation
		s.departmentRepository.Delete(department.ID)
		return nil, fmt.Errorf("failed to create PIC pelayanan: %w", err)
	}

	return s.GetByID(department.ID)
}

func (s *departmentService) GetAll() ([]dto.DepartmentResponse, error) {
	departments, err := s.departmentRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.DepartmentResponse
	for _, department := range departments {
		responses = append(responses, *s.toResponse(&department))
	}

	return responses, nil
}

func (s *departmentService) GetByID(id uuid.UUID) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(department), nil
}

func (s *departmentService) Update(id uuid.UUID, req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldName := department.Name
	department.Name = req.Name
	department.Description = req.Description

	if err := s.departmentRepository.Update(department); err != nil {
		return nil, err
	}

	// Update PIC pelayanan name if department name changed
	if oldName != req.Name {
		picPelayanan, err := s.pelayananRepository.GetPelayananByDepartmentAndPic(context.Background(), id, true)
		if err == nil {
			picPelayanan.Pelayanan = fmt.Sprintf("PIC %s", req.Name)
			picPelayanan.Description = fmt.Sprintf("Person in Charge untuk departemen %s", req.Name)
			s.pelayananRepository.UpdatePelayanan(context.Background(), picPelayanan)
		}
	}

	return s.GetByID(id)
}

func (s *departmentService) Delete(id uuid.UUID) error {
	// Delete PIC pelayanan first
	picPelayanan, err := s.pelayananRepository.GetPelayananByDepartmentAndPic(context.Background(), id, true)
	if err == nil {
		s.pelayananRepository.DeletePelayanan(context.Background(), picPelayanan.ID)
	}

	return s.departmentRepository.Delete(id)
}

func (s *departmentService) toResponse(department *entity.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		Description: department.Description,
		CreatedAt:   department.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   department.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
