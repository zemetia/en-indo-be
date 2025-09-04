package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupVisitorMemberRepository interface {
	Create(ctx context.Context, member *entity.LifeGroupVisitorMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeGroupVisitorMember, error)
	GetByLifeGroupID(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupVisitorMember, error)
	GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]entity.LifeGroupVisitorMember, error)
	GetByLifeGroupAndVisitorID(ctx context.Context, lifeGroupID uuid.UUID, visitorID uuid.UUID) (*entity.LifeGroupVisitorMember, error)
	Update(ctx context.Context, member *entity.LifeGroupVisitorMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByLifeGroupAndVisitorID(ctx context.Context, lifeGroupID uuid.UUID, visitorID uuid.UUID) (bool, error)
}

type lifeGroupVisitorMemberRepository struct {
	db *gorm.DB
}

func NewLifeGroupVisitorMemberRepository(db *gorm.DB) LifeGroupVisitorMemberRepository {
	return &lifeGroupVisitorMemberRepository{
		db: db,
	}
}

func (r *lifeGroupVisitorMemberRepository) Create(ctx context.Context, member *entity.LifeGroupVisitorMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *lifeGroupVisitorMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeGroupVisitorMember, error) {
	var member entity.LifeGroupVisitorMember
	err := r.db.WithContext(ctx).
		Preload("Visitor").
		Preload("Visitor.Kabupaten").
		Preload("LifeGroup").
		Where("id = ? AND is_active = ?", id, true).
		First(&member).Error
	return &member, err
}

func (r *lifeGroupVisitorMemberRepository) GetByLifeGroupID(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupVisitorMember, error) {
	var members []entity.LifeGroupVisitorMember
	err := r.db.WithContext(ctx).
		Preload("Visitor").
		Preload("Visitor.Kabupaten").
		Where("life_group_id = ? AND is_active = ?", lifeGroupID, true).
		Order("joined_date ASC").
		Find(&members).Error
	return members, err
}

func (r *lifeGroupVisitorMemberRepository) GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]entity.LifeGroupVisitorMember, error) {
	var members []entity.LifeGroupVisitorMember
	err := r.db.WithContext(ctx).
		Preload("LifeGroup").
		Preload("LifeGroup.Church").
		Where("visitor_id = ? AND is_active = ?", visitorID, true).
		Find(&members).Error
	return members, err
}

func (r *lifeGroupVisitorMemberRepository) GetByLifeGroupAndVisitorID(ctx context.Context, lifeGroupID uuid.UUID, visitorID uuid.UUID) (*entity.LifeGroupVisitorMember, error) {
	var member entity.LifeGroupVisitorMember
	err := r.db.WithContext(ctx).
		Preload("Visitor").
		Preload("Visitor.Kabupaten").
		Preload("LifeGroup").
		Where("life_group_id = ? AND visitor_id = ? AND is_active = ?", lifeGroupID, visitorID, true).
		First(&member).Error
	return &member, err
}

func (r *lifeGroupVisitorMemberRepository) Update(ctx context.Context, member *entity.LifeGroupVisitorMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *lifeGroupVisitorMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.LifeGroupVisitorMember{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *lifeGroupVisitorMemberRepository) ExistsByLifeGroupAndVisitorID(ctx context.Context, lifeGroupID uuid.UUID, visitorID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.LifeGroupVisitorMember{}).
		Where("life_group_id = ? AND visitor_id = ? AND is_active = ?", lifeGroupID, visitorID, true).
		Count(&count).Error
	return count > 0, err
}