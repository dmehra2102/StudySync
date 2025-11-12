package domain

import "time"

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" gorm:"default:medium"`
	Status      string     `json:"status" gorm:"default:pending"`
	DueDate     time.Time  `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskRepository interface {
	Create(task *Task) error
	FindByID(id uint) (*Task, error)
	FindByUserID(userID uint) ([]Task, error)
	FindOverdue(userID uint) ([]Task, error)
	Update(task *Task) error
	Delete(id uint) error
}
