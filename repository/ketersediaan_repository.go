package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type KetersediaanRepository interface {
	Create(ketersediaan *entity.Ketersediaan) error
	GetByID(id uuid.UUID) (*entity.Ketersediaan, error)
	Update(ketersediaan *entity.Ketersediaan) error
	Delete(id uuid.UUID) error
	FindByPersonAndEvent(personID, eventID uuid.UUID, occurrenceDate time.Time) (*entity.Ketersediaan, error)
	FindByPerson(personID uuid.UUID, startDate, endDate *time.Time, status *string, page, limit int) ([]entity.Ketersediaan, int64, error)
	FindByEvent(eventID uuid.UUID, occurrenceDate *time.Time) ([]entity.Ketersediaan, error)
	FindByDateRange(personID *uuid.UUID, startDate, endDate time.Time) ([]entity.Ketersediaan, error)
	BulkCreate(ketersediaan []entity.Ketersediaan) error
	GetAvailabilitySummary(eventID uuid.UUID, occurrenceDate time.Time) (map[string]int, error)
}

type ketersediaanRepository struct {
	db *gorm.DB
}

func NewKetersediaanRepository(db *gorm.DB) KetersediaanRepository {
	return &ketersediaanRepository{
		db: db,
	}
}

func (r *ketersediaanRepository) Create(ketersediaan *entity.Ketersediaan) error {
	return r.db.Create(ketersediaan).Error
}

func (r *ketersediaanRepository) GetByID(id uuid.UUID) (*entity.Ketersediaan, error) {
	var ketersediaan entity.Ketersediaan
	err := r.db.Preload("Person").Preload("Event").First(&ketersediaan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ketersediaan, nil
}

func (r *ketersediaanRepository) Update(ketersediaan *entity.Ketersediaan) error {
	return r.db.Save(ketersediaan).Error
}

func (r *ketersediaanRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Ketersediaan{}, "id = ?", id).Error
}

func (r *ketersediaanRepository) FindByPersonAndEvent(personID, eventID uuid.UUID, occurrenceDate time.Time) (*entity.Ketersediaan, error) {
	var ketersediaan entity.Ketersediaan
	err := r.db.Where("person_id = ? AND event_id = ? AND occurrence_date = ?", personID, eventID, occurrenceDate).
		Preload("Person").
		Preload("Event").
		First(&ketersediaan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ketersediaan, nil
}

func (r *ketersediaanRepository) FindByPerson(personID uuid.UUID, startDate, endDate *time.Time, status *string, page, limit int) ([]entity.Ketersediaan, int64, error) {
	var ketersediaan []entity.Ketersediaan
	var total int64

	query := r.db.Model(&entity.Ketersediaan{}).Where("person_id = ?", personID)

	if startDate != nil {
		query = query.Where("occurrence_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("occurrence_date <= ?", endDate)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Preload("Person").Preload("Event").Order("occurrence_date ASC").Find(&ketersediaan).Error
	if err != nil {
		return nil, 0, err
	}

	return ketersediaan, total, nil
}

func (r *ketersediaanRepository) FindByEvent(eventID uuid.UUID, occurrenceDate *time.Time) ([]entity.Ketersediaan, error) {
	var ketersediaan []entity.Ketersediaan
	query := r.db.Where("event_id = ?", eventID)

	if occurrenceDate != nil {
		query = query.Where("occurrence_date = ?", occurrenceDate)
	}

	err := query.Preload("Person").Preload("Event").Order("person_id ASC").Find(&ketersediaan).Error
	if err != nil {
		return nil, err
	}

	return ketersediaan, nil
}

func (r *ketersediaanRepository) FindByDateRange(personID *uuid.UUID, startDate, endDate time.Time) ([]entity.Ketersediaan, error) {
	var ketersediaan []entity.Ketersediaan
	query := r.db.Where("occurrence_date >= ? AND occurrence_date <= ?", startDate, endDate)

	if personID != nil {
		query = query.Where("person_id = ?", personID)
	}

	err := query.Preload("Person").Preload("Event").Order("occurrence_date ASC, event_id ASC").Find(&ketersediaan).Error
	if err != nil {
		return nil, err
	}

	return ketersediaan, nil
}

func (r *ketersediaanRepository) BulkCreate(ketersediaan []entity.Ketersediaan) error {
	return r.db.Create(&ketersediaan).Error
}

func (r *ketersediaanRepository) GetAvailabilitySummary(eventID uuid.UUID, occurrenceDate time.Time) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := r.db.Model(&entity.Ketersediaan{}).
		Select("status, COUNT(*) as count").
		Where("event_id = ? AND occurrence_date = ?", eventID, occurrenceDate).
		Group("status").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]int)
	summary["available"] = 0
	summary["unavailable"] = 0
	summary["tentative"] = 0

	for _, result := range results {
		summary[result.Status] = result.Count
	}

	return summary, nil
}
