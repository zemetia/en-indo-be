package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPersonID(ctx context.Context, personID uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateActivationStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	HasActivePelayanan(ctx context.Context, personID uuid.UUID) (bool, error)
	RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error)
	GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error)
	CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	if err := r.db.WithContext(ctx).Preload("Person").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Person").First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Person").First(&user, "email = ?", email).Error
	return &user, err
}

func (r *userRepository) GetByPersonID(ctx context.Context, personID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Person").First(&user, "person_id = ?", personID).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *userRepository) UpdateActivationStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *userRepository) HasActivePelayanan(ctx context.Context, personID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.PersonPelayananGereja{}).Where("person_id = ?", personID).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User
	var err error
	var count int64

	req.Default()

	query := tx.WithContext(ctx).Model(&entity.User{})
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(req)).Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := TotalPage(count, int64(req.PerPage))
	return dto.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Preload("Person").Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}
