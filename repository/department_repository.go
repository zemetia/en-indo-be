package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{
		db: db,
	}
}

func (r *DepartmentRepository) Create(department *entity.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) GetAll() ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) GetByID(id uuid.UUID) (*entity.Department, error) {
	var department entity.Department
	err := r.db.First(&department, "id = ?", id).Error
	return &department, err
}

func (r *DepartmentRepository) GetByChurchID(churchID uuid.UUID) ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.Where("church_id = ?", churchID).Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) Update(department *entity.Department) error {
	return r.db.Save(department).Error
}

func (r *DepartmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Department{}, "id = ?", id).Error
}
