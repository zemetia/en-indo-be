package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type EventPICRepository interface {
	// Basic CRUD operations
	Create(eventPIC *entity.EventPIC) error
	GetByID(id uuid.UUID) (*entity.EventPIC, error)
	Update(eventPIC *entity.EventPIC) error
	Delete(id uuid.UUID) error
	
	// Event-specific PIC operations
	GetPICsByEventID(eventID uuid.UUID) ([]entity.EventPIC, error)
	GetPrimaryPICByEventID(eventID uuid.UUID) (*entity.EventPIC, error)
	GetActivePICsByEventID(eventID uuid.UUID) ([]entity.EventPIC, error)
	
	// Person-specific PIC operations
	GetPICsByPersonID(personID uuid.UUID) ([]entity.EventPIC, error)
	GetActivePICsByPersonID(personID uuid.UUID) ([]entity.EventPIC, error)
	
	// Validation and business logic support
	HasPersonPICRoleForEvent(eventID, personID uuid.UUID) (bool, error)
	IsPrimaryPICForEvent(eventID, personID uuid.UUID) (bool, error)
	CountActivePICsForEvent(eventID uuid.UUID) (int64, error)
	CountPrimaryPICsForEvent(eventID uuid.UUID) (int64, error)
	
	// Advanced queries
	List(filters EventPICFilters) ([]entity.EventPIC, int64, error)
	GetPICsWithPermission(eventID uuid.UUID, permission string) ([]entity.EventPIC, error)
	GetExpiringPICs(days int) ([]entity.EventPIC, error)
	
	// Bulk operations
	CreateMultiple(eventPICs []entity.EventPIC) error
	UpdateMultiple(eventPICs []entity.EventPIC) error
	DeactivateByEventID(eventID uuid.UUID) error
	TransferPICRole(fromPersonID, toPersonID, eventID uuid.UUID, transferType string) error
	
	// PIC Role management
	CreateRole(role *entity.EventPICRole) error
	GetRoleByID(id uuid.UUID) (*entity.EventPICRole, error)
	GetRoleByName(name string) (*entity.EventPICRole, error)
	UpdateRole(role *entity.EventPICRole) error
	DeleteRole(id uuid.UUID) error
	ListRoles(filters EventPICRoleFilters) ([]entity.EventPICRole, int64, error)
	
	// History and audit
	CreateHistory(history *entity.EventPICHistory) error
	GetHistoryByEventID(eventID uuid.UUID) ([]entity.EventPICHistory, error)
	GetHistoryByPersonID(personID uuid.UUID) ([]entity.EventPICHistory, error)
	ListHistory(filters EventPICHistoryFilters) ([]entity.EventPICHistory, int64, error)
}

type EventPICFilters struct {
	EventID   *uuid.UUID
	PersonID  *uuid.UUID
	Role      string
	IsActive  *bool
	IsPrimary *bool
	Search    string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

type EventPICRoleFilters struct {
	Name     string
	IsActive *bool
	Search   string
	Limit    int
	Offset   int
}

type EventPICHistoryFilters struct {
	EventID   *uuid.UUID
	PersonID  *uuid.UUID
	Action    string
	ChangedBy *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

type eventPICRepository struct {
	db *gorm.DB
}

func NewEventPICRepository(db *gorm.DB) EventPICRepository {
	return &eventPICRepository{db: db}
}

// Basic CRUD operations
func (r *eventPICRepository) Create(eventPIC *entity.EventPIC) error {
	return r.db.Create(eventPIC).Error
}

func (r *eventPICRepository) GetByID(id uuid.UUID) (*entity.EventPIC, error) {
	var eventPIC entity.EventPIC
	err := r.db.Preload("Event").Preload("Person").First(&eventPIC, id).Error
	if err != nil {
		return nil, err
	}
	return &eventPIC, nil
}

func (r *eventPICRepository) Update(eventPIC *entity.EventPIC) error {
	return r.db.Save(eventPIC).Error
}

func (r *eventPICRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.EventPIC{}, id).Error
}

// Event-specific PIC operations
func (r *eventPICRepository) GetPICsByEventID(eventID uuid.UUID) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	err := r.db.Preload("Person").
		Where("event_id = ?", eventID).
		Order("is_primary DESC, created_at ASC").
		Find(&pics).Error
	return pics, err
}

func (r *eventPICRepository) GetPrimaryPICByEventID(eventID uuid.UUID) (*entity.EventPIC, error) {
	var pic entity.EventPIC
	err := r.db.Preload("Person").
		Where("event_id = ? AND is_primary = ? AND is_active = ?", eventID, true, true).
		First(&pic).Error
	if err != nil {
		return nil, err
	}
	return &pic, nil
}

func (r *eventPICRepository) GetActivePICsByEventID(eventID uuid.UUID) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	err := r.db.Preload("Person").
		Where("event_id = ? AND is_active = ?", eventID, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Order("is_primary DESC, created_at ASC").
		Find(&pics).Error
	return pics, err
}

// Person-specific PIC operations
func (r *eventPICRepository) GetPICsByPersonID(personID uuid.UUID) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	err := r.db.Preload("Event").
		Where("person_id = ?", personID).
		Order("created_at DESC").
		Find(&pics).Error
	return pics, err
}

func (r *eventPICRepository) GetActivePICsByPersonID(personID uuid.UUID) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	err := r.db.Preload("Event").
		Where("person_id = ? AND is_active = ?", personID, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Order("created_at DESC").
		Find(&pics).Error
	return pics, err
}

// Validation and business logic support
func (r *eventPICRepository) HasPersonPICRoleForEvent(eventID, personID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.EventPIC{}).
		Where("event_id = ? AND person_id = ? AND is_active = ?", eventID, personID, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Count(&count).Error
	return count > 0, err
}

func (r *eventPICRepository) IsPrimaryPICForEvent(eventID, personID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.EventPIC{}).
		Where("event_id = ? AND person_id = ? AND is_primary = ? AND is_active = ?", 
			eventID, personID, true, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Count(&count).Error
	return count > 0, err
}

func (r *eventPICRepository) CountActivePICsForEvent(eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.EventPIC{}).
		Where("event_id = ? AND is_active = ?", eventID, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Count(&count).Error
	return count, err
}

func (r *eventPICRepository) CountPrimaryPICsForEvent(eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.EventPIC{}).
		Where("event_id = ? AND is_primary = ? AND is_active = ?", eventID, true, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02")).
		Count(&count).Error
	return count, err
}

// Advanced queries
func (r *eventPICRepository) List(filters EventPICFilters) ([]entity.EventPIC, int64, error) {
	var pics []entity.EventPIC
	var count int64

	query := r.db.Model(&entity.EventPIC{}).
		Preload("Event").
		Preload("Person")

	// Apply filters
	if filters.EventID != nil {
		query = query.Where("event_id = ?", *filters.EventID)
	}
	if filters.PersonID != nil {
		query = query.Where("person_id = ?", *filters.PersonID)
	}
	if filters.Role != "" {
		query = query.Where("role ILIKE ?", "%"+filters.Role+"%")
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filters.IsPrimary)
	}
	if filters.Search != "" {
		query = query.Joins("JOIN people ON event_pics.person_id = people.id").
			Where("people.nama ILIKE ? OR event_pics.role ILIKE ? OR event_pics.description ILIKE ?",
				"%"+filters.Search+"%", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}
	if filters.StartDate != nil {
		query = query.Where("start_date >= ?", filters.StartDate.Format("2006-01-02"))
	}
	if filters.EndDate != nil {
		query = query.Where("end_date IS NULL OR end_date <= ?", filters.EndDate.Format("2006-01-02"))
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

	err := query.Order("is_primary DESC, created_at DESC").Find(&pics).Error
	return pics, count, err
}

func (r *eventPICRepository) GetPICsWithPermission(eventID uuid.UUID, permission string) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	query := r.db.Preload("Person").
		Where("event_id = ? AND is_active = ?", eventID, true).
		Where("end_date IS NULL OR end_date >= ?", time.Now().Format("2006-01-02"))

	switch permission {
	case "edit":
		query = query.Where("can_edit = ?", true)
	case "delete":
		query = query.Where("can_delete = ?", true)
	case "assign_pic":
		query = query.Where("can_assign_pic = ?", true)
	}

	err := query.Find(&pics).Error
	return pics, err
}

func (r *eventPICRepository) GetExpiringPICs(days int) ([]entity.EventPIC, error) {
	var pics []entity.EventPIC
	cutoffDate := time.Now().AddDate(0, 0, days)
	
	err := r.db.Preload("Event").Preload("Person").
		Where("is_active = ? AND end_date IS NOT NULL AND end_date <= ?", true, cutoffDate.Format("2006-01-02")).
		Find(&pics).Error
	return pics, err
}

// Bulk operations
func (r *eventPICRepository) CreateMultiple(eventPICs []entity.EventPIC) error {
	return r.db.Create(&eventPICs).Error
}

func (r *eventPICRepository) UpdateMultiple(eventPICs []entity.EventPIC) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, pic := range eventPICs {
			if err := tx.Save(&pic).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *eventPICRepository) DeactivateByEventID(eventID uuid.UUID) error {
	return r.db.Model(&entity.EventPIC{}).
		Where("event_id = ?", eventID).
		Update("is_active", false).Error
}

func (r *eventPICRepository) TransferPICRole(fromPersonID, toPersonID, eventID uuid.UUID, transferType string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get the current PIC assignment
		var fromPIC entity.EventPIC
		if err := tx.Where("event_id = ? AND person_id = ? AND is_active = ?", 
			eventID, fromPersonID, true).First(&fromPIC).Error; err != nil {
			return err
		}

		switch transferType {
		case "replace":
			// Deactivate the old PIC
			fromPIC.IsActive = false
			endDate := time.Now()
			fromPIC.EndDate = &endDate
			if err := tx.Save(&fromPIC).Error; err != nil {
				return err
			}

			// Create new PIC with the same role and permissions
			newPIC := entity.EventPIC{
				EventID:           eventID,
				PersonID:          toPersonID,
				Role:              fromPIC.Role,
				Description:       fromPIC.Description,
				IsActive:          true,
				IsPrimary:         fromPIC.IsPrimary,
				StartDate:         time.Now(),
				CanEdit:           fromPIC.CanEdit,
				CanDelete:         fromPIC.CanDelete,
				CanAssignPIC:      fromPIC.CanAssignPIC,
				NotifyOnChanges:   fromPIC.NotifyOnChanges,
				NotifyOnReminders: fromPIC.NotifyOnReminders,
			}
			return tx.Create(&newPIC).Error

		case "add_as_secondary":
			// Keep the old PIC, but make them secondary if they were primary
			if fromPIC.IsPrimary {
				fromPIC.IsPrimary = false
				if err := tx.Save(&fromPIC).Error; err != nil {
					return err
				}
			}

			// Create new primary PIC
			newPIC := entity.EventPIC{
				EventID:           eventID,
				PersonID:          toPersonID,
				Role:              fromPIC.Role,
				Description:       fromPIC.Description,
				IsActive:          true,
				IsPrimary:         true,
				StartDate:         time.Now(),
				CanEdit:           fromPIC.CanEdit,
				CanDelete:         fromPIC.CanDelete,
				CanAssignPIC:      fromPIC.CanAssignPIC,
				NotifyOnChanges:   true,
				NotifyOnReminders: true,
			}
			return tx.Create(&newPIC).Error
		}

		return nil
	})
}

// PIC Role management
func (r *eventPICRepository) CreateRole(role *entity.EventPICRole) error {
	return r.db.Create(role).Error
}

func (r *eventPICRepository) GetRoleByID(id uuid.UUID) (*entity.EventPICRole, error) {
	var role entity.EventPICRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *eventPICRepository) GetRoleByName(name string) (*entity.EventPICRole, error) {
	var role entity.EventPICRole
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *eventPICRepository) UpdateRole(role *entity.EventPICRole) error {
	return r.db.Save(role).Error
}

func (r *eventPICRepository) DeleteRole(id uuid.UUID) error {
	return r.db.Delete(&entity.EventPICRole{}, id).Error
}

func (r *eventPICRepository) ListRoles(filters EventPICRoleFilters) ([]entity.EventPICRole, int64, error) {
	var roles []entity.EventPICRole
	var count int64

	query := r.db.Model(&entity.EventPICRole{})

	// Apply filters
	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+filters.Search+"%", "%"+filters.Search+"%")
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

	err := query.Order("name ASC").Find(&roles).Error
	return roles, count, err
}

// History and audit
func (r *eventPICRepository) CreateHistory(history *entity.EventPICHistory) error {
	return r.db.Create(history).Error
}

func (r *eventPICRepository) GetHistoryByEventID(eventID uuid.UUID) ([]entity.EventPICHistory, error) {
	var history []entity.EventPICHistory
	err := r.db.Preload("Person").Preload("ChangedByPerson").
		Where("event_id = ?", eventID).
		Order("action_date DESC").
		Find(&history).Error
	return history, err
}

func (r *eventPICRepository) GetHistoryByPersonID(personID uuid.UUID) ([]entity.EventPICHistory, error) {
	var history []entity.EventPICHistory
	err := r.db.Preload("Event").Preload("ChangedByPerson").
		Where("person_id = ?", personID).
		Order("action_date DESC").
		Find(&history).Error
	return history, err
}

func (r *eventPICRepository) ListHistory(filters EventPICHistoryFilters) ([]entity.EventPICHistory, int64, error) {
	var history []entity.EventPICHistory
	var count int64

	query := r.db.Model(&entity.EventPICHistory{}).
		Preload("Event").
		Preload("Person").
		Preload("ChangedByPerson")

	// Apply filters
	if filters.EventID != nil {
		query = query.Where("event_id = ?", *filters.EventID)
	}
	if filters.PersonID != nil {
		query = query.Where("person_id = ?", *filters.PersonID)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.ChangedBy != nil {
		query = query.Where("changed_by = ?", *filters.ChangedBy)
	}
	if filters.StartDate != nil {
		query = query.Where("action_date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("action_date <= ?", *filters.EndDate)
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

	err := query.Order("action_date DESC").Find(&history).Error
	return history, count, err
}