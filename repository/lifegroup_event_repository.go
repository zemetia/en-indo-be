package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupEventRepository interface {
	Create(event *entity.LifeGroupEvent) error
	Update(event *entity.LifeGroupEvent) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.LifeGroupEvent, error)
	FindByLifeGroupID(lifeGroupID uuid.UUID) ([]entity.LifeGroupEvent, error)
}

type lifeGroupEventRepository struct {
	db *gorm.DB
}

func NewLifeGroupEventRepository(db *gorm.DB) LifeGroupEventRepository {
	return &lifeGroupEventRepository{db: db}
}

func (r *lifeGroupEventRepository) Create(event *entity.LifeGroupEvent) error {
	return r.db.Create(event).Error
}

func (r *lifeGroupEventRepository) Update(event *entity.LifeGroupEvent) error {
	return r.db.Save(event).Error
}

func (r *lifeGroupEventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.LifeGroupEvent{}, id).Error
}

func (r *lifeGroupEventRepository) FindByID(id uuid.UUID) (*entity.LifeGroupEvent, error) {
	var event entity.LifeGroupEvent
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *lifeGroupEventRepository) FindByLifeGroupID(lifeGroupID uuid.UUID) ([]entity.LifeGroupEvent, error) {
	var events []entity.LifeGroupEvent
	err := r.db.Where("life_group_id = ?", lifeGroupID).Find(&events).Error
	return events, err
}
