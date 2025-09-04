package service

import (
	"fmt"
	"log"

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
	log.Printf("[INFO] Church service: Creating church with data: %+v", req)
	
	// Check if church code already exists (unless it's empty)
	if req.ChurchCode != "" {
		// Check for existing church with same code
		existingChurches, err := s.churchRepository.GetAll()
		if err != nil {
			log.Printf("[ERROR] Church service: Failed to check existing churches: %v", err)
			return nil, err
		}
		
		for _, existing := range existingChurches {
			if existing.ChurchCode == req.ChurchCode {
				log.Printf("[ERROR] Church service: Church code '%s' already exists", req.ChurchCode)
				return nil, fmt.Errorf("church code '%s' already exists", req.ChurchCode)
			}
		}
	}
	
	church := &entity.Church{
		ID:          uuid.New(),
		Name:        req.Name,
		Address:     req.Address,
		ChurchCode:  req.ChurchCode,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		KabupatenID: req.KabupatenID,
	}

	log.Printf("[INFO] Church service: About to create church entity with ID: %s", church.ID)
	if err := s.churchRepository.Create(church); err != nil {
		log.Printf("[ERROR] Church service: Database error creating church: %v", err)
		log.Printf("[ERROR] Church service: Error type: %T", err)
		return nil, err
	}

	log.Printf("[INFO] Church service: Successfully created church with ID: %s", church.ID)
	log.Printf("[INFO] Church service: Fetching created church details")
	
	result, err := s.GetByID(church.ID)
	if err != nil {
		log.Printf("[ERROR] Church service: Failed to fetch created church: %v", err)
		return nil, err
	}
	
	return result, nil
}

func (s *ChurchService) GetAll() ([]dto.ChurchResponse, error) {
	churches, err := s.churchRepository.GetAll()
	if err != nil {
		log.Printf("[ERROR] Church service: Failed to get all churches: %v", err)
		return nil, err
	}

	log.Printf("[INFO] Church service: Retrieved %d churches from database", len(churches))
	
	var responses []dto.ChurchResponse
	for _, church := range churches {
		if response := s.toResponse(&church); response != nil {
			responses = append(responses, *response)
		}
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
	church.ChurchCode = req.ChurchCode
	church.Phone = req.Phone
	church.Email = req.Email
	church.Website = req.Website
	church.Latitude = req.Latitude
	church.Longitude = req.Longitude
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
		log.Printf("[WARN] Church service: toResponse received nil church")
		return nil
	}

	var kabupatenName string
	var provinsiName string
	var provinsiID uint
	
	// First try to get data from preloaded relationships
	if church.Kabupaten.ID > 0 {
		kabupatenName = church.Kabupaten.Name
		if church.Kabupaten.Provinsi.ID > 0 {
			provinsiID = church.Kabupaten.Provinsi.ID
			provinsiName = church.Kabupaten.Provinsi.Name
		} else {
			log.Printf("[WARN] Church service: Provinsi not preloaded for church %s (kabupaten: %s)", church.Name, kabupatenName)
		}
	} else if church.KabupatenID > 0 {
		// Fallback: fetch kabupaten data if not preloaded
		log.Printf("[WARN] Church service: Kabupaten not preloaded for church %s, fetching separately", church.Name)
		kabupaten, err := s.kabupatenRepository.GetByID(church.KabupatenID)
		if err != nil {
			log.Printf("[ERROR] Church service: Failed to get kabupaten ID %d for church %s: %v", church.KabupatenID, church.Name, err)
		} else if kabupaten != nil {
			kabupatenName = kabupaten.Name
			
			// Get provinsi name
			if kabupaten.ProvinsiID > 0 {
				provinsiID = kabupaten.ProvinsiID
				provinsi, err := s.provinsiRepository.GetByID(provinsiID)
				if err != nil {
					log.Printf("[ERROR] Church service: Failed to get provinsi ID %d for kabupaten %s: %v", provinsiID, kabupatenName, err)
				} else if provinsi != nil {
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
		ChurchCode:  church.ChurchCode,
		Phone:       church.Phone,
		Email:       church.Email,
		Website:     church.Website,
		Latitude:    church.Latitude,
		Longitude:   church.Longitude,
		KabupatenID: church.KabupatenID,
		Kabupaten:   kabupatenName,
		ProvinsiID:  provinsiID,
		Provinsi:    provinsiName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
