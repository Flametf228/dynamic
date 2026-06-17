package repository

import (
	"context"
	"time"
	"zzz/internal/models"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}
func (r *TaskRepo) StopActiveTasks(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("status = ?", "in_progress").
		Updates(map[string]interface{}{
			"status":   "stopped",
			"end_time": now,
		}).Error
}

func (r *TaskRepo) CreateTask(ctx context.Context) (string, error) {
	task := models.Task{
		Status: "in_progress",
	}
	err := r.db.WithContext(ctx).Create(&task).Error
	return task.ID, err
}

func (r *TaskRepo) CompleteTask(ctx context.Context, taskID string, status string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status":   status,
			"end_time": now,
		}).Error
}
