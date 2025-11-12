package repository

import (
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
	"gorm.io/gorm"
)

type studySessionRepository struct {
	db *gorm.DB
}

func NewStudySessionRepository(db *gorm.DB) domain.StudySessionRepository {
	return &studySessionRepository{db: db}
}

func (r *studySessionRepository) Create(session *domain.StudySession) error {
	return r.db.Create(session).Error
}

func (r *studySessionRepository) FindByID(id uint) (*domain.StudySession, error) {
	var session domain.StudySession
	err := r.db.Find(&session, id).Error
	return &session, err
}

func (r *studySessionRepository) FindByUserID(userID uint) ([]domain.StudySession, error) {
	var sessions []domain.StudySession
	err := r.db.Where("user_id = ?", userID).Order("planned_for ASC").Find(&sessions).Error
	return sessions, err
}

func (r *studySessionRepository) FindUpcoming(userID uint, from time.Time) ([]domain.StudySession, error) {
	var sessions []domain.StudySession
	err := r.db.Where("user_id = ? AND planned_for >= ? AND completed = ?",
		userID, from, false).Order("planned_for ASC").Find(&sessions).Error
	return sessions, err
}

func (r *studySessionRepository) Update(session *domain.StudySession) error {
	return r.db.Save(session).Error
}

func (r *studySessionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.StudySession{}, id).Error
}

func (r *studySessionRepository) GetUserStats(userID uint) (*domain.StudyStats, error) {
	var stats domain.StudyStats

	err := r.db.Model(&domain.StudySession{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_sessions, " +
			"SUM(CASE WHEN completed = true THEN 1 ELSE 0 END) as completed, " +
			"COALESCE(SUM(duration), 0) as total_time, " +
			"COALESCE(AVG(duration), 0) as average_duration",
		).Scan(&stats).Error

	return &stats, err
}
