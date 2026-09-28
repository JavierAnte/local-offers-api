package services

import (
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo notificationRepository
}

type notificationRepository interface {
	FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error)
}

func NewNotificationService(repo notificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error) {
	return s.repo.FindByRecipientID(recipientID)
}
