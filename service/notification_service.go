package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type NotificationService struct {
	notificationRepository *repository.NotificationRepository
}

func NewNotificationService(notificationRepository *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) Create(req *dto.NotificationRequest) (*dto.NotificationResponse, error) {
	notification := &entity.Notification{
		Title:         req.Title,
		Message:       req.Message,
		Type:          req.Type,
		UserID:        req.UserID,
		IsRead:        req.IsRead,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
	}

	if err := s.notificationRepository.Create(notification); err != nil {
		return nil, err
	}

	return s.GetByID(notification.ID)
}

func (s *NotificationService) GetAll() ([]dto.NotificationResponse, error) {
	notifications, err := s.notificationRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *s.toResponse(&notification))
	}

	return responses, nil
}

func (s *NotificationService) GetByID(id uuid.UUID) (*dto.NotificationResponse, error) {
	notification, err := s.notificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(notification), nil
}

func (s *NotificationService) GetByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error) {
	notifications, err := s.notificationRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *s.toResponse(&notification))
	}

	return responses, nil
}

func (s *NotificationService) GetUnreadByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error) {
	notifications, err := s.notificationRepository.GetUnreadByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *s.toResponse(&notification))
	}

	return responses, nil
}

func (s *NotificationService) Update(id uuid.UUID, req *dto.NotificationRequest) (*dto.NotificationResponse, error) {
	notification, err := s.notificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	notification.Title = req.Title
	notification.Message = req.Message
	notification.Type = req.Type
	notification.UserID = req.UserID
	notification.IsRead = req.IsRead
	notification.ReferenceID = req.ReferenceID
	notification.ReferenceType = req.ReferenceType

	if err := s.notificationRepository.Update(notification); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *NotificationService) Delete(id uuid.UUID) error {
	return s.notificationRepository.Delete(id)
}

func (s *NotificationService) MarkAsRead(id uuid.UUID) error {
	return s.notificationRepository.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepository.MarkAllAsRead(userID)
}

func (s *NotificationService) DeleteByUserID(userID uuid.UUID) error {
	return s.notificationRepository.DeleteByUserID(userID)
}

func (s *NotificationService) toResponse(notification *entity.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:            notification.ID,
		Title:         notification.Title,
		Message:       notification.Message,
		Type:          notification.Type,
		UserID:        notification.UserID,
		IsRead:        notification.IsRead,
		ReferenceID:   notification.ReferenceID,
		ReferenceType: notification.ReferenceType,
		CreatedAt:     notification.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     notification.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
