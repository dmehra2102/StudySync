package repository

import (
	"time"

	"github.com/dmehra2102/StudySync/internal/domain"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *taskRepository) FindByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("user_id = ?", userID).Order("due_date ASC").Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) FindOverdue(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("user_id = ? AND due_date < ? AND status != ?",
		userID, time.Now(), "completed").Order("due_date ASC").Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Task{}, id).Error
}

func (r *taskRepository) FindOverdueTasks() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("due_date < ? AND status = ?",
		time.Now(), "pending").Find(&tasks).Error
	return tasks, err
}
