package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type LifeGroupEventService interface {
	Create(event *entity.LifeGroupEvent) error
	Update(event *entity.LifeGroupEvent) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.LifeGroupEvent, error)
	FindByLifeGroupID(lifeGroupID uuid.UUID) ([]entity.LifeGroupEvent, error)
}

type lifeGroupEventService struct {
	repo repository.LifeGroupEventRepository
}

func NewLifeGroupEventService(repo repository.LifeGroupEventRepository) LifeGroupEventService {
	return &lifeGroupEventService{repo: repo}
}

func (s *lifeGroupEventService) Create(event *entity.LifeGroupEvent) error {
	event.ID = uuid.New()
	return s.repo.Create(event)
}

func (s *lifeGroupEventService) Update(event *entity.LifeGroupEvent) error {
	return s.repo.Update(event)
}

func (s *lifeGroupEventService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *lifeGroupEventService) FindByID(id uuid.UUID) (*entity.LifeGroupEvent, error) {
	return s.repo.FindByID(id)
}

func (s *lifeGroupEventService) FindByLifeGroupID(lifeGroupID uuid.UUID) ([]entity.LifeGroupEvent, error) {
	return s.repo.FindByLifeGroupID(lifeGroupID)
}
