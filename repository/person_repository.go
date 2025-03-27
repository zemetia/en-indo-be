package repository

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type PersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *PersonRepository {
	return &PersonRepository{
		db: db,
	}
}

func (r *PersonRepository) Create(person *entity.Person) error {
	return r.db.Create(person).Error
}

func (r *PersonRepository) GetAll() ([]entity.Person, error) {
	var persons []entity.Person
	err := r.db.Find(&persons).Error
	return persons, err
}

func (r *PersonRepository) GetByID(id uint) (*entity.Person, error) {
	var person entity.Person
	err := r.db.First(&person, id).Error
	return &person, err
}

func (r *PersonRepository) Update(person *entity.Person) error {
	return r.db.Save(person).Error
}

func (r *PersonRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Person{}, id).Error
}

func (r *PersonRepository) GetByChurchID(churchID uint) ([]entity.Person, error) {
	var persons []entity.Person
	err := r.db.Where("church_id = ?", churchID).Find(&persons).Error
	return persons, err
}
