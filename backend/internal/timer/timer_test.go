package timer

import (
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/database/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStopTimerCreatesSession(t *testing.T) {
	svc, db := buildTimerServiceForTest(t)
	userID := uint64(1)
	subjectID := uint64(2)
	startedAt := time.Date(2026, 4, 3, 10, 0, 0, 0, time.Local)
	stoppedAt := time.Date(2026, 4, 3, 11, 0, 0, 0, time.Local)
	note := "???"

	svc.now = func() time.Time { return startedAt }
	state, err := svc.Start(userID, subjectID)
	if err != nil {
		t.Fatalf("start timer: %v", err)
	}
	if state == nil {
		t.Fatal("expected timer state after start")
	}

	svc.now = func() time.Time { return stoppedAt }
	session, err := svc.Stop(userID, &note)
	if err != nil {
		t.Fatalf("stop timer: %v", err)
	}
	if session == nil {
		t.Fatal("expected study session after stop")
	}
	if session.DurationMinutes != 60 {
		t.Fatalf("expected duration minutes 60, got %d", session.DurationMinutes)
	}

	var persisted models.StudySession
	if err := db.Where("user_id = ? AND record_type = ?", userID, "TIMER").First(&persisted).Error; err != nil {
		t.Fatalf("load persisted study session: %v", err)
	}
	if persisted.DurationMinutes != 60 {
		t.Fatalf("expected persisted duration minutes 60, got %d", persisted.DurationMinutes)
	}
	if persisted.Note == nil || *persisted.Note != note {
		t.Fatalf("expected persisted note %q, got %#v", note, persisted.Note)
	}
}

func buildTimerServiceForTest(t *testing.T) (*Service, *gorm.DB) {
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

	return NewService(NewRepository(db)), db
}

func shouldSkipForCGO(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
