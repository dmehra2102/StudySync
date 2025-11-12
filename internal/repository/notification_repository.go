package repository

import (
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uint) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.First(&notification, id).Error
	return &notification, err
}

func (r *notificationRepository) FindByUserID(userID uint, unreadOnly bool) ([]domain.Notification, error) {
	var notification []domain.Notification
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")

	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	err := query.Find(&notification).Error
	return notification, err
}

func (r *notificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&domain.Notification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"read":    true,
			"read_at": time.Now(),
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]any{
			"read":    true,
			"read_at": time.Now(),
		}).Error
}

func (r *notificationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Notification{}, id).Error
}

func (r *notificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}
