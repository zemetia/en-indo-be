package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type LaguRepository interface {
	Create(ctx context.Context, lagu *entity.Lagu) error
	FindAll(ctx context.Context) ([]entity.Lagu, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lagu, error)
	Update(ctx context.Context, lagu *entity.Lagu) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type laguRepository struct {
	db *gorm.DB
}

func NewLaguRepository(db *gorm.DB) LaguRepository {
	return &laguRepository{
		db: db,
	}
}

func (r *laguRepository) Create(ctx context.Context, lagu *entity.Lagu) error {
	return r.db.WithContext(ctx).Create(lagu).Error
}

func (r *laguRepository) FindAll(ctx context.Context) ([]entity.Lagu, error) {
	var lagus []entity.Lagu
	err := r.db.WithContext(ctx).
		Where("original_lagu_id IS NULL").
		Preload("Versions").
		Find(&lagus).Error
	return lagus, err
}

func (r *laguRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Lagu, error) {
	var lagu entity.Lagu
	err := r.db.WithContext(ctx).
		Preload("Versions").
		Preload("OriginalLagu").
		First(&lagu, "id = ?", id).Error
	return &lagu, err
}

func (r *laguRepository) Update(ctx context.Context, lagu *entity.Lagu) error {
	return r.db.WithContext(ctx).Save(lagu).Error
}

func (r *laguRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Lagu{}, "id = ?", id).Error
}
