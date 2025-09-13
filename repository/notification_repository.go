package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	GetAll() ([]entity.Notification, error)
	GetByID(id uuid.UUID) (*entity.Notification, error)
	GetByUserID(userID uuid.UUID) ([]entity.Notification, error)
	GetUnreadByUserID(userID uuid.UUID) ([]entity.Notification, error)
	Update(notification *entity.Notification) error
	Delete(id uuid.UUID) error
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetAll() ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetByID(id uuid.UUID) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	return &notification, err
}

func (r *notificationRepository) GetByUserID(userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) GetUnreadByUserID(userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) Update(notification *entity.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Notification{}, "id = ?", id).Error
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *notificationRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.Notification{}).Error
}
