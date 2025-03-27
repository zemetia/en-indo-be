package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetByID(id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error
	return &role, err
}

func (r *RoleRepository) Update(role *entity.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Role{}, "id = ?", id).Error
}

func (r *RoleRepository) AddPermissions(id uuid.UUID, permissionIDs []uuid.UUID) error {
	var permissions []entity.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.Role{}).Where("id = ?", id).Association("Permissions").Append(permissions)
}

func (r *RoleRepository) RemovePermissions(id uuid.UUID, permissionIDs []uuid.UUID) error {
	var permissions []entity.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.Role{}).Where("id = ?", id).Association("Permissions").Delete(permissions)
}

func (r *RoleRepository) AssignToUser(userID uuid.UUID, roleIDs []uuid.UUID) error {
	// Hapus semua relasi role yang ada
	if err := r.db.Model(&entity.User{}).Where("id = ?", userID).Association("Roles").Clear(); err != nil {
		return err
	}

	// Tambahkan role baru
	var roles []entity.Role
	if err := r.db.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.User{}).Where("id = ?", userID).Association("Roles").Replace(roles)
}

func (r *RoleRepository) RemoveFromUser(userID uuid.UUID, roleIDs []uuid.UUID) error {
	var roles []entity.Role
	if err := r.db.Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&entity.User{}).Where("id = ?", userID).Association("Roles").Delete(roles)
}
