package repository

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type KabupatenRepository interface {
	GetAll() ([]entity.Kabupaten, error)
	GetByID(id uint) (*entity.Kabupaten, error)
	GetByProvinsiID(provinsiID uint) ([]entity.Kabupaten, error)
}

type kabupatenRepository struct {
	db *gorm.DB
}

func NewKabupatenRepository(db *gorm.DB) KabupatenRepository {
	return &kabupatenRepository{
		db: db,
	}
}

func (r *kabupatenRepository) GetAll() ([]entity.Kabupaten, error) {
	var kabupaten []entity.Kabupaten
	err := r.db.Preload("Provinsi").Find(&kabupaten).Error
	return kabupaten, err
}

func (r *kabupatenRepository) GetByID(id uint) (*entity.Kabupaten, error) {
	var kabupaten entity.Kabupaten
	err := r.db.Preload("Provinsi").First(&kabupaten, "id = ?", id).Error
	return &kabupaten, err
}

func (r *kabupatenRepository) GetByProvinsiID(provinsiID uint) ([]entity.Kabupaten, error) {
	var kabupaten []entity.Kabupaten
	err := r.db.Preload("Provinsi").Where("provinsi_id = ?", provinsiID).Find(&kabupaten).Error
	return kabupaten, err
}
