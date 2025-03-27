package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupRepository struct {
	db *gorm.DB
}

func NewLifeGroupRepository(db *gorm.DB) *LifeGroupRepository {
	return &LifeGroupRepository{
		db: db,
	}
}

func (r *LifeGroupRepository) Create(lifeGroup *entity.LifeGroup) error {
	return r.db.Create(lifeGroup).Error
}

func (r *LifeGroupRepository) GetAll() ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *LifeGroupRepository) GetByID(id uuid.UUID) (*entity.LifeGroup, error) {
	var lifeGroup entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").First(&lifeGroup, "id = ?", id).Error
	return &lifeGroup, err
}

func (r *LifeGroupRepository) Update(lifeGroup *entity.LifeGroup) error {
	return r.db.Save(lifeGroup).Error
}

func (r *LifeGroupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.LifeGroup{}, "id = ?", id).Error
}

func (r *LifeGroupRepository) UpdateLeader(id uuid.UUID, leaderID uuid.UUID) error {
	return r.db.Model(&entity.LifeGroup{}).Where("id = ?", id).Update("leader_id", leaderID).Error
}

func (r *LifeGroupRepository) UpdateMembers(id uuid.UUID, memberIDs []uuid.UUID) error {
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

func (r *LifeGroupRepository) UpdatePersons(id uuid.UUID, personIDs []uuid.UUID) error {
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

func (r *LifeGroupRepository) GetByChurchID(churchID uuid.UUID) ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("Leader").Preload("Members").Preload("Persons").
		Where("church_id = ?", churchID).
		Find(&lifeGroups).Error
	return lifeGroups, err
}
