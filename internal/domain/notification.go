package domain

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Read      bool      `json:"read" gorm:"default:false"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    time.Time `json:"read_at"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	FindByID(id uint) (*Notification, error)
	FindByUserID(userID uint, unreadOnly bool) ([]Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(userID uint) error
	Delete(id uint) error
	GetUnreadCount(userID uint) (int64, error)
}
