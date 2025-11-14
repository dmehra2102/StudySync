package service

import (
	"fmt"
	"log"
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	cron             *cron.Cron
	studySessionRepo domain.StudySessionRepository
	taskRepo         domain.TaskRepository
	notificationRepo domain.NotificationRepository
	notificationSvc  *NotificationService
}

func NewSchedulerService(
	studySessionRepo domain.StudySessionRepository,
	taskRepo domain.TaskRepository,
	notificationRepo domain.NotificationRepository,
	notificationSvc *NotificationService,
) *SchedulerService {
	return &SchedulerService{
		cron:             cron.New(),
		studySessionRepo: studySessionRepo,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
		notificationSvc:  notificationSvc,
	}
}

func (s *SchedulerService) Start() {
	// Check for upcoming study sessions every 5 minutes
	s.cron.AddFunc("*/5 * * * *", s.checkUpcomingSessions)

	// Check for overdue tasks every hour
	s.cron.AddFunc("0 * * * *", s.checkOverdueTasks)

	s.cron.Start()
	log.Println("Background scheduler started")
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
	log.Println("Background scheduler stopped")
}

func (s *SchedulerService) checkUpcomingSessions() {
	now := time.Now()
	upcomingTime := now.Add(30 * time.Minute)

	sessions, err := s.studySessionRepo.FindUpcomingSessions(upcomingTime)
	if err != nil {
		log.Printf("Error fetching upcoming sessions: %v", err)
		return
	}

	for _, session := range sessions {
		notification := CreateNotificationRequest{
			UserID:  session.UserId,
			Title:   "Upcoming Study Session",
			Message: fmt.Sprintf("You have a study session '%s' starting in 30 minutes", session.Title),
			Type:    "reminder",
			Data: map[string]any{
				"session_id": session.ID,
				"start_time": session.PlannedFor,
				"duration":   session.Duration,
			},
		}

		if err := s.notificationSvc.CreateNotification(notification); err != nil {
			log.Printf("Error creating notification for session %d: %v", session.ID, err)
		}
	}
}

func (s *SchedulerService) checkOverdueTasks() {
	tasks, err := s.taskRepo.FindOverdueTasks()
	if err != nil {
		log.Printf("Error fetching overdue tasks: %v", err)
		return
	}

	for _, task := range tasks {
		notification := CreateNotificationRequest{
			UserID:  task.UserID,
			Title:   "Overdue Task",
			Message: fmt.Sprintf("Task '%s' is overdue", task.Title),
			Type:    "alert",
			Data: map[string]any{
				"task_id":  task.ID,
				"due_date": task.DueDate,
			},
		}

		if err := s.notificationSvc.CreateNotification(notification); err != nil {
			log.Printf("Error creating notification for task %d: %v", task.ID, err)
		}
	}
}
