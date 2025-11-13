package service

import (
	"encoding/json"
	"fmt"

	"github.com/dmehra2102/StudySync/internal/domain"
	"github.com/dmehra2102/StudySync/pkg/redis"
)

type NotificationService struct {
	notificationRepo domain.NotificationRepository
	redis            *redis.Client
}

func NewNotificationService(notficationRepo domain.NotificationRepository, redis *redis.Client) *NotificationService {
	return &NotificationService{
		notificationRepo: notficationRepo,
		redis:            redis,
	}
}

type CreateNotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Data    any    `json:"data"`
}

func (s *NotificationService) CreateNotification(req CreateNotificationRequest) error {
	var dataStr string
	if req.Data != nil {
		dataBytes, err := json.Marshal(req.Data)
		if err == nil {
			dataStr = string(dataBytes)
		}
	}

	notification := &domain.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		Data:    dataStr,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return err
	}

	s.pushToRedis(notification)

	return nil
}

func (s *NotificationService) GetUserNotifications(userID uint, unreadOnly bool) ([]domain.Notification, error) {
	return s.notificationRepo.FindByUserID(userID, unreadOnly)
}

func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	notification, err := s.notificationRepo.FindByID(notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return domain.ErrUnauthorized
	}

	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *NotificationService) pushToRedis(notification *domain.Notification) {
	notificationData := map[string]any{
		"id":         notification.ID,
		"user_id":    notification.UserID,
		"title":      notification.Title,
		"message":    notification.Message,
		"type":       notification.Type,
		"created_at": notification.CreatedAt,
	}

	channel := s.getUserNotificationChannel(notification.UserID)
	s.redis.Publish(channel, notificationData)
}

func (s *NotificationService) getUserNotificationChannel(userID uint) string {
	return fmt.Sprintf("notifications:user:%d", userID)
}
