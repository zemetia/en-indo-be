package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type ChurchService struct {
	churchRepository    repository.ChurchRepository
	kabupatenRepository repository.KabupatenRepository
	provinsiRepository  repository.ProvinsiRepository
}

func NewChurchService(churchRepository repository.ChurchRepository, kabupatenRepository repository.KabupatenRepository, provinsiRepository repository.ProvinsiRepository) *ChurchService {
	return &ChurchService{
		churchRepository:    churchRepository,
		kabupatenRepository: kabupatenRepository,
		provinsiRepository:  provinsiRepository,
	}
}

func (s *ChurchService) Create(req *dto.ChurchRequest) (*dto.ChurchResponse, error) {
	church := &entity.Church{
		Name:        req.Name,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		KabupatenID: req.KabupatenID,
	}

	if err := s.churchRepository.Create(church); err != nil {
		return nil, err
	}

	return s.GetByID(church.ID)
}

func (s *ChurchService) GetAll() ([]dto.ChurchResponse, error) {
	churches, err := s.churchRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ChurchResponse
	for _, church := range churches {
		responses = append(responses, *s.toResponse(&church))
	}

	return responses, nil
}

func (s *ChurchService) GetByID(id uuid.UUID) (*dto.ChurchResponse, error) {
	church, err := s.churchRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(church), nil
}

func (s *ChurchService) GetByKabupatenID(kabupatenID uint) ([]dto.ChurchResponse, error) {
	churches, err := s.churchRepository.GetByKabupatenID(kabupatenID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ChurchResponse
	for _, church := range churches {
		responses = append(responses, *s.toResponse(&church))
	}

	return responses, nil
}

func (s *ChurchService) GetByProvinsiID(provinsiID uint) ([]dto.ChurchResponse, error) {
	churches, err := s.churchRepository.GetByProvinsiID(provinsiID)
	if err != nil {
		return nil, err
	}

	var responses []dto.ChurchResponse
	for _, church := range churches {
		responses = append(responses, *s.toResponse(&church))
	}

	return responses, nil
}

func (s *ChurchService) Update(id uuid.UUID, req *dto.ChurchRequest) (*dto.ChurchResponse, error) {
	church, err := s.churchRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	church.Name = req.Name
	church.Address = req.Address
	church.Phone = req.Phone
	church.Email = req.Email
	church.KabupatenID = req.KabupatenID

	if err := s.churchRepository.Update(church); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *ChurchService) Delete(id uuid.UUID) error {
	return s.churchRepository.Delete(id)
}

func (s *ChurchService) toResponse(church *entity.Church) *dto.ChurchResponse {
	if church == nil {
		return nil
	}

	// Get kabupaten name with better error handling
	var kabupatenName string
	var provinsiName string
	var provinsiID uint
	
	if church.KabupatenID > 0 {
		kabupaten, err := s.kabupatenRepository.GetByID(church.KabupatenID)
		if err == nil && kabupaten != nil {
			kabupatenName = kabupaten.Name
			
			// Get provinsi name
			if kabupaten.ProvinsiID > 0 {
				provinsiID = kabupaten.ProvinsiID
				provinsi, err := s.provinsiRepository.GetByID(provinsiID)
				if err == nil && provinsi != nil {
					provinsiName = provinsi.Name
				}
			}
		}
	}

	// Format dates safely
	var createdAt, updatedAt string
	if !church.CreatedAt.IsZero() {
		createdAt = church.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !church.UpdatedAt.IsZero() {
		updatedAt = church.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &dto.ChurchResponse{
		ID:          church.ID,
		Name:        church.Name,
		Address:     church.Address,
		Phone:       church.Phone,
		Email:       church.Email,
		KabupatenID: church.KabupatenID,
		Kabupaten:   kabupatenName,
		ProvinsiID:  provinsiID,
		Provinsi:    provinsiName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
