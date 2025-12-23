package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type LaguService interface {
	Create(ctx context.Context, lagu *entity.Lagu) error
	FindAll(ctx context.Context) ([]entity.Lagu, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lagu, error)
	Update(ctx context.Context, lagu *entity.Lagu) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type laguService struct {
	laguRepository repository.LaguRepository
}

func NewLaguService(laguRepository repository.LaguRepository) LaguService {
	return &laguService{
		laguRepository: laguRepository,
	}
}

func (s *laguService) Create(ctx context.Context, lagu *entity.Lagu) error {
	lagu.ID = uuid.New()
	return s.laguRepository.Create(ctx, lagu)
}

func (s *laguService) FindAll(ctx context.Context) ([]entity.Lagu, error) {
	return s.laguRepository.FindAll(ctx)
}

func (s *laguService) FindByID(ctx context.Context, id uuid.UUID) (*entity.Lagu, error) {
	return s.laguRepository.FindByID(ctx, id)
}

func (s *laguService) Update(ctx context.Context, lagu *entity.Lagu) error {
	return s.laguRepository.Update(ctx, lagu)
}

func (s *laguService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.laguRepository.Delete(ctx, id)
}
