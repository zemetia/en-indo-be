package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(department *entity.Department) error
	GetAll() ([]entity.Department, error)
	GetByID(id uuid.UUID) (*entity.Department, error)
	GetByChurchID(churchID uuid.UUID) ([]entity.Department, error)
	Update(department *entity.Department) error
	Delete(id uuid.UUID) error
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) Create(department *entity.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) GetAll() ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) GetByID(id uuid.UUID) (*entity.Department, error) {
	var department entity.Department
	err := r.db.First(&department, "id = ?", id).Error
	return &department, err
}

func (r *departmentRepository) GetByChurchID(churchID uuid.UUID) ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.Where("church_id = ?", churchID).Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) Update(department *entity.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Department{}, "id = ?", id).Error
}
