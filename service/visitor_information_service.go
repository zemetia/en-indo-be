package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type VisitorInformationService interface {
	Create(ctx context.Context, req *dto.VisitorInformationRequest) (*dto.VisitorInformationResponse, error)
	GetAll(ctx context.Context) ([]dto.VisitorInformationResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.VisitorInformationResponse, error)
	GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]dto.VisitorInformationResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.VisitorInformationRequest) (*dto.VisitorInformationResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitorInformationService struct {
	visitorInfoRepo repository.VisitorInformationRepository
	visitorRepo     repository.VisitorRepository
}

func NewVisitorInformationService(
	visitorInfoRepo repository.VisitorInformationRepository,
	visitorRepo repository.VisitorRepository,
) VisitorInformationService {
	return &visitorInformationService{
		visitorInfoRepo: visitorInfoRepo,
		visitorRepo:     visitorRepo,
	}
}

func (s *visitorInformationService) Create(ctx context.Context, req *dto.VisitorInformationRequest) (*dto.VisitorInformationResponse, error) {
	// Validate visitor exists
	_, err := s.visitorRepo.GetByID(ctx, req.VisitorID)
	if err != nil {
		return nil, err
	}

	visitorInfo := &entity.VisitorInformation{
		VisitorID: req.VisitorID,
		Label:     req.Label,
		Value:     req.Value,
	}

	err = s.visitorInfoRepo.Create(ctx, visitorInfo)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitorInfo), nil
}

func (s *visitorInformationService) GetAll(ctx context.Context) ([]dto.VisitorInformationResponse, error) {
	visitorInfos, err := s.visitorInfoRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.VisitorInformationResponse
	for _, info := range visitorInfos {
		responses = append(responses, *s.entityToResponse(&info))
	}

	return responses, nil
}

func (s *visitorInformationService) GetByID(ctx context.Context, id uuid.UUID) (*dto.VisitorInformationResponse, error) {
	visitorInfo, err := s.visitorInfoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitorInfo), nil
}

func (s *visitorInformationService) GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]dto.VisitorInformationResponse, error) {
	visitorInfos, err := s.visitorInfoRepo.GetByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, err
	}

	var responses []dto.VisitorInformationResponse
	for _, info := range visitorInfos {
		responses = append(responses, *s.entityToResponse(&info))
	}

	return responses, nil
}

func (s *visitorInformationService) Update(ctx context.Context, id uuid.UUID, req *dto.VisitorInformationRequest) (*dto.VisitorInformationResponse, error) {
	visitorInfo, err := s.visitorInfoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	visitorInfo.VisitorID = req.VisitorID
	visitorInfo.Label = req.Label
	visitorInfo.Value = req.Value

	err = s.visitorInfoRepo.Update(ctx, visitorInfo)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitorInfo), nil
}

func (s *visitorInformationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.visitorInfoRepo.Delete(ctx, id)
}

func (s *visitorInformationService) entityToResponse(visitorInfo *entity.VisitorInformation) *dto.VisitorInformationResponse {
	response := &dto.VisitorInformationResponse{
		ID:        visitorInfo.ID,
		VisitorID: visitorInfo.VisitorID,
		Label:     visitorInfo.Label,
		Value:     visitorInfo.Value,
		CreatedAt: visitorInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: visitorInfo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if visitorInfo.Visitor.ID != uuid.Nil {
		response.Visitor = dto.VisitorSimpleResponse{
			ID:          visitorInfo.Visitor.ID,
			Name:        visitorInfo.Visitor.Name,
			IGUsername:  visitorInfo.Visitor.IGUsername,
			PhoneNumber: visitorInfo.Visitor.PhoneNumber,
			KabupatenID: visitorInfo.Visitor.KabupatenID,
			Kabupaten: func() string {
				if visitorInfo.Visitor.KabupatenID != nil {
					return visitorInfo.Visitor.Kabupaten.Name
				}
				return ""
			}(),
		}
	}

	return response
}
