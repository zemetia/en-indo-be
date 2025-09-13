package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type VisitorService interface {
	Create(ctx context.Context, req *dto.VisitorRequest) (*dto.VisitorResponse, error)
	GetAll(ctx context.Context) ([]dto.VisitorResponse, error)
	Search(ctx context.Context, search *dto.VisitorSearchDto) ([]dto.VisitorResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.VisitorResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.VisitorRequest) (*dto.VisitorResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitorService struct {
	visitorRepo repository.VisitorRepository
}

func NewVisitorService(visitorRepo repository.VisitorRepository) VisitorService {
	return &visitorService{
		visitorRepo: visitorRepo,
	}
}

func (s *visitorService) Create(ctx context.Context, req *dto.VisitorRequest) (*dto.VisitorResponse, error) {
	visitor := &entity.Visitor{
		Name:        req.Name,
		IGUsername:  req.IGUsername,
		PhoneNumber: req.PhoneNumber,
		KabupatenID: req.KabupatenID,
	}

	err := s.visitorRepo.Create(ctx, visitor)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitor), nil
}

func (s *visitorService) GetAll(ctx context.Context) ([]dto.VisitorResponse, error) {
	visitors, err := s.visitorRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.VisitorResponse
	for _, visitor := range visitors {
		responses = append(responses, *s.entityToResponse(&visitor))
	}

	return responses, nil
}

func (s *visitorService) Search(ctx context.Context, search *dto.VisitorSearchDto) ([]dto.VisitorResponse, error) {
	visitors, err := s.visitorRepo.Search(ctx, search)
	if err != nil {
		return nil, err
	}

	var responses []dto.VisitorResponse
	for _, visitor := range visitors {
		responses = append(responses, *s.entityToResponse(&visitor))
	}

	return responses, nil
}

func (s *visitorService) GetByID(ctx context.Context, id uuid.UUID) (*dto.VisitorResponse, error) {
	visitor, err := s.visitorRepo.GetWithInformation(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitor), nil
}

func (s *visitorService) Update(ctx context.Context, id uuid.UUID, req *dto.VisitorRequest) (*dto.VisitorResponse, error) {
	visitor, err := s.visitorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	visitor.Name = req.Name
	visitor.IGUsername = req.IGUsername
	visitor.PhoneNumber = req.PhoneNumber
	visitor.KabupatenID = req.KabupatenID

	err = s.visitorRepo.Update(ctx, visitor)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(visitor), nil
}

func (s *visitorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.visitorRepo.Delete(ctx, id)
}

func (s *visitorService) entityToResponse(visitor *entity.Visitor) *dto.VisitorResponse {
	response := &dto.VisitorResponse{
		ID:          visitor.ID,
		Name:        visitor.Name,
		IGUsername:  visitor.IGUsername,
		PhoneNumber: visitor.PhoneNumber,
		KabupatenID: visitor.KabupatenID,
		Kabupaten: func() string {
			if visitor.KabupatenID != nil {
				return visitor.Kabupaten.Name
			}
			return ""
		}(),
		ProvinsiID: func() *uint {
			if visitor.KabupatenID != nil {
				return &visitor.Kabupaten.ProvinsiID
			}
			return nil
		}(),
		Provinsi: func() string {
			if visitor.KabupatenID != nil {
				return visitor.Kabupaten.Provinsi.Name
			}
			return ""
		}(),
		CreatedAt: visitor.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: visitor.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	var information []dto.VisitorInformationResponse
	for _, info := range visitor.Information {
		information = append(information, dto.VisitorInformationResponse{
			ID:        info.ID,
			VisitorID: info.VisitorID,
			Label:     info.Label,
			Value:     info.Value,
			CreatedAt: info.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: info.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.Information = information

	return response
}
