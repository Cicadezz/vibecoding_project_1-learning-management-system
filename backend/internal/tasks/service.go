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
	CarryCount  int
	Ext         []byte
}

type UpdateTaskInput struct {
	Title       *string
	Priority    *string
	DueDate     *time.Time
	PlanDate    *time.Time
	Status      *string
	CompletedAt *time.Time
	CarryCount  *int
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

	priority := strings.ToUpper(strings.TrimSpace(input.Priority))
	if priority == "" {
		priority = "MEDIUM"
	}

	status := strings.ToUpper(strings.TrimSpace(input.Status))
	if status == "" {
		status = "PENDING"
	}

	task := &models.Task{
		UserID:      input.UserID,
		Title:       title,
		Priority:    priority,
		PlanDate:    normalizeDate(input.PlanDate),
		Status:      status,
		CompletedAt: input.CompletedAt,
		CarryCount:  input.CarryCount,
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
		task.Priority = strings.ToUpper(strings.TrimSpace(*input.Priority))
	}
	if input.PlanDate != nil && !input.PlanDate.IsZero() {
		task.PlanDate = normalizeDate(*input.PlanDate)
	}
	if input.DueDate != nil {
		dueDate := normalizeDate(*input.DueDate)
		task.DueDate = &dueDate
	}
	if input.Status != nil {
		task.Status = strings.ToUpper(strings.TrimSpace(*input.Status))
	}
	if input.CompletedAt != nil {
		task.CompletedAt = input.CompletedAt
	}
	if input.CarryCount != nil {
		task.CarryCount = *input.CarryCount
	}
	if input.Ext != nil {
		task.Ext = datatypes.JSON(input.Ext)
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

func isTaskNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrTaskNotFound)
}
