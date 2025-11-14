package service

import (
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
)

type TaskService struct {
	taskRepo domain.TaskRepository
}

func NewTaskService(taskRepo domain.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" binding:"oneof=low medium high"`
	DueDate     time.Time `json:"due_date" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" binding:"oneof=low medium high"`
	Status      string    `json:"status" binding:"oneof=pending in_progress completed"`
	DueDate     time.Time `json:"due_date"`
}

func (s *TaskService) CreateTask(userID uint, req CreateTaskRequest) (*domain.Task, error) {
	task := &domain.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      "pending",
		DueDate:     req.DueDate,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]domain.Task, error) {
	return s.taskRepo.FindByUserID(userID)
}

func (s *TaskService) GetOverdueTasks(userID uint) ([]domain.Task, error) {
	return s.taskRepo.FindOverdue(userID)
}

func (s *TaskService) GetTask(userID, taskID uint) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if task.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return task, nil
}

func (s *TaskService) UpdateTask(userID, taskID uint, req UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.GetTask(userID, taskID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.Status != "" {
		task.Status = req.Status
		if req.Status == "completed" && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		} else if req.Status != "completed" {
			task.CompletedAt = nil
		}
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(userID, taskID uint) error {
	task, err := s.GetTask(userID, taskID)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(task.ID)
}
