package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type VisitorRepository interface {
	Create(ctx context.Context, visitor *entity.Visitor) error
	GetAll(ctx context.Context) ([]entity.Visitor, error)
	Search(ctx context.Context, search *dto.VisitorSearchDto) ([]entity.Visitor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Visitor, error)
	Update(ctx context.Context, visitor *entity.Visitor) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithInformation(ctx context.Context, id uuid.UUID) (*entity.Visitor, error)
}

type visitorRepository struct {
	db *gorm.DB
}

func NewVisitorRepository(db *gorm.DB) VisitorRepository {
	return &visitorRepository{
		db: db,
	}
}

func (r *visitorRepository) Create(ctx context.Context, visitor *entity.Visitor) error {
	return r.db.WithContext(ctx).Create(visitor).Error
}

func (r *visitorRepository) GetAll(ctx context.Context) ([]entity.Visitor, error) {
	var visitors []entity.Visitor
	err := r.db.WithContext(ctx).
		Preload("Information").
		Preload("Kabupaten").
		Preload("Kabupaten.Provinsi").
		Find(&visitors).Error
	return visitors, err
}

func (r *visitorRepository) Search(ctx context.Context, search *dto.VisitorSearchDto) ([]entity.Visitor, error) {
	var visitors []entity.Visitor
	query := r.db.WithContext(ctx).
		Preload("Information").
		Preload("Kabupaten").
		Preload("Kabupaten.Provinsi")

	if search.Name != nil {
		query = query.Where("name LIKE ?", "%"+*search.Name+"%")
	}

	if search.IGUsername != nil {
		query = query.Where("ig_username LIKE ?", "%"+*search.IGUsername+"%")
	}

	if search.PhoneNumber != nil {
		query = query.Where("phone_number LIKE ?", "%"+*search.PhoneNumber+"%")
	}

	if search.KabupatenID != nil {
		query = query.Where("kabupaten_id = ?", *search.KabupatenID)
	}

	err := query.Find(&visitors).Error
	return visitors, err
}

func (r *visitorRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Visitor, error) {
	var visitor entity.Visitor
	err := r.db.WithContext(ctx).
		Preload("Kabupaten").
		Preload("Kabupaten.Provinsi").
		Where("id = ?", id).
		First(&visitor).Error
	if err != nil {
		return nil, err
	}
	return &visitor, nil
}

func (r *visitorRepository) Update(ctx context.Context, visitor *entity.Visitor) error {
	return r.db.WithContext(ctx).Save(visitor).Error
}

func (r *visitorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Visitor{}, id).Error
}

func (r *visitorRepository) GetWithInformation(ctx context.Context, id uuid.UUID) (*entity.Visitor, error) {
	var visitor entity.Visitor
	err := r.db.WithContext(ctx).
		Preload("Information").
		Preload("Kabupaten").
		Preload("Kabupaten.Provinsi").
		Where("id = ?", id).
		First(&visitor).Error
	if err != nil {
		return nil, err
	}
	return &visitor, nil
}
