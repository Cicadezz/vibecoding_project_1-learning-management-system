package stats

import (
	"context"
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/database/models"
	"learning-growth-platform/internal/subjects"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestOverviewUsesMondayWeekStart(t *testing.T) {
	svc, subjectRepo := buildStatsServiceForTest(t)
	userID := uint64(1)
	subjectID := createStatsSubjectForUser(t, subjectRepo, userID, "math")
	ref := time.Date(2026, 4, 2, 12, 0, 0, 0, time.Local)
	svc.now = func() time.Time { return ref }

	if err := createStatsStudySession(t, svc.repo.db, userID, subjectID, time.Date(2026, 3, 30, 9, 0, 0, 0, time.Local), 60); err != nil {
		t.Fatalf("create monday session: %v", err)
	}
	if err := createStatsStudySession(t, svc.repo.db, userID, subjectID, time.Date(2026, 3, 29, 9, 0, 0, 0, time.Local), 90); err != nil {
		t.Fatalf("create sunday session: %v", err)
	}

	overview, err := svc.Overview(userID)
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if overview.WeekMinutes != 60 {
		t.Fatalf("expected week minutes 60, got %d", overview.WeekMinutes)
	}
}

func buildStatsServiceForTest(t *testing.T) (*Service, *subjects.Repository) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGOStats(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGOStats(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	subjectRepo := subjects.NewRepository(db)
	return NewService(NewRepository(db)), subjectRepo
}

func createStatsSubjectForUser(t *testing.T, repo *subjects.Repository, userID uint64, name string) uint64 {
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

func createStatsStudySession(t *testing.T, db *gorm.DB, userID, subjectID uint64, startAt time.Time, durationMinutes int) error {
	t.Helper()

	return db.Create(&models.StudySession{
		UserID:          userID,
		SubjectID:       subjectID,
		RecordType:      "MANUAL",
		StartAt:         startAt,
		EndAt:           startAt.Add(time.Duration(durationMinutes) * time.Minute),
		DurationMinutes: durationMinutes,
	}).Error
}

func shouldSkipForCGOStats(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
