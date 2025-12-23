package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type KetersediaanService interface {
	CreateKetersediaan(req *dto.CreateKetersediaanRequest) (*dto.KetersediaanResponse, error)
	GetKetersediaan(id uuid.UUID) (*dto.KetersediaanResponse, error)
	UpdateKetersediaan(id uuid.UUID, req *dto.UpdateKetersediaanRequest) (*dto.KetersediaanResponse, error)
	DeleteKetersediaan(id uuid.UUID) error
	GetPersonAvailability(req *dto.GetKetersediaanRequest) (*dto.KetersediaanListResponse, error)
	GetEventAvailability(eventID uuid.UUID, occurrenceDate *time.Time) ([]dto.KetersediaanResponse, error)
	BulkCreateKetersediaan(req *dto.BulkCreateKetersediaanRequest) ([]dto.KetersediaanResponse, error)
	GetEventAvailabilitySummary(req *dto.GetEventAvailabilitySummaryRequest) (*dto.EventAvailabilitySummaryResponse, error)
	UpsertKetersediaan(req *dto.CreateKetersediaanRequest) (*dto.KetersediaanResponse, error)
	CleanupOldKetersediaan() error
}

type ketersediaanService struct {
	ketersediaanRepo repository.KetersediaanRepository
	personRepo       repository.PersonRepository
	eventRepo        repository.EventRepository
}

func NewKetersediaanService(
	ketersediaanRepo repository.KetersediaanRepository,
	personRepo repository.PersonRepository,
	eventRepo repository.EventRepository,
) KetersediaanService {
	return &ketersediaanService{
		ketersediaanRepo: ketersediaanRepo,
		personRepo:       personRepo,
		eventRepo:        eventRepo,
	}
}

func (s *ketersediaanService) CreateKetersediaan(req *dto.CreateKetersediaanRequest) (*dto.KetersediaanResponse, error) {
	// Validate person exists
	_, err := s.personRepo.GetByID(context.Background(), req.PersonID)
	if err != nil {
		return nil, fmt.Errorf("person not found: %w", err)
	}

	// Validate event exists
	_, err = s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// Check if availability already exists
	existing, err := s.ketersediaanRepo.FindByPersonAndEvent(req.PersonID, req.EventID, occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing availability: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("availability already exists for this person and event occurrence")
	}

	// Create entity
	ketersediaan := &entity.Ketersediaan{
		ID:             uuid.New(),
		PersonID:       req.PersonID,
		EventID:        req.EventID,
		OccurrenceDate: occurrenceDate,
		Status:         entity.KetersediaanStatus(req.Status),
		Notes:          req.Notes,
	}

	// Save to database
	if err := s.ketersediaanRepo.Create(ketersediaan); err != nil {
		return nil, fmt.Errorf("failed to create availability: %w", err)
	}

	// Fetch complete data with relations
	created, err := s.ketersediaanRepo.GetByID(ketersediaan.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created availability: %w", err)
	}

	return s.toResponse(created), nil
}

func (s *ketersediaanService) GetKetersediaan(id uuid.UUID) (*dto.KetersediaanResponse, error) {
	ketersediaan, err := s.ketersediaanRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("availability not found: %w", err)
	}

	return s.toResponse(ketersediaan), nil
}

func (s *ketersediaanService) UpdateKetersediaan(id uuid.UUID, req *dto.UpdateKetersediaanRequest) (*dto.KetersediaanResponse, error) {
	// Get existing
	ketersediaan, err := s.ketersediaanRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("availability not found: %w", err)
	}

	// Update fields
	if req.Status != nil {
		ketersediaan.Status = entity.KetersediaanStatus(*req.Status)
	}
	if req.Notes != nil {
		ketersediaan.Notes = *req.Notes
	}

	// Save
	if err := s.ketersediaanRepo.Update(ketersediaan); err != nil {
		return nil, fmt.Errorf("failed to update availability: %w", err)
	}

	// Fetch updated data
	updated, err := s.ketersediaanRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated availability: %w", err)
	}

	return s.toResponse(updated), nil
}

func (s *ketersediaanService) DeleteKetersediaan(id uuid.UUID) error {
	// Check exists
	_, err := s.ketersediaanRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("availability not found: %w", err)
	}

	if err := s.ketersediaanRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete availability: %w", err)
	}

	return nil
}

func (s *ketersediaanService) GetPersonAvailability(req *dto.GetKetersediaanRequest) (*dto.KetersediaanListResponse, error) {
	if req.PersonID == nil {
		return nil, fmt.Errorf("personId is required")
	}

	// Set defaults
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid startDate format: %w", err)
		}
		startDate = &parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format: %w", err)
		}
		endDate = &parsed
	}

	var statusPtr *string
	if req.Status != "" {
		statusPtr = &req.Status
	}

	ketersediaan, total, err := s.ketersediaanRepo.FindByPerson(*req.PersonID, startDate, endDate, statusPtr, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability: %w", err)
	}

	responses := make([]dto.KetersediaanResponse, len(ketersediaan))
	for i, k := range ketersediaan {
		responses[i] = *s.toResponse(&k)
	}

	return &dto.KetersediaanListResponse{
		Availabilities: responses,
		TotalCount:     int(total),
		Page:           page,
		Limit:          limit,
	}, nil
}

func (s *ketersediaanService) GetEventAvailability(eventID uuid.UUID, occurrenceDate *time.Time) ([]dto.KetersediaanResponse, error) {
	ketersediaan, err := s.ketersediaanRepo.FindByEvent(eventID, occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get event availability: %w", err)
	}

	responses := make([]dto.KetersediaanResponse, len(ketersediaan))
	for i, k := range ketersediaan {
		responses[i] = *s.toResponse(&k)
	}

	return responses, nil
}

func (s *ketersediaanService) BulkCreateKetersediaan(req *dto.BulkCreateKetersediaanRequest) ([]dto.KetersediaanResponse, error) {
	entities := make([]entity.Ketersediaan, 0, len(req.Availabilities))

	for _, avail := range req.Availabilities {
		// Validate person exists
		_, err := s.personRepo.GetByID(context.Background(), avail.PersonID)
		if err != nil {
			return nil, fmt.Errorf("person %s not found: %w", avail.PersonID, err)
		}

		// Validate event exists
		_, err = s.eventRepo.GetByID(avail.EventID)
		if err != nil {
			return nil, fmt.Errorf("event %s not found: %w", avail.EventID, err)
		}

		// Parse occurrence date
		occurrenceDate, err := time.Parse("2006-01-02", avail.OccurrenceDate)
		if err != nil {
			return nil, fmt.Errorf("invalid occurrence date format: %w", err)
		}

		// Check if already exists
		existing, err := s.ketersediaanRepo.FindByPersonAndEvent(avail.PersonID, avail.EventID, occurrenceDate)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing availability: %w", err)
		}
		if existing != nil {
			continue // Skip if already exists
		}

		entities = append(entities, entity.Ketersediaan{
			ID:             uuid.New(),
			PersonID:       avail.PersonID,
			EventID:        avail.EventID,
			OccurrenceDate: occurrenceDate,
			Status:         entity.KetersediaanStatus(avail.Status),
			Notes:          avail.Notes,
		})
	}

	if len(entities) == 0 {
		return []dto.KetersediaanResponse{}, nil
	}

	if err := s.ketersediaanRepo.BulkCreate(entities); err != nil {
		return nil, fmt.Errorf("failed to bulk create availability: %w", err)
	}

	// Fetch created entities with relations
	responses := make([]dto.KetersediaanResponse, len(entities))
	for i, entity := range entities {
		fetched, err := s.ketersediaanRepo.GetByID(entity.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch created availability: %w", err)
		}
		responses[i] = *s.toResponse(fetched)
	}

	return responses, nil
}

func (s *ketersediaanService) GetEventAvailabilitySummary(req *dto.GetEventAvailabilitySummaryRequest) (*dto.EventAvailabilitySummaryResponse, error) {
	// Validate event exists
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// Get summary
	summary, err := s.ketersediaanRepo.GetAvailabilitySummary(req.EventID, occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability summary: %w", err)
	}

	// Get all responses
	ketersediaan, err := s.ketersediaanRepo.FindByEvent(req.EventID, &occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get availability responses: %w", err)
	}

	responses := make([]dto.KetersediaanResponse, len(ketersediaan))
	for i, k := range ketersediaan {
		responses[i] = *s.toResponse(&k)
	}

	total := summary["available"] + summary["unavailable"] + summary["tentative"]

	return &dto.EventAvailabilitySummaryResponse{
		EventID:        req.EventID,
		EventTitle:     event.Title,
		OccurrenceDate: occurrenceDate,
		TotalResponses: total,
		Available:      summary["available"],
		Unavailable:    summary["unavailable"],
		Tentative:      summary["tentative"],
		Responses:      responses,
	}, nil
}

func (s *ketersediaanService) UpsertKetersediaan(req *dto.CreateKetersediaanRequest) (*dto.KetersediaanResponse, error) {
	// Parse occurrence date
	occurrenceDate, err := time.Parse("2006-01-02", req.OccurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid occurrence date format: %w", err)
	}

	// Check if availability already exists
	existing, err := s.ketersediaanRepo.FindByPersonAndEvent(req.PersonID, req.EventID, occurrenceDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing availability: %w", err)
	}

	if existing != nil {
		// Update existing
		existing.Status = entity.KetersediaanStatus(req.Status)
		existing.Notes = req.Notes

		if err := s.ketersediaanRepo.Update(existing); err != nil {
			return nil, fmt.Errorf("failed to update availability: %w", err)
		}

		updated, err := s.ketersediaanRepo.GetByID(existing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch updated availability: %w", err)
		}

		return s.toResponse(updated), nil
	}

	// Create new
	return s.CreateKetersediaan(req)
}

func (s *ketersediaanService) toResponse(k *entity.Ketersediaan) *dto.KetersediaanResponse {
	response := &dto.KetersediaanResponse{
		ID:             k.ID,
		PersonID:       k.PersonID,
		EventID:        k.EventID,
		OccurrenceDate: k.OccurrenceDate,
		Status:         string(k.Status),
		Notes:          k.Notes,
		CreatedAt:      k.CreatedAt,
		UpdatedAt:      k.UpdatedAt,
	}

	if k.Person.ID != uuid.Nil {
		response.PersonName = k.Person.Nama
		response.Person = &dto.PersonBasicResponse{
			ID:    k.Person.ID,
			Nama:  k.Person.Nama,
			Email: k.Person.Email,
		}
	}

	if k.Event.ID != uuid.Nil {
		response.EventTitle = k.Event.Title
	}

	return response
}

func (s *ketersediaanService) CleanupOldKetersediaan() error {
	// Calculate date 2 months ago
	retentionDate := time.Now().AddDate(0, -2, 0)

	if err := s.ketersediaanRepo.CleanupOldKetersediaan(retentionDate); err != nil {
		return fmt.Errorf("failed to cleanup old availability: %w", err)
	}

	return nil
}
