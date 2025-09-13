package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupRepository interface {
	Create(lifeGroup *entity.LifeGroup) error
	GetAll() ([]entity.LifeGroup, error)
	GetByID(id uuid.UUID) (*entity.LifeGroup, error)
	Update(lifeGroup *entity.LifeGroup) error
	Delete(id uuid.UUID) error
	Search(ctx context.Context, search *dto.PersonSearchDto) ([]entity.LifeGroup, error)
	UpdateLeader(id uuid.UUID, leaderID uuid.UUID) error
	GetByChurchID(churchID uuid.UUID) ([]entity.LifeGroup, error)
	GetByUserID(userID uuid.UUID) ([]entity.LifeGroup, error)
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
	err := r.db.Preload("Church").Preload("PersonMembers").Preload("PersonMembers.Person").Preload("VisitorMembers").Preload("VisitorMembers.Visitor").Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *lifeGroupRepository) Search(ctx context.Context, search *dto.PersonSearchDto) ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	query := r.db.Preload("Church").Preload("PersonMembers").Preload("PersonMembers.Person").Preload("VisitorMembers").Preload("VisitorMembers.Visitor")

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

	err := query.Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *lifeGroupRepository) GetByID(id uuid.UUID) (*entity.LifeGroup, error) {
	var lifeGroup entity.LifeGroup
	err := r.db.Preload("Church").Preload("PersonMembers").Preload("PersonMembers.Person").Preload("VisitorMembers").Preload("VisitorMembers.Visitor").First(&lifeGroup, "id = ?", id).Error
	return &lifeGroup, err
}

func (r *lifeGroupRepository) Update(lifeGroup *entity.LifeGroup) error {
	return r.db.Save(lifeGroup).Error
}

func (r *lifeGroupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.LifeGroup{}, "id = ?", id).Error
}

func (r *lifeGroupRepository) UpdateLeader(id uuid.UUID, leaderID uuid.UUID) error {
	// Leader management is now handled through PersonMember API
	// This method is deprecated
	return nil
}

func (r *lifeGroupRepository) GetByChurchID(churchID uuid.UUID) ([]entity.LifeGroup, error) {
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("PersonMembers").Preload("PersonMembers.Person").Preload("VisitorMembers").Preload("VisitorMembers.Visitor").
		Where("church_id = ?", churchID).
		Find(&lifeGroups).Error
	return lifeGroups, err
}

func (r *lifeGroupRepository) GetByUserID(userID uuid.UUID) ([]entity.LifeGroup, error) {
	// This method now searches through PersonMembers table instead of leader columns
	var lifeGroups []entity.LifeGroup
	err := r.db.Preload("Church").Preload("PersonMembers").Preload("PersonMembers.Person").Preload("VisitorMembers").Preload("VisitorMembers.Visitor").
		Joins("JOIN life_group_person_members lgpm ON lgpm.life_group_id = life_groups.id").
		Joins("JOIN persons p ON p.id = lgpm.person_id").
		Joins("JOIN users u ON u.person_id = p.id").
		Where("u.id = ? AND lgpm.is_active = true", userID).
		Find(&lifeGroups).Error
	return lifeGroups, err
}
