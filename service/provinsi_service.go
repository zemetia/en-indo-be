package service

import (
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type ProvinsiService interface {
	GetAll() ([]dto.ProvinsiResponse, error)
	GetByID(id uint) (*dto.ProvinsiResponse, error)
}

type provinsiService struct {
	provinsiRepository repository.ProvinsiRepository
}

func NewProvinsiService(provinsiRepository repository.ProvinsiRepository) ProvinsiService {
	return &provinsiService{
		provinsiRepository: provinsiRepository,
	}
}

func (s *provinsiService) GetAll() ([]dto.ProvinsiResponse, error) {
	provinsi, err := s.provinsiRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ProvinsiResponse
	for _, p := range provinsi {
		responses = append(responses, *s.toResponse(&p))
	}

	return responses, nil
}

func (s *provinsiService) GetByID(id uint) (*dto.ProvinsiResponse, error) {
	provinsi, err := s.provinsiRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(provinsi), nil
}

func (s *provinsiService) toResponse(provinsi *entity.Provinsi) *dto.ProvinsiResponse {
	return &dto.ProvinsiResponse{
		ID:   provinsi.ID,
		Name: provinsi.Name,
	}
}
