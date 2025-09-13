package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LifeGroupPersonMemberRepository interface {
	Create(ctx context.Context, member *entity.LifeGroupPersonMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeGroupPersonMember, error)
	GetByLifeGroupID(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupPersonMember, error)
	GetByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.LifeGroupPersonMember, error)
	GetByLifeGroupAndPersonID(ctx context.Context, lifeGroupID uuid.UUID, personID uuid.UUID) (*entity.LifeGroupPersonMember, error)
	Update(ctx context.Context, member *entity.LifeGroupPersonMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePosition(ctx context.Context, lifeGroupID uuid.UUID, personID uuid.UUID, position entity.PersonMemberPosition) error
	GetCurrentLeader(ctx context.Context, lifeGroupID uuid.UUID) (*entity.LifeGroupPersonMember, error)
	GetCoLeaders(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupPersonMember, error)
	CountByPosition(ctx context.Context, lifeGroupID uuid.UUID, position entity.PersonMemberPosition) (int64, error)
	DemoteCurrentLeader(ctx context.Context, lifeGroupID uuid.UUID) error
}

type lifeGroupPersonMemberRepository struct {
	db *gorm.DB
}

func NewLifeGroupPersonMemberRepository(db *gorm.DB) LifeGroupPersonMemberRepository {
	return &lifeGroupPersonMemberRepository{
		db: db,
	}
}

func (r *lifeGroupPersonMemberRepository) Create(ctx context.Context, member *entity.LifeGroupPersonMember) error {
	// Check if position is LEADER and there's already a leader
	if member.Position == entity.PersonMemberPositionLeader {
		count, err := r.CountByPosition(ctx, member.LifeGroupID, entity.PersonMemberPositionLeader)
		if err != nil {
			return err
		}
		if count > 0 {
			// Demote existing leader to CO_LEADER
			if err := r.DemoteCurrentLeader(ctx, member.LifeGroupID); err != nil {
				return err
			}
		}
	}

	return r.db.WithContext(ctx).Create(member).Error
}

func (r *lifeGroupPersonMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeGroupPersonMember, error) {
	var member entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("LifeGroup").
		Where("id = ? AND is_active = ?", id, true).
		First(&member).Error
	return &member, err
}

func (r *lifeGroupPersonMemberRepository) GetByLifeGroupID(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupPersonMember, error) {
	var members []entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("life_group_id = ? AND is_active = ?", lifeGroupID, true).
		Order("position ASC, joined_date ASC").
		Find(&members).Error
	return members, err
}

func (r *lifeGroupPersonMemberRepository) GetByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.LifeGroupPersonMember, error) {
	var members []entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("LifeGroup").
		Where("person_id = ? AND is_active = ?", personID, true).
		Find(&members).Error
	return members, err
}

func (r *lifeGroupPersonMemberRepository) GetByLifeGroupAndPersonID(ctx context.Context, lifeGroupID uuid.UUID, personID uuid.UUID) (*entity.LifeGroupPersonMember, error) {
	var member entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("Person").
		Preload("LifeGroup").
		Where("life_group_id = ? AND person_id = ? AND is_active = ?", lifeGroupID, personID, true).
		First(&member).Error
	return &member, err
}

func (r *lifeGroupPersonMemberRepository) Update(ctx context.Context, member *entity.LifeGroupPersonMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *lifeGroupPersonMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.LifeGroupPersonMember{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *lifeGroupPersonMemberRepository) UpdatePosition(ctx context.Context, lifeGroupID uuid.UUID, personID uuid.UUID, position entity.PersonMemberPosition) error {
	// Check if trying to set LEADER position
	if position == entity.PersonMemberPositionLeader {
		count, err := r.CountByPosition(ctx, lifeGroupID, entity.PersonMemberPositionLeader)
		if err != nil {
			return err
		}
		if count > 0 {
			// Demote existing leader to CO_LEADER
			if err := r.DemoteCurrentLeader(ctx, lifeGroupID); err != nil {
				return err
			}
		}
	}

	return r.db.WithContext(ctx).
		Model(&entity.LifeGroupPersonMember{}).
		Where("life_group_id = ? AND person_id = ? AND is_active = ?", lifeGroupID, personID, true).
		Update("position", position).Error
}

func (r *lifeGroupPersonMemberRepository) GetCurrentLeader(ctx context.Context, lifeGroupID uuid.UUID) (*entity.LifeGroupPersonMember, error) {
	var member entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("life_group_id = ? AND position = ? AND is_active = ?", lifeGroupID, entity.PersonMemberPositionLeader, true).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *lifeGroupPersonMemberRepository) GetCoLeaders(ctx context.Context, lifeGroupID uuid.UUID) ([]entity.LifeGroupPersonMember, error) {
	var members []entity.LifeGroupPersonMember
	err := r.db.WithContext(ctx).
		Preload("Person").
		Where("life_group_id = ? AND position = ? AND is_active = ?", lifeGroupID, entity.PersonMemberPositionCoLeader, true).
		Find(&members).Error
	return members, err
}

func (r *lifeGroupPersonMemberRepository) CountByPosition(ctx context.Context, lifeGroupID uuid.UUID, position entity.PersonMemberPosition) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.LifeGroupPersonMember{}).
		Where("life_group_id = ? AND position = ? AND is_active = ?", lifeGroupID, position, true).
		Count(&count).Error
	return count, err
}

func (r *lifeGroupPersonMemberRepository) DemoteCurrentLeader(ctx context.Context, lifeGroupID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.LifeGroupPersonMember{}).
		Where("life_group_id = ? AND position = ? AND is_active = ?", lifeGroupID, entity.PersonMemberPositionLeader, true).
		Update("position", entity.PersonMemberPositionCoLeader).Error
}
