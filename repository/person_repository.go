package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type PersonRepository interface {
	Create(ctx context.Context, person *entity.Person) error
	GetAll(ctx context.Context) ([]entity.Person, error)
	Search(ctx context.Context, search *dto.PersonSearchDto) ([]entity.Person, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Person, error)
	GetByChurchID(ctx context.Context, churchID uuid.UUID) ([]entity.Person, error)
	GetByKabupatenID(ctx context.Context, kabupatenID uuid.UUID) ([]entity.Person, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Person, error)
	Update(ctx context.Context, person *entity.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
	RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
}

type personRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{
		db: db,
	}
}

func (r *personRepository) Create(ctx context.Context, person *entity.Person) error {
	return r.db.WithContext(ctx).Create(person).Error
}

func (r *personRepository) GetAll(ctx context.Context) ([]entity.Person, error) {
	var persons []entity.Person
	err := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten").
		Find(&persons).Error
	return persons, err
}

func (r *personRepository) Search(ctx context.Context, search *dto.PersonSearchDto) ([]entity.Person, error) {
	var persons []entity.Person
	query := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten")

	if search.Name != nil {
		query = query.Where("nama LIKE ?", "%"+*search.Name+"%")
	}

	if search.ChurchID != nil {
		query = query.Where("church_id = ?", search.ChurchID)
	}

	if search.KabupatenID != nil {
		query = query.Where("kabupaten_id = ?", search.KabupatenID)
	}

	if search.UserID != nil {
		query = query.Where("user_id = ?", search.UserID)
	}

	err := query.Find(&persons).Error
	return persons, err
}

func (r *personRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Person, error) {
	var person entity.Person
	err := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten").
		First(&person, "id = ?", id).Error
	return &person, err
}

func (r *personRepository) GetByChurchID(ctx context.Context, churchID uuid.UUID) ([]entity.Person, error) {
	var persons []entity.Person
	err := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten").
		Where("church_id = ?", churchID).Find(&persons).Error
	return persons, err
}

func (r *personRepository) GetByKabupatenID(ctx context.Context, kabupatenID uuid.UUID) ([]entity.Person, error) {
	var persons []entity.Person
	err := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten").
		Where("kabupaten_id = ?", kabupatenID).Find(&persons).Error
	return persons, err
}

func (r *personRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Person, error) {
	var person entity.Person
	err := r.db.WithContext(ctx).
		Preload("Pasangan").
		Preload("Church").
		Preload("Kabupaten").
		Where("user_id = ?", userID).First(&person).Error
	return &person, err
}

func (r *personRepository) Update(ctx context.Context, person *entity.Person) error {
	return r.db.WithContext(ctx).Save(person).Error
}

func (r *personRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Person{}, "id = ?", id).Error
}

func (r *personRepository) AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
	var person entity.Person
	if err := r.db.WithContext(ctx).First(&person, "id = ?", personID).Error; err != nil {
		return err
	}

	var lifeGroup entity.LifeGroup
	if err := r.db.WithContext(ctx).First(&lifeGroup, "id = ?", lifeGroupID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&person).Association("LifeGroups").Append(&lifeGroup)
}

func (r *personRepository) RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
	var person entity.Person
	if err := r.db.WithContext(ctx).First(&person, "id = ?", personID).Error; err != nil {
		return err
	}

	var lifeGroup entity.LifeGroup
	if err := r.db.WithContext(ctx).First(&lifeGroup, "id = ?", lifeGroupID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&person).Association("LifeGroups").Delete(&lifeGroup)
}
