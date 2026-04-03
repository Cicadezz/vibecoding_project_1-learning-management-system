package tasks

import (
	"context"
	"errors"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *Repository) ListTasksByDate(ctx context.Context, userID uint64, planDate time.Time) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND plan_date = ?", userID, normalizeDate(planDate)).
		Order("id ASC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID, userID uint64) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *Repository) DeleteTask(ctx context.Context, taskID, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		Delete(&models.Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CarryOverPending(ctx context.Context, userID uint64, fromDate, toDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("user_id = ? AND plan_date = ? AND status = ?", userID, normalizeDate(fromDate), "PENDING").
		Updates(map[string]any{
			"plan_date":   normalizeDate(toDate),
			"carry_count": gorm.Expr("carry_count + 1"),
		})
	return result.RowsAffected, result.Error
}

func normalizeDate(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	loc := t.Location()
	if loc == nil {
		loc = time.Local
	}

	year, month, day := t.In(loc).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}
