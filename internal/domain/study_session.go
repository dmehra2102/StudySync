package domain

import "time"

type StudySession struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserId      uint      `json:"user_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Subject     string    `json:"subject" gorm:"not null"`
	PlannedFor  time.Time `json:"panned_for" gorm:"not null"`
	Duration    int       `json:"duration" gorm:"not null"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StudySessionRepository interface {
	Create(session *StudySession) error
	FindByID(id uint) (*StudySession, error)
	FindByUserID(userID uint) ([]StudySession, error)
	FindUpcoming(userID uint, from time.Time) ([]StudySession, error)
	Update(session *StudySession) error
	Delete(id uint) error
	GetUserStats(userID uint) (*StudyStats, error)
	FindUpcomingSessions(upcomingTime time.Time) ([]StudySession, error)
	FindCompletedSessions(userID uint, from, to time.Time) ([]StudySession, error)
}

type StudyStats struct {
	TotalSessions   int     `json:"total_sessions"`
	Completed       int     `json:"completed_Sessions"`
	TotalTime       int     `json:"total_time_minutes"`
	AverageDuration float64 `json:"average_duration_minutes"`
}
