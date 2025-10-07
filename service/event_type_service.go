package service

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type EventTypeService struct {
	eventTypeRepository repository.EventTypeRepository
}

func NewEventTypeService(eventTypeRepository repository.EventTypeRepository) *EventTypeService {
	return &EventTypeService{
		eventTypeRepository: eventTypeRepository,
	}
}

// toSnakeCase converts a string to snake_case format
// Example: "Spiritual Journey" -> "spiritual_journey"
func toSnakeCase(str string) string {
	// Convert to lowercase
	str = strings.ToLower(str)

	// Replace spaces and special characters with underscore
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	str = reg.ReplaceAllString(str, "_")

	// Remove leading/trailing underscores
	str = strings.Trim(str, "_")

	// Replace multiple consecutive underscores with single underscore
	reg = regexp.MustCompile(`_+`)
	str = reg.ReplaceAllString(str, "_")

	return str
}

func (s *EventTypeService) Create(req *dto.EventTypeRequest) (*dto.EventTypeResponse, error) {
	log.Printf("[INFO] EventType service: Creating event type with data: %+v", req)

	// Auto-generate Name from DisplayName if not provided
	name := req.Name
	if name == "" {
		name = toSnakeCase(req.DisplayName)
		log.Printf("[INFO] EventType service: Auto-generated name '%s' from display name '%s'", name, req.DisplayName)
	}

	// Check if event type name already exists
	existing, err := s.eventTypeRepository.GetByName(name)
	if err == nil && existing != nil {
		log.Printf("[ERROR] EventType service: Event type name '%s' already exists", name)
		return nil, fmt.Errorf("event type name '%s' already exists", name)
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#2980b9"
	}

	eventType := &entity.EventType{
		ID:          uuid.New(),
		Name:        name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Color:       color,
		IsActive:    true,
		SortOrder:   req.SortOrder,
	}

	// Override IsActive if explicitly provided
	if req.IsActive != nil {
		eventType.IsActive = *req.IsActive
	}

	log.Printf("[INFO] EventType service: About to create event type entity with ID: %s", eventType.ID)
	if err := s.eventTypeRepository.Create(eventType); err != nil {
		log.Printf("[ERROR] EventType service: Database error creating event type: %v", err)
		return nil, err
	}

	log.Printf("[INFO] EventType service: Successfully created event type with ID: %s", eventType.ID)

	result, err := s.GetByID(eventType.ID)
	if err != nil {
		log.Printf("[ERROR] EventType service: Failed to fetch created event type: %v", err)
		return nil, err
	}

	return result, nil
}

func (s *EventTypeService) GetAll() ([]dto.EventTypeResponse, error) {
	eventTypes, err := s.eventTypeRepository.GetAll()
	if err != nil {
		log.Printf("[ERROR] EventType service: Failed to get all event types: %v", err)
		return nil, err
	}

	log.Printf("[INFO] EventType service: Retrieved %d event types from database", len(eventTypes))

	var responses []dto.EventTypeResponse
	for _, eventType := range eventTypes {
		if response := s.toResponse(&eventType); response != nil {
			responses = append(responses, *response)
		}
	}

	return responses, nil
}

func (s *EventTypeService) GetActive() ([]dto.EventTypeResponse, error) {
	eventTypes, err := s.eventTypeRepository.GetActive()
	if err != nil {
		log.Printf("[ERROR] EventType service: Failed to get active event types: %v", err)
		return nil, err
	}

	log.Printf("[INFO] EventType service: Retrieved %d active event types from database", len(eventTypes))

	var responses []dto.EventTypeResponse
	for _, eventType := range eventTypes {
		if response := s.toResponse(&eventType); response != nil {
			responses = append(responses, *response)
		}
	}

	return responses, nil
}

func (s *EventTypeService) GetByID(id uuid.UUID) (*dto.EventTypeResponse, error) {
	eventType, err := s.eventTypeRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(eventType), nil
}

func (s *EventTypeService) Update(id uuid.UUID, req *dto.EventTypeRequest) (*dto.EventTypeResponse, error) {
	eventType, err := s.eventTypeRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Auto-generate Name from DisplayName if not provided
	name := req.Name
	if name == "" {
		name = toSnakeCase(req.DisplayName)
		log.Printf("[INFO] EventType service: Auto-generated name '%s' from display name '%s'", name, req.DisplayName)
	}

	// Check if new name conflicts with another event type
	if name != eventType.Name {
		existing, err := s.eventTypeRepository.GetByName(name)
		if err == nil && existing != nil && existing.ID != id {
			log.Printf("[ERROR] EventType service: Event type name '%s' already exists", name)
			return nil, fmt.Errorf("event type name '%s' already exists", name)
		}
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#2980b9"
	}

	eventType.Name = name
	eventType.DisplayName = req.DisplayName
	eventType.Description = req.Description
	eventType.Color = color
	eventType.SortOrder = req.SortOrder

	// Update IsActive if explicitly provided
	if req.IsActive != nil {
		eventType.IsActive = *req.IsActive
	}

	if err := s.eventTypeRepository.Update(eventType); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *EventTypeService) Delete(id uuid.UUID) error {
	return s.eventTypeRepository.Delete(id)
}

func (s *EventTypeService) toResponse(eventType *entity.EventType) *dto.EventTypeResponse {
	if eventType == nil {
		log.Printf("[WARN] EventType service: toResponse received nil event type")
		return nil
	}

	// Format dates safely
	var createdAt, updatedAt string
	if !eventType.CreatedAt.IsZero() {
		createdAt = eventType.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !eventType.UpdatedAt.IsZero() {
		updatedAt = eventType.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &dto.EventTypeResponse{
		ID:          eventType.ID,
		Name:        eventType.Name,
		DisplayName: eventType.DisplayName,
		Description: eventType.Description,
		Color:       eventType.Color,
		IsActive:    eventType.IsActive,
		SortOrder:   eventType.SortOrder,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
