package service

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/dto"
	"github.com/zemetia/en-indo-be/entity"
	"github.com/zemetia/en-indo-be/repository"
)

type NotificationService interface {
	Create(req *dto.NotificationRequest) (*dto.NotificationResponse, error)
	GetAll() ([]dto.NotificationResponse, error)
	GetByID(id uuid.UUID) (*dto.NotificationResponse, error)
	GetByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error)
	GetUnreadByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error)
	Update(id uuid.UUID, req *dto.NotificationRequest) (*dto.NotificationResponse, error)
	Delete(id uuid.UUID) error
	MarkAsRead(id uuid.UUID) error
	MarkAllAsRead(userID uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
}

type notificationService struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationService(notificationRepository repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *notificationService) Create(req *dto.NotificationRequest) (*dto.NotificationResponse, error) {
	notification := &entity.Notification{
		Title:    req.Title,
		Message:  req.Message,
		Type:     req.Type,
		UserID:   req.UserID,
		IsRead:   req.IsRead,
		ChurchID: &req.ChurchID,
	}

	if err := s.notificationRepository.Create(notification); err != nil {
		return nil, err
	}

	return s.GetByID(notification.ID)
}

func (s *notificationService) GetAll() ([]dto.NotificationResponse, error) {
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

func (s *notificationService) GetByID(id uuid.UUID) (*dto.NotificationResponse, error) {
	notification, err := s.notificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(notification), nil
}

func (s *notificationService) GetByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error) {
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

func (s *notificationService) GetUnreadByUserID(userID uuid.UUID) ([]dto.NotificationResponse, error) {
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

func (s *notificationService) Update(id uuid.UUID, req *dto.NotificationRequest) (*dto.NotificationResponse, error) {
	notification, err := s.notificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	notification.Title = req.Title
	notification.Message = req.Message
	notification.Type = req.Type
	notification.UserID = req.UserID
	notification.IsRead = req.IsRead
	notification.ChurchID = &req.ChurchID

	if err := s.notificationRepository.Update(notification); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *notificationService) Delete(id uuid.UUID) error {
	return s.notificationRepository.Delete(id)
}

func (s *notificationService) MarkAsRead(id uuid.UUID) error {
	return s.notificationRepository.MarkAsRead(id)
}

func (s *notificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepository.MarkAllAsRead(userID)
}

func (s *notificationService) DeleteByUserID(userID uuid.UUID) error {
	return s.notificationRepository.DeleteByUserID(userID)
}

func (s *notificationService) toResponse(notification *entity.Notification) *dto.NotificationResponse {
	var churchID uuid.UUID
	if notification.ChurchID != nil {
		churchID = *notification.ChurchID
	}

	return &dto.NotificationResponse{
		ID:        notification.ID,
		Title:     notification.Title,
		Message:   notification.Message,
		Type:      notification.Type,
		UserID:    notification.UserID,
		IsRead:    notification.IsRead,
		ChurchID:  churchID,
		CreatedAt: notification.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: notification.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
