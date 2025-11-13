package service

import (
	"errors"
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
)

type StudySessionService struct {
	sessionRepo domain.StudySessionRepository
}

func NewStudySessionService(sessionRepo domain.StudySessionRepository) *StudySessionService {
	return &StudySessionService{
		sessionRepo: sessionRepo,
	}
}

type CreateSessionRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Subject     string    `json:"subject" binding:"required"`
	PlannedFor  time.Time `json:"planned_for" binding:"required"`
	Duration    int       `json:"duration" binding:"required,min=1"`
}

type UpdateSessionRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	PlannedFor  time.Time `json:"planned_for"`
	Duration    int       `json:"duration" binding:"min=1"`
	Completed   bool      `json:"completed"`
}

func (s *StudySessionService) CreateSession(userID uint, req CreateSessionRequest) (*domain.StudySession, error) {
	session := &domain.StudySession{
		UserId:      userID,
		Title:       req.Title,
		Description: req.Description,
		Subject:     req.Subject,
		PlannedFor:  req.PlannedFor,
		Duration:    req.Duration,
		Completed:   false,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *StudySessionService) GetUserSessions(userID uint) ([]domain.StudySession, error) {
	return s.sessionRepo.FindByUserID(userID)
}

func (s *StudySessionService) GetUpcomingSessions(userID uint) ([]domain.StudySession, error) {
	return s.sessionRepo.FindUpcoming(userID, time.Now())
}

func (s *StudySessionService) GetSession(userID, sessionID uint) (*domain.StudySession, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	if session.UserId != userID {
		return nil, errors.New("unauthorized")
	}

	return session, nil
}

func (s *StudySessionService) UpdateSession(userID, sessionID uint, req UpdateSessionRequest) (*domain.StudySession, error) {
	session, err := s.GetSession(userID, sessionID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		session.Title = req.Title
	}
	if req.Description != "" {
		session.Description = req.Description
	}
	if req.Subject != "" {
		session.Subject = req.Subject
	}
	if !req.PlannedFor.IsZero() {
		session.PlannedFor = req.PlannedFor
	}
	if req.Duration > 0 {
		session.Duration = req.Duration
	}

	if req.Completed && !session.Completed {
		session.Completed = true
		session.CompletedAt = time.Now()
	} else if !req.Completed && session.Completed {
		session.Completed = false
		session.CompletedAt = time.Time{}
	}

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *StudySessionService) DeleteSession(userID, sessionID uint) error {
	session, err := s.GetSession(userID, sessionID)
	if err != nil {
		return err
	}

	return s.sessionRepo.Delete(session.ID)
}

func (s *StudySessionService) GetUserStats(userID uint) (*domain.StudyStats, error) {
	return s.sessionRepo.GetUserStats(userID)
}
