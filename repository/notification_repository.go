package repository

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetAll() ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetByID(id uuid.UUID) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	return &notification, err
}

func (r *NotificationRepository) GetByUserID(userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) GetUnreadByUserID(userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) Update(notification *entity.Notification) error {
	return r.db.Save(notification).Error
}

func (r *NotificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Notification{}, "id = ?", id).Error
}

func (r *NotificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *NotificationRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.Notification{}).Error
}
