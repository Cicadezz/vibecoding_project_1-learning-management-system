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

	_, _ = svc.Create(CreateTaskInput{UserID: userID, Title: "刷题", PlanDate: yesterday})
	_, _ = svc.CarryOverPendingTasks(userID, today)
	_, _ = svc.CarryOverPendingTasks(userID, today)

	list, _ := svc.ListByDate(userID, today)
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
