package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type DepartmentService interface {
	Create(req *dto.DepartmentRequest) (*dto.DepartmentResponse, error)
	GetAll() ([]dto.DepartmentResponse, error)
	GetByID(id uuid.UUID) (*dto.DepartmentResponse, error)
	GetByChurchID(churchID uuid.UUID) ([]dto.DepartmentResponse, error)
	Update(id uuid.UUID, req *dto.DepartmentRequest) (*dto.DepartmentResponse, error)
	Delete(id uuid.UUID) error
}

type departmentService struct {
	departmentRepository repository.DepartmentRepository
}

func NewDepartmentService(departmentRepository repository.DepartmentRepository) DepartmentService {
	return &departmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *departmentService) Create(req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department := &entity.Department{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.departmentRepository.Create(department); err != nil {
		return nil, err
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

func (s *departmentService) GetByChurchID(churchID uuid.UUID) ([]dto.DepartmentResponse, error) {
	departments, err := s.departmentRepository.GetByChurchID(churchID)
	if err != nil {
		return nil, err
	}

	var responses []dto.DepartmentResponse
	for _, department := range departments {
		responses = append(responses, *s.toResponse(&department))
	}

	return responses, nil
}

func (s *departmentService) Update(id uuid.UUID, req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	department.Name = req.Name
	department.Description = req.Description

	if err := s.departmentRepository.Update(department); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *departmentService) Delete(id uuid.UUID) error {
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
