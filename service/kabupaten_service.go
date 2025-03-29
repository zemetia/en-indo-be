package service

import (
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type KabupatenService struct {
	kabupatenRepository repository.KabupatenRepository
}

func NewKabupatenService(kabupatenRepository repository.KabupatenRepository) *KabupatenService {
	return &KabupatenService{
		kabupatenRepository: kabupatenRepository,
	}
}

func (s *KabupatenService) GetAll() ([]dto.KabupatenResponse, error) {
	kabupaten, err := s.kabupatenRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.KabupatenResponse
	for _, k := range kabupaten {
		responses = append(responses, *s.toResponse(&k))
	}

	return responses, nil
}

func (s *KabupatenService) GetByID(id uint) (*dto.KabupatenResponse, error) {
	kabupaten, err := s.kabupatenRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(kabupaten), nil
}

func (s *KabupatenService) GetByProvinsiID(provinsiID uint) ([]dto.KabupatenResponse, error) {
	kabupaten, err := s.kabupatenRepository.GetByProvinsiID(provinsiID)
	if err != nil {
		return nil, err
	}

	var responses []dto.KabupatenResponse
	for _, k := range kabupaten {
		responses = append(responses, *s.toResponse(&k))
	}

	return responses, nil
}

func (s *KabupatenService) toResponse(kabupaten *entity.Kabupaten) *dto.KabupatenResponse {
	return &dto.KabupatenResponse{
		ID:         kabupaten.ID,
		Name:       kabupaten.Name,
		ProvinsiID: kabupaten.ProvinsiID,
		Provinsi:   kabupaten.Provinsi.Name,
	}
}
