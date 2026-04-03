package tasks

import (
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCarryOverPendingTasksIsIdempotent(t *testing.T) {
	svc := buildTaskServiceForTest(t)
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

func buildTaskServiceForTest(t *testing.T) *Service {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	return NewService(NewRepository(db))
}

func shouldSkipForCGO(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
