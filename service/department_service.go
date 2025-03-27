package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type DepartmentService struct {
	departmentRepository *repository.DepartmentRepository
}

func NewDepartmentService(departmentRepository *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *DepartmentService) Create(req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department := &entity.Department{
		Name:        req.Name,
		Description: req.Description,
		ChurchID:    req.ChurchID,
	}

	if err := s.departmentRepository.Create(department); err != nil {
		return nil, err
	}

	return s.GetByID(department.ID)
}

func (s *DepartmentService) GetAll() ([]dto.DepartmentResponse, error) {
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

func (s *DepartmentService) GetByID(id uuid.UUID) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(department), nil
}

func (s *DepartmentService) GetByChurchID(churchID uuid.UUID) ([]dto.DepartmentResponse, error) {
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

func (s *DepartmentService) Update(id uuid.UUID, req *dto.DepartmentRequest) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	department.Name = req.Name
	department.Description = req.Description
	department.ChurchID = req.ChurchID

	if err := s.departmentRepository.Update(department); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *DepartmentService) Delete(id uuid.UUID) error {
	return s.departmentRepository.Delete(id)
}

func (s *DepartmentService) toResponse(department *entity.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		Description: department.Description,
		ChurchID:    department.ChurchID,
		CreatedAt:   department.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   department.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
