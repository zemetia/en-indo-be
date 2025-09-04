package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type VisitorInformationRepository interface {
	Create(ctx context.Context, visitorInfo *entity.VisitorInformation) error
	GetAll(ctx context.Context) ([]entity.VisitorInformation, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.VisitorInformation, error)
	GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]entity.VisitorInformation, error)
	Update(ctx context.Context, visitorInfo *entity.VisitorInformation) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByVisitorID(ctx context.Context, visitorID uuid.UUID) error
}

type visitorInformationRepository struct {
	db *gorm.DB
}

func NewVisitorInformationRepository(db *gorm.DB) VisitorInformationRepository {
	return &visitorInformationRepository{
		db: db,
	}
}

func (r *visitorInformationRepository) Create(ctx context.Context, visitorInfo *entity.VisitorInformation) error {
	return r.db.WithContext(ctx).Create(visitorInfo).Error
}

func (r *visitorInformationRepository) GetAll(ctx context.Context) ([]entity.VisitorInformation, error) {
	var visitorInfos []entity.VisitorInformation
	err := r.db.WithContext(ctx).
		Preload("Visitor").
		Find(&visitorInfos).Error
	return visitorInfos, err
}

func (r *visitorInformationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.VisitorInformation, error) {
	var visitorInfo entity.VisitorInformation
	err := r.db.WithContext(ctx).
		Preload("Visitor").
		Where("id = ?", id).
		First(&visitorInfo).Error
	if err != nil {
		return nil, err
	}
	return &visitorInfo, nil
}

func (r *visitorInformationRepository) GetByVisitorID(ctx context.Context, visitorID uuid.UUID) ([]entity.VisitorInformation, error) {
	var visitorInfos []entity.VisitorInformation
	err := r.db.WithContext(ctx).
		Where("visitor_id = ?", visitorID).
		Find(&visitorInfos).Error
	return visitorInfos, err
}

func (r *visitorInformationRepository) Update(ctx context.Context, visitorInfo *entity.VisitorInformation) error {
	return r.db.WithContext(ctx).Save(visitorInfo).Error
}

func (r *visitorInformationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.VisitorInformation{}, id).Error
}

func (r *visitorInformationRepository) DeleteByVisitorID(ctx context.Context, visitorID uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Where("visitor_id = ?", visitorID).Delete(&entity.VisitorInformation{}).Error
}