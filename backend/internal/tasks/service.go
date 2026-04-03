package tasks

import (
	"context"
	"errors"
	"strings"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var ErrInvalidTaskInput = errors.New("invalid task input")
var ErrTaskNotFound = errors.New("task not found")

var allowedPriorities = map[string]struct{}{
	"HIGH":   {},
	"MEDIUM": {},
	"LOW":    {},
}

var allowedStatuses = map[string]struct{}{
	"PENDING": {},
	"DONE":    {},
}

type Service struct {
	repo *Repository
	now  func() time.Time
}

type CreateTaskInput struct {
	UserID      uint64
	Title       string
	Priority    string
	DueDate     *time.Time
	PlanDate    time.Time
	Status      string
	CompletedAt *time.Time
	Ext         []byte
}

type UpdateTaskInput struct {
	Title       *string
	Priority    *string
	DueDate     *time.Time
	PlanDate    *time.Time
	Status      *string
	CompletedAt *time.Time
	Ext         []byte
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) Create(input CreateTaskInput) (*models.Task, error) {
	title := strings.TrimSpace(input.Title)
	if input.UserID == 0 || title == "" || input.PlanDate.IsZero() {
		return nil, ErrInvalidTaskInput
	}

	priority, err := normalizePriority(input.Priority)
	if err != nil {
		return nil, err
	}
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if err := validateTaskConsistency(status, input.CompletedAt); err != nil {
		return nil, err
	}

	task := &models.Task{
		UserID:      input.UserID,
		Title:       title,
		Priority:    priority,
		PlanDate:    normalizeDate(input.PlanDate),
		Status:      status,
		CompletedAt: input.CompletedAt,
		Ext:         datatypes.JSON(input.Ext),
	}
	if input.DueDate != nil {
		dueDate := normalizeDate(*input.DueDate)
		task.DueDate = &dueDate
	}

	if err := s.repo.CreateTask(context.Background(), task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) ListByDate(userID uint64, date time.Time) ([]models.Task, error) {
	return s.repo.ListTasksByDate(context.Background(), userID, normalizeDate(date))
}

func (s *Service) Update(taskID, userID uint64, input UpdateTaskInput) (*models.Task, error) {
	if taskID == 0 || userID == 0 {
		return nil, ErrInvalidTaskInput
	}

	task, err := s.repo.GetTaskByID(context.Background(), taskID, userID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, ErrInvalidTaskInput
		}
		task.Title = title
	}
	if input.Priority != nil {
		priority, err := normalizePriority(*input.Priority)
		if err != nil {
			return nil, err
		}
		task.Priority = priority
	}
	if input.PlanDate != nil && !input.PlanDate.IsZero() {
		task.PlanDate = normalizeDate(*input.PlanDate)
	}
	if input.DueDate != nil {
		dueDate := normalizeDate(*input.DueDate)
		task.DueDate = &dueDate
	}
	if input.Status != nil {
		status, err := normalizeStatus(*input.Status)
		if err != nil {
			return nil, err
		}
		task.Status = status
	}
	if input.CompletedAt != nil {
		task.CompletedAt = input.CompletedAt
	}
	if input.Ext != nil {
		task.Ext = datatypes.JSON(input.Ext)
	}

	if task.Status == "PENDING" {
		if input.CompletedAt != nil {
			return nil, ErrInvalidTaskInput
		}
		task.CompletedAt = nil
	} else if task.CompletedAt == nil {
		return nil, ErrInvalidTaskInput
	}

	if err := s.repo.UpdateTask(context.Background(), task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) Delete(taskID, userID uint64) error {
	if taskID == 0 || userID == 0 {
		return ErrInvalidTaskInput
	}
	return s.repo.DeleteTask(context.Background(), taskID, userID)
}

func (s *Service) CarryOverPendingTasks(userID uint64, today time.Time) (int64, error) {
	if userID == 0 || today.IsZero() {
		return 0, ErrInvalidTaskInput
	}

	toDate := normalizeDate(today)
	fromDate := toDate.AddDate(0, 0, -1)
	return s.repo.CarryOverPending(context.Background(), userID, fromDate, toDate)
}

func normalizePriority(value string) (string, error) {
	priority := strings.ToUpper(strings.TrimSpace(value))
	if priority == "" {
		return "MEDIUM", nil
	}
	if _, ok := allowedPriorities[priority]; !ok {
		return "", ErrInvalidTaskInput
	}
	return priority, nil
}

func normalizeStatus(value string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" {
		return "PENDING", nil
	}
	if _, ok := allowedStatuses[status]; !ok {
		return "", ErrInvalidTaskInput
	}
	return status, nil
}

func validateTaskConsistency(status string, completedAt *time.Time) error {
	switch status {
	case "DONE":
		if completedAt == nil {
			return ErrInvalidTaskInput
		}
	case "PENDING":
		if completedAt != nil {
			return ErrInvalidTaskInput
		}
	}
	return nil
}

func isTaskNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrTaskNotFound)
}
