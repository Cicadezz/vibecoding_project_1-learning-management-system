package tasks

import (
	"context"
	"testing"
	"time"

	"learning-growth-platform/internal/database/models"
)

func TestCarryOverPendingTasksIsIdempotentWithoutSQLite(t *testing.T) {
	repo := newFakeTaskRepository()
	svc := NewService(repo)
	userID := uint64(1)
	yesterday := time.Date(2026, 4, 1, 0, 0, 0, 0, time.Local)
	today := time.Date(2026, 4, 2, 0, 0, 0, 0, time.Local)

	created, err := svc.Create(CreateTaskInput{UserID: userID, Title: "刷题", PlanDate: yesterday})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if created == nil {
		t.Fatal("expected created task")
	}

	rows1, err := svc.CarryOverPendingTasks(userID, today)
	if err != nil {
		t.Fatalf("carry over first run: %v", err)
	}
	rows2, err := svc.CarryOverPendingTasks(userID, today)
	if err != nil {
		t.Fatalf("carry over second run: %v", err)
	}
	if rows1 != 1 || rows2 != 0 {
		t.Fatalf("expected carry-over rows 1 then 0, got %d then %d", rows1, rows2)
	}

	list, err := svc.ListByDate(userID, today)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	if len(list) != 1 || list[0].CarryCount != 1 {
		t.Fatalf("expected exactly one carried task with carry_count=1")
	}
}

func TestCarryOverPendingTasksKeepsListStableWithoutSQLite(t *testing.T) {
	repo := newFakeTaskRepository()
	svc := NewService(repo)
	userID := uint64(1)
	yesterday := time.Date(2026, 4, 1, 0, 0, 0, 0, time.Local)
	today := time.Date(2026, 4, 2, 0, 0, 0, 0, time.Local)

	_, err := svc.Create(CreateTaskInput{UserID: userID, Title: "错题整理", PlanDate: yesterday})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	_, err = svc.Create(CreateTaskInput{UserID: userID, Title: "已完成任务", PlanDate: yesterday, Status: "DONE", CompletedAt: ptrTime(time.Date(2026, 4, 1, 18, 0, 0, 0, time.Local))})
	if err != nil {
		t.Fatalf("create done task: %v", err)
	}

	rows, err := svc.CarryOverPendingTasks(userID, today)
	if err != nil {
		t.Fatalf("carry over: %v", err)
	}
	if rows != 1 {
		t.Fatalf("expected 1 carried row, got %d", rows)
	}

	list, err := svc.ListByDate(userID, today)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 task on today, got %d", len(list))
	}
	if list[0].Title != "错题整理" || list[0].CarryCount != 1 {
		t.Fatalf("expected carried pending task with carry_count=1, got %+v", list[0])
	}
}

type fakeTaskRepository struct {
	tasks []*models.Task
	next  uint64
}

func newFakeTaskRepository() *fakeTaskRepository {
	return &fakeTaskRepository{next: 1}
}

func (r *fakeTaskRepository) CreateTask(_ context.Context, task *models.Task) error {
	task.ID = r.next
	r.next++
	r.tasks = append(r.tasks, task)
	return nil
}

func (r *fakeTaskRepository) ListTasksByDate(_ context.Context, userID uint64, planDate time.Time) ([]models.Task, error) {
	out := make([]models.Task, 0)
	for _, task := range r.tasks {
		if task.UserID == userID && sameDay(task.PlanDate, planDate) {
			out = append(out, *task)
		}
	}
	return out, nil
}

func (r *fakeTaskRepository) GetTaskByID(_ context.Context, taskID, userID uint64) (*models.Task, error) {
	for _, task := range r.tasks {
		if task.ID == taskID && task.UserID == userID {
			return task, nil
		}
	}
	return nil, ErrTaskNotFound
}

func (r *fakeTaskRepository) UpdateTask(_ context.Context, task *models.Task) error {
	for i, existing := range r.tasks {
		if existing.ID == task.ID && existing.UserID == task.UserID {
			r.tasks[i] = task
			return nil
		}
	}
	return ErrTaskNotFound
}

func (r *fakeTaskRepository) DeleteTask(_ context.Context, taskID, userID uint64) error {
	for i, task := range r.tasks {
		if task.ID == taskID && task.UserID == userID {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return ErrTaskNotFound
}

func (r *fakeTaskRepository) CarryOverPending(_ context.Context, userID uint64, fromDate, toDate time.Time) (int64, error) {
	var rows int64
	for _, task := range r.tasks {
		if task.UserID != userID {
			continue
		}
		if !sameDay(task.PlanDate, fromDate) || task.Status != "PENDING" {
			continue
		}
		task.PlanDate = normalizeDate(toDate)
		task.CarryCount++
		rows++
	}
	return rows, nil
}

func sameDay(a, b time.Time) bool {
	a = normalizeDate(a)
	b = normalizeDate(b)
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
