package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type ChurchRepository struct {
	db *gorm.DB
}

func NewChurchRepository(db *gorm.DB) *ChurchRepository {
	return &ChurchRepository{
		db: db,
	}
}

func (r *ChurchRepository) Create(church *entity.Church) error {
	return r.db.Create(church).Error
}

func (r *ChurchRepository) GetAll() ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Find(&churches).Error
	return churches, err
}

func (r *ChurchRepository) GetByID(id uuid.UUID) (*entity.Church, error) {
	var church entity.Church
	err := r.db.First(&church, "id = ?", id).Error
	return &church, err
}

func (r *ChurchRepository) GetByCityID(cityID uuid.UUID) ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Where("city_id = ?", cityID).Find(&churches).Error
	return churches, err
}

func (r *ChurchRepository) GetByProvinceID(provinceID uuid.UUID) ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Where("province_id = ?", provinceID).Find(&churches).Error
	return churches, err
}

func (r *ChurchRepository) Update(church *entity.Church) error {
	return r.db.Save(church).Error
}

func (r *ChurchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Church{}, "id = ?", id).Error
}
