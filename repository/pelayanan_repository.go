package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type PelayananRepository interface {
	// Pelayanan entity CRUD
	CreatePelayanan(ctx context.Context, pelayanan *entity.Pelayanan) error
	UpdatePelayanan(ctx context.Context, pelayanan *entity.Pelayanan) error
	DeletePelayanan(ctx context.Context, id uuid.UUID) error
	GetPelayananByID(ctx context.Context, id uuid.UUID) (*entity.Pelayanan, error)
	GetAllPelayanan(ctx context.Context) ([]entity.Pelayanan, error)
	GetAllPelayananByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entity.Pelayanan, error)
	GetPelayananByDepartmentAndPic(ctx context.Context, departmentID uuid.UUID, isPic bool) (*entity.Pelayanan, error)

	// Assignment operations
	GetPelayananByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.PersonPelayananGereja, error)
	GetAllPelayananAssignments(ctx context.Context, req dto.PaginationRequest) ([]entity.PersonPelayananGereja, *dto.PaginationResponse, error)
	CreatePelayananAssignment(ctx context.Context, assignment *entity.PersonPelayananGereja) error
	DeletePelayananAssignment(ctx context.Context, id uuid.UUID) error
	GetAssignmentByID(ctx context.Context, id uuid.UUID) (*entity.PersonPelayananGereja, error)
	GetAssignmentByPersonPelayananChurch(ctx context.Context, personID, pelayananID, churchID uuid.UUID) (*entity.PersonPelayananGereja, error)
	UpdatePelayananAssignment(ctx context.Context, assignment *entity.PersonPelayananGereja) error
}

type pelayananRepository struct {
	db *gorm.DB
}

func NewPelayananRepository(db *gorm.DB) PelayananRepository {
	return &pelayananRepository{
		db: db,
	}
}

// Pelayanan entity CRUD methods
func (r *pelayananRepository) CreatePelayanan(ctx context.Context, pelayanan *entity.Pelayanan) error {
	return r.db.WithContext(ctx).Create(pelayanan).Error
}

func (r *pelayananRepository) UpdatePelayanan(ctx context.Context, pelayanan *entity.Pelayanan) error {
	return r.db.WithContext(ctx).Save(pelayanan).Error
}

func (r *pelayananRepository) DeletePelayanan(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Pelayanan{}, id).Error
}

func (r *pelayananRepository) GetPelayananByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.PersonPelayananGereja, error) {
	var assignments []entity.PersonPelayananGereja

	if err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("Pelayanan").
		Preload("Pelayanan.Department").
		Preload("Church").
		Where("person_id = ?", personID).
		Find(&assignments).Error; err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *pelayananRepository) GetAllPelayananAssignments(ctx context.Context, req dto.PaginationRequest) ([]entity.PersonPelayananGereja, *dto.PaginationResponse, error) {
	var assignments []entity.PersonPelayananGereja
	var total int64

	// Set defaults
	req.Default()

	// Count total records
	if err := r.db.WithContext(ctx).Model(&entity.PersonPelayananGereja{}).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get paginated data
	if err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("Pelayanan").
		Preload("Pelayanan.Department").
		Preload("Church").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&assignments).Error; err != nil {
		return nil, nil, err
	}

	// Calculate pagination info
	maxPage := (total + int64(req.PerPage) - 1) / int64(req.PerPage)
	pagination := &dto.PaginationResponse{
		Page:    req.Page,
		PerPage: req.PerPage,
		MaxPage: maxPage,
		Count:   total,
	}

	return assignments, pagination, nil
}

func (r *pelayananRepository) GetPelayananByID(ctx context.Context, id uuid.UUID) (*entity.Pelayanan, error) {
	var pelayanan entity.Pelayanan

	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("id = ?", id).
		First(&pelayanan).Error; err != nil {
		return nil, err
	}

	return &pelayanan, nil
}

func (r *pelayananRepository) GetAllPelayanan(ctx context.Context) ([]entity.Pelayanan, error) {
	var pelayanan []entity.Pelayanan

	if err := r.db.WithContext(ctx).
		Preload("Department").
		Find(&pelayanan).Error; err != nil {
		return nil, err
	}

	return pelayanan, nil
}

func (r *pelayananRepository) GetAllPelayananByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entity.Pelayanan, error) {
	var pelayanan []entity.Pelayanan

	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("department_id = ?", departmentID).
		Find(&pelayanan).Error; err != nil {
		return nil, err
	}

	return pelayanan, nil
}

func (r *pelayananRepository) GetPelayananByDepartmentAndPic(ctx context.Context, departmentID uuid.UUID, isPic bool) (*entity.Pelayanan, error) {
	var pelayanan entity.Pelayanan

	if err := r.db.WithContext(ctx).
		Preload("Department").
		Where("department_id = ? AND is_pic = ?", departmentID, isPic).
		First(&pelayanan).Error; err != nil {
		return nil, err
	}

	return &pelayanan, nil
}

func (r *pelayananRepository) CreatePelayananAssignment(ctx context.Context, assignment *entity.PersonPelayananGereja) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *pelayananRepository) DeletePelayananAssignment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.PersonPelayananGereja{}, id).Error
}

func (r *pelayananRepository) GetAssignmentByID(ctx context.Context, id uuid.UUID) (*entity.PersonPelayananGereja, error) {
	var assignment entity.PersonPelayananGereja

	if err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("Pelayanan").
		Preload("Pelayanan.Department").
		Preload("Church").
		Where("id = ?", id).
		First(&assignment).Error; err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *pelayananRepository) GetAssignmentByPersonPelayananChurch(ctx context.Context, personID, pelayananID, churchID uuid.UUID) (*entity.PersonPelayananGereja, error) {
	var assignment entity.PersonPelayananGereja

	err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("Pelayanan").
		Preload("Pelayanan.Department").
		Preload("Church").
		Where("person_id = ? AND pelayanan_id = ? AND church_id = ?", personID, pelayananID, churchID).
		First(&assignment).Error

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *pelayananRepository) UpdatePelayananAssignment(ctx context.Context, assignment *entity.PersonPelayananGereja) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}
