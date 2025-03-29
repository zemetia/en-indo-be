package repository

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type ProvinsiRepository interface {
	GetAll() ([]entity.Provinsi, error)
	GetByID(id uint) (*entity.Provinsi, error)
}

type provinsiRepository struct {
	db *gorm.DB
}

func NewProvinsiRepository(db *gorm.DB) ProvinsiRepository {
	return &provinsiRepository{
		db: db,
	}
}

func (r *provinsiRepository) GetAll() ([]entity.Provinsi, error) {
	var provinsi []entity.Provinsi
	err := r.db.Find(&provinsi).Error
	return provinsi, err
}

func (r *provinsiRepository) GetByID(id uint) (*entity.Provinsi, error) {
	var provinsi entity.Provinsi
	err := r.db.First(&provinsi, "id = ?", id).Error
	return &provinsi, err
}
