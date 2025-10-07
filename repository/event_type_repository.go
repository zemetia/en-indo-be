package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type EventTypeRepository interface {
	Create(eventType *entity.EventType) error
	GetAll() ([]entity.EventType, error)
	GetActive() ([]entity.EventType, error)
	GetByID(id uuid.UUID) (*entity.EventType, error)
	GetByName(name string) (*entity.EventType, error)
	Update(eventType *entity.EventType) error
	Delete(id uuid.UUID) error
}

type eventTypeRepository struct {
	db *gorm.DB
}

func NewEventTypeRepository(db *gorm.DB) EventTypeRepository {
	return &eventTypeRepository{
		db: db,
	}
}

func (r *eventTypeRepository) Create(eventType *entity.EventType) error {
	return r.db.Create(eventType).Error
}

func (r *eventTypeRepository) GetAll() ([]entity.EventType, error) {
	var eventTypes []entity.EventType
	err := r.db.Order("sort_order ASC, display_name ASC").Find(&eventTypes).Error
	return eventTypes, err
}

func (r *eventTypeRepository) GetActive() ([]entity.EventType, error) {
	var eventTypes []entity.EventType
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, display_name ASC").Find(&eventTypes).Error
	return eventTypes, err
}

func (r *eventTypeRepository) GetByID(id uuid.UUID) (*entity.EventType, error) {
	var eventType entity.EventType
	err := r.db.First(&eventType, "id = ?", id).Error
	return &eventType, err
}

func (r *eventTypeRepository) GetByName(name string) (*entity.EventType, error) {
	var eventType entity.EventType
	err := r.db.First(&eventType, "name = ?", name).Error
	return &eventType, err
}

func (r *eventTypeRepository) Update(eventType *entity.EventType) error {
	return r.db.Save(eventType).Error
}

func (r *eventTypeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.EventType{}, "id = ?", id).Error
}
