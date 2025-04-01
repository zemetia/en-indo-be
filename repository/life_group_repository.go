package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupRepository interface {
	Create(lifeGroup *entity.LifeGroup) error
	GetAll() ([]entity.LifeGroup, error)
	GetByID(id uuid.UUID) (*entity.LifeGroup, error)
	Update(lifeGroup *entity.LifeGroup) error
	Delete(id uuid.UUID) error
	UpdateLeader(id uuid.UUID, leaderID uuid.UUID) error
	UpdateMembers(id uuid.UUID, memberIDs []uuid.UUID) error
	AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
	RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error
}

type lifeGroupRepository struct {
	db *gorm.DB
}

func NewLifeGroupRepository(db *gorm.DB) LifeGroupRepository {
	return &lifeGroupRepository{
		db: db,
	}
}

func (r *lifeGroupRepository) Create(lifeGroup *entity.LifeGroup) error {
	return r.db.Create(lifeGroup).Error
}

func (r *lifeGroupRepository) GetAll() ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *lifeGroupRepository) GetByID(id uuid.UUID) (*entity.LifeGroup, error) {
	var lifeGroup entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").First(&lifeGroup, "id = ?", id).Error
	return &lifeGroup, err
}

func (r *lifeGroupRepository) Update(lifeGroup *entity.LifeGroup) error {
	return r.db.Save(lifeGroup).Error
}

func (r *lifeGroupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.LifeGroup{}, "id = ?", id).Error
}

func (r *lifeGroupRepository) UpdateLeader(id uuid.UUID, leaderID uuid.UUID) error {
	return r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Update("leader_id", leaderID).Error
}

func (r *lifeGroupRepository) UpdateMembers(id uuid.UUID, memberIDs []uuid.UUID) error {
	// Hapus semua relasi member yang ada
	if err := r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Association("Members").Clear(); err != nil {
		return err
	}

	// Tambahkan member baru
	var members []entity.User
	if err := r.db.Find(&members, memberIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Association("Members").Replace(members)
}

func (r *lifeGroupRepository) UpdatePersons(id uuid.UUID, personIDs []uuid.UUID) error {
	// Hapus semua relasi person yang ada
	if err := r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Association("Persons").Clear(); err != nil {
		return err
	}

	// Tambahkan person baru
	var persons []entity.Person
	if err := r.db.Find(&persons, personIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Association("Persons").Replace(persons)
}

func (r *lifeGroupRepository) GetByChurchID(churchID uuid.UUID) ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").
		Where("church_id = ?", churchID).
		Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *lifeGroupRepository) AddToLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
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

func (r *lifeGroupRepository) RemoveFromLifeGroup(ctx context.Context, personID uuid.UUID, lifeGroupID uuid.UUID) error {
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
