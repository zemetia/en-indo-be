package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type EventDepartmentRepository interface {
	GetEventDepartments(ctx context.Context, eventID uuid.UUID) ([]entity.Department, error)
	UpdateEventDepartments(ctx context.Context, eventID uuid.UUID, departmentIDs []uuid.UUID) error
}

type eventDepartmentRepository struct {
	db *gorm.DB
}

func NewEventDepartmentRepository(db *gorm.DB) EventDepartmentRepository {
	return &eventDepartmentRepository{
		db: db,
	}
}

func (r *eventDepartmentRepository) GetEventDepartments(ctx context.Context, eventID uuid.UUID) ([]entity.Department, error) {
	var departments []entity.Department

	err := r.db.WithContext(ctx).
		Joins("JOIN event_departments ON event_departments.department_id = departments.id").
		Where("event_departments.event_id = ?", eventID).
		Find(&departments).Error

	if err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *eventDepartmentRepository) UpdateEventDepartments(ctx context.Context, eventID uuid.UUID, departmentIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all existing event-department associations
		if err := tx.Where("event_id = ?", eventID).Delete(&entity.EventDepartment{}).Error; err != nil {
			return err
		}

		// Create new associations
		for _, departmentID := range departmentIDs {
			eventDept := entity.EventDepartment{
				ID:           uuid.New(),
				EventID:      eventID,
				DepartmentID: departmentID,
			}
			if err := tx.Create(&eventDept).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
