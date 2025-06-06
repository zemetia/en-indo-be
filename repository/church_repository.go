package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type ChurchRepository interface {
	Create(church *entity.Church) error
	GetAll() ([]entity.Church, error)
	GetByID(id uuid.UUID) (*entity.Church, error)
	GetByKabupatenID(kabupatenID uuid.UUID) ([]entity.Church, error)
	GetByProvinsiID(provinsiID uuid.UUID) ([]entity.Church, error)
	Update(church *entity.Church) error
	Delete(id uuid.UUID) error
}

type churchRepository struct {
	db *gorm.DB
}

func NewChurchRepository(db *gorm.DB) ChurchRepository {
	return &churchRepository{
		db: db,
	}
}

func (r *churchRepository) Create(church *entity.Church) error {
	return r.db.Create(church).Error
}

func (r *churchRepository) GetAll() ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Preload("Kabupaten").Preload("Provinsi").Find(&churches).Error
	return churches, err
}

func (r *churchRepository) GetByID(id uuid.UUID) (*entity.Church, error) {
	var church entity.Church
	err := r.db.Preload("Kabupaten").Preload("Provinsi").First(&church, "id = ?", id).Error
	return &church, err
}

func (r *churchRepository) GetByKabupatenID(kabupatenID uuid.UUID) ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Preload("Kabupaten").Preload("Provinsi").Where("kabupaten_id = ?", kabupatenID).Find(&churches).Error
	return churches, err
}

func (r *churchRepository) GetByProvinsiID(provinsiID uuid.UUID) ([]entity.Church, error) {
	var churches []entity.Church
	err := r.db.Preload("Kabupaten").Preload("Provinsi").Where("provinsi_id = ?", provinsiID).Find(&churches).Error
	return churches, err
}

func (r *churchRepository) Update(church *entity.Church) error {
	return r.db.Save(church).Error
}

func (r *churchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Church{}, "id = ?", id).Error
}
