package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/repository"
)

type EventDepartmentService interface {
	GetEventDepartments(ctx context.Context, eventID uuid.UUID) ([]dto.EventDepartmentResponse, error)
	UpdateEventDepartments(ctx context.Context, eventID uuid.UUID, req dto.UpdateEventDepartmentsRequest) error
}

type eventDepartmentService struct {
	eventDepartmentRepo repository.EventDepartmentRepository
}

func NewEventDepartmentService(eventDepartmentRepo repository.EventDepartmentRepository) EventDepartmentService {
	return &eventDepartmentService{
		eventDepartmentRepo: eventDepartmentRepo,
	}
}

func (s *eventDepartmentService) GetEventDepartments(ctx context.Context, eventID uuid.UUID) ([]dto.EventDepartmentResponse, error) {
	departments, err := s.eventDepartmentRepo.GetEventDepartments(ctx, eventID)
	if err != nil {
		return nil, err
	}

	var response []dto.EventDepartmentResponse
	for _, dept := range departments {
		response = append(response, dto.EventDepartmentResponse{
			ID:          dept.ID,
			Name:        dept.Name,
			Description: dept.Description,
		})
	}

	return response, nil
}

func (s *eventDepartmentService) UpdateEventDepartments(ctx context.Context, eventID uuid.UUID, req dto.UpdateEventDepartmentsRequest) error {
	return s.eventDepartmentRepo.UpdateEventDepartments(ctx, eventID, req.DepartmentIDs)
}
