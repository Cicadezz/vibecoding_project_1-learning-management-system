package timer

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/database/models"
	"learning-growth-platform/internal/subjects"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStopTimerCreatesSessionAndResetsStateToIdle(t *testing.T) {
	svc, db, subjectRepo := buildTimerServiceForTest(t)
	userID := uint64(1)
	subjectID := createSubjectForUser(t, subjectRepo, userID, "math")
	startedAt := time.Date(2026, 4, 3, 10, 0, 0, 0, time.Local)
	stoppedAt := time.Date(2026, 4, 3, 11, 0, 0, 0, time.Local)
	note := "???"

	svc.now = func() time.Time { return startedAt }
	if _, err := svc.Start(userID, subjectID); err != nil {
		t.Fatalf("start timer: %v", err)
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

	var state models.TimerState
	if err := db.Where("user_id = ?", userID).First(&state).Error; err != nil {
		t.Fatalf("load timer state: %v", err)
	}
	if state.Status != "IDLE" {
		t.Fatalf("expected timer state IDLE after stop, got %q", state.Status)
	}
	if state.SubjectID != nil {
		t.Fatalf("expected subject cleared after stop, got %#v", state.SubjectID)
	}
	if state.StartedAt != nil {
		t.Fatalf("expected started_at cleared after stop, got %#v", state.StartedAt)
	}
	if state.LastResumedAt != nil {
		t.Fatalf("expected last_resumed_at cleared after stop, got %#v", state.LastResumedAt)
	}
	if state.PausedSeconds != 0 {
		t.Fatalf("expected paused_seconds reset to 0, got %d", state.PausedSeconds)
	}
}

func TestStartTimerRejectsAlreadyRunning(t *testing.T) {
	svc, _, subjectRepo := buildTimerServiceForTest(t)
	userID := uint64(1)
	subjectID := createSubjectForUser(t, subjectRepo, userID, "math")
	startAt := time.Date(2026, 4, 3, 9, 0, 0, 0, time.Local)

	svc.now = func() time.Time { return startAt }
	if _, err := svc.Start(userID, subjectID); err != nil {
		t.Fatalf("start timer: %v", err)
	}

	svc.now = func() time.Time { return startAt.Add(5 * time.Minute) }
	if _, err := svc.Start(userID, subjectID); !errors.Is(err, ErrTimerAlreadyRunning) {
		t.Fatalf("expected ErrTimerAlreadyRunning, got %v", err)
	}
}

func TestStopTimerRejectsWhenNotRunning(t *testing.T) {
	svc, _, _ := buildTimerServiceForTest(t)

	if _, err := svc.Stop(1, nil); !errors.Is(err, ErrTimerNotRunning) {
		t.Fatalf("expected ErrTimerNotRunning, got %v", err)
	}
}

func TestStartTimerRejectsMissingOrForeignSubject(t *testing.T) {
	svc, _, subjectRepo := buildTimerServiceForTest(t)
	userID := uint64(1)
	foreignSubjectID := createSubjectForUser(t, subjectRepo, 2, "history")
	startAt := time.Date(2026, 4, 3, 8, 0, 0, 0, time.Local)
	svc.now = func() time.Time { return startAt }

	t.Run("missing subject", func(t *testing.T) {
		if _, err := svc.Start(userID, 9999); !errors.Is(err, ErrSubjectNotFound) {
			t.Fatalf("expected ErrSubjectNotFound, got %v", err)
		}
	})

	t.Run("foreign subject", func(t *testing.T) {
		if _, err := svc.Start(userID, foreignSubjectID); !errors.Is(err, ErrSubjectNotFound) {
			t.Fatalf("expected ErrSubjectNotFound, got %v", err)
		}
	})
}

func buildTimerServiceForTest(t *testing.T) (*Service, *gorm.DB, *subjects.Repository) {
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

	subjectRepo := subjects.NewRepository(db)
	return NewService(NewRepository(db), subjectRepo), db, subjectRepo
}

func createSubjectForUser(t *testing.T, repo *subjects.Repository, userID uint64, name string) uint64 {
	t.Helper()

	subject := &models.Subject{
		UserID: userID,
		Name:   name,
	}
	if err := repo.CreateSubject(context.Background(), subject); err != nil {
		t.Fatalf("create subject: %v", err)
	}
	return subject.ID
}

func shouldSkipForCGO(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
