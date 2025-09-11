package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *entity.Event) error
	GetByID(id uuid.UUID) (*entity.Event, error)
	Update(event *entity.Event) error
	Delete(id uuid.UUID) error
	List(filters EventFilters) ([]entity.Event, int64, error)
	GetByDateRange(startDate, endDate time.Time, filters EventFilters) ([]entity.Event, error)

	// Recurrence related methods
	CreateRecurrenceRule(rule *entity.RecurrenceRule) error
	UpdateRecurrenceRule(rule *entity.RecurrenceRule) error
	DeleteRecurrenceRule(id uuid.UUID) error

	// Exception handling
	CreateRecurrenceException(exception *entity.RecurrenceException) error
	GetRecurrenceExceptions(eventID uuid.UUID) ([]entity.RecurrenceException, error)
	UpdateRecurrenceException(exception *entity.RecurrenceException) error
	DeleteRecurrenceException(id uuid.UUID) error
	GetExceptionByEventAndDate(eventID uuid.UUID, date time.Time) (*entity.RecurrenceException, error)

	// Bulk operations for recurring events
	DeleteFutureOccurrences(eventID uuid.UUID, fromDate time.Time) error
	SetRecurrenceUntilDate(eventID uuid.UUID, untilDate time.Time) error
	GetEventsWithRecurrenceInRange(startDate, endDate time.Time) ([]entity.Event, error)
}

type EventFilters struct {
	Type     string
	IsPublic *bool
	Search   string
	Limit    int
	Offset   int
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *entity.Event) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create recurrence rule first if exists
		if event.RecurrenceRule != nil {
			event.RecurrenceRule.ID = uuid.New()
			if err := tx.Create(event.RecurrenceRule).Error; err != nil {
				return err
			}
			event.RecurrenceRuleID = &event.RecurrenceRule.ID
		}

		// Create the event
		event.ID = uuid.New()
		return tx.Create(event).Error
	})
}

func (r *eventRepository) GetByID(id uuid.UUID) (*entity.Event, error) {
	var event entity.Event
	err := r.db.Preload("RecurrenceRule").
		Preload("Lagu").
		Preload("DiscipleshipJourney").
		First(&event, id).Error

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Update(event *entity.Event) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Handle recurrence rule update
		if event.RecurrenceRule != nil {
			if event.RecurrenceRuleID == nil {
				// Create new recurrence rule
				event.RecurrenceRule.ID = uuid.New()
				if err := tx.Create(event.RecurrenceRule).Error; err != nil {
					return err
				}
				event.RecurrenceRuleID = &event.RecurrenceRule.ID
			} else {
				// Update existing recurrence rule
				if err := tx.Save(event.RecurrenceRule).Error; err != nil {
					return err
				}
			}
		} else if event.RecurrenceRuleID != nil {
			// Remove recurrence rule if it was cleared
			if err := tx.Delete(&entity.RecurrenceRule{}, event.RecurrenceRuleID).Error; err != nil {
				return err
			}
			event.RecurrenceRuleID = nil
		}

		// Update the event
		return tx.Save(event).Error
	})
}

func (r *eventRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get event with recurrence rule
		var event entity.Event
		if err := tx.Preload("RecurrenceRule").First(&event, id).Error; err != nil {
			return err
		}

		// Delete recurrence exceptions first
		if err := tx.Where("event_id = ?", id).Delete(&entity.RecurrenceException{}).Error; err != nil {
			return err
		}

		// Delete recurrence rule if exists
		if event.RecurrenceRuleID != nil {
			if err := tx.Delete(&entity.RecurrenceRule{}, event.RecurrenceRuleID).Error; err != nil {
				return err
			}
		}

		// Delete the event
		return tx.Delete(&event).Error
	})
}

func (r *eventRepository) List(filters EventFilters) ([]entity.Event, int64, error) {
	var events []entity.Event
	var count int64

	query := r.db.Model(&entity.Event{}).
		Preload("RecurrenceRule").
		Preload("Lagu").
		Preload("DiscipleshipJourney")

	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.Search != "" {
		query = query.Where(
			"title ILIKE ? OR description ILIKE ? OR event_location ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%",
		)
	}

	// Count total records
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Order("event_date ASC, start_datetime ASC").Find(&events).Error
	return events, count, err
}

func (r *eventRepository) GetByDateRange(startDate, endDate time.Time, filters EventFilters) ([]entity.Event, error) {
	var events []entity.Event

	query := r.db.Model(&entity.Event{}).
		Preload("RecurrenceRule").
		Preload("Lagu").
		Preload("DiscipleshipJourney").
		Where("event_date >= ? AND event_date <= ?", startDate, endDate)

	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.Search != "" {
		query = query.Where(
			"title ILIKE ? OR description ILIKE ? OR event_location ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%",
		)
	}

	err := query.Order("event_date ASC, start_datetime ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) CreateRecurrenceRule(rule *entity.RecurrenceRule) error {
	rule.ID = uuid.New()
	return r.db.Create(rule).Error
}

func (r *eventRepository) UpdateRecurrenceRule(rule *entity.RecurrenceRule) error {
	return r.db.Save(rule).Error
}

func (r *eventRepository) DeleteRecurrenceRule(id uuid.UUID) error {
	return r.db.Delete(&entity.RecurrenceRule{}, id).Error
}

func (r *eventRepository) CreateRecurrenceException(exception *entity.RecurrenceException) error {
	exception.ID = uuid.New()
	return r.db.Create(exception).Error
}

func (r *eventRepository) GetRecurrenceExceptions(eventID uuid.UUID) ([]entity.RecurrenceException, error) {
	var exceptions []entity.RecurrenceException
	err := r.db.Where("event_id = ?", eventID).
		Order("exception_date ASC").
		Find(&exceptions).Error
	return exceptions, err
}

func (r *eventRepository) UpdateRecurrenceException(exception *entity.RecurrenceException) error {
	return r.db.Save(exception).Error
}

func (r *eventRepository) DeleteRecurrenceException(id uuid.UUID) error {
	return r.db.Delete(&entity.RecurrenceException{}, id).Error
}

func (r *eventRepository) GetExceptionByEventAndDate(eventID uuid.UUID, date time.Time) (*entity.RecurrenceException, error) {
	var exception entity.RecurrenceException
	err := r.db.Where("event_id = ? AND exception_date = ?", eventID, date.Format("2006-01-02")).
		First(&exception).Error

	if err != nil {
		return nil, err
	}
	return &exception, nil
}

func (r *eventRepository) DeleteFutureOccurrences(eventID uuid.UUID, fromDate time.Time) error {
	// This creates exceptions to skip all future occurrences from the given date
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get the event with recurrence rule
		var event entity.Event
		if err := tx.Preload("RecurrenceRule").First(&event, eventID).Error; err != nil {
			return err
		}

		if event.RecurrenceRule == nil {
			return nil // No recurrence to handle
		}

		// Update recurrence rule to end at the given date
		yesterday := fromDate.AddDate(0, 0, -1)
		event.RecurrenceRule.Until = &yesterday

		return tx.Save(event.RecurrenceRule).Error
	})
}

func (r *eventRepository) SetRecurrenceUntilDate(eventID uuid.UUID, untilDate time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get the event with recurrence rule
		var event entity.Event
		if err := tx.Preload("RecurrenceRule").First(&event, eventID).Error; err != nil {
			return err
		}

		if event.RecurrenceRule == nil {
			return nil // No recurrence to handle
		}

		// Update recurrence rule to end at the given date
		event.RecurrenceRule.Until = &untilDate

		return tx.Save(event.RecurrenceRule).Error
	})
}

func (r *eventRepository) GetEventsWithRecurrenceInRange(startDate, endDate time.Time) ([]entity.Event, error) {
	var events []entity.Event

	// Get all recurring events that could have occurrences in the date range
	// This includes events that start before the range but have recurrence rules
	err := r.db.Preload("RecurrenceRule").
		Preload("Lagu").
		Preload("DiscipleshipJourney").
		Joins("LEFT JOIN recurrence_rules ON events.recurrence_rule_id = recurrence_rules.id").
		Where(`
			(events.recurrence_rule_id IS NOT NULL) AND 
			(events.event_date <= ? OR recurrence_rules.until IS NULL OR recurrence_rules.until >= ?)
		`, endDate, startDate).
		Find(&events).Error

	return events, err
}
