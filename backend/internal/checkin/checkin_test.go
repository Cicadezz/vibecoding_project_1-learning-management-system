package checkin

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

func TestCheckinRequiresStudyAndIsIdempotent(t *testing.T) {
	svc, db, subjectRepo := buildCheckinServiceForTest(t)
	userID := uint64(1)
	subjectID := createCheckinSubjectForUser(t, subjectRepo, userID, "math")
	day := time.Date(2026, 4, 3, 0, 0, 0, 0, time.Local)

	if _, err := svc.CheckinToday(userID, day); !errors.Is(err, ErrStudySessionRequired) {
		t.Fatalf("expected ErrStudySessionRequired, got %v", err)
	}

	if err := createCheckinStudySession(t, db, userID, subjectID, day); err != nil {
		t.Fatalf("create study session: %v", err)
	}

	first, err := svc.CheckinToday(userID, day)
	if err != nil {
		t.Fatalf("checkin after study session: %v", err)
	}
	if first == nil {
		t.Fatal("expected checkin record")
	}

	second, err := svc.CheckinToday(userID, day)
	if err != nil {
		t.Fatalf("repeat checkin: %v", err)
	}
	if second == nil {
		t.Fatal("expected repeated checkin record")
	}

	var count int64
	if err := db.Model(&models.DailyCheckin{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		t.Fatalf("count checkins: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one checkin row, got %d", count)
	}

	streak, err := svc.GetStreak(userID)
	if err != nil {
		t.Fatalf("get streak: %v", err)
	}
	if streak != 1 {
		t.Fatalf("expected streak 1, got %d", streak)
	}
}

func TestCheckinStreakAnchorsOnToday(t *testing.T) {
	t.Run("continuous streak includes today", func(t *testing.T) {
		svc, db, subjectRepo := buildCheckinServiceForTest(t)
		today := time.Date(2026, 4, 5, 0, 0, 0, 0, time.Local)
		svc.now = func() time.Time { return today.Add(12 * time.Hour) }
		userID := uint64(1)
		subjectID := createCheckinSubjectForUser(t, subjectRepo, userID, "math")

		for _, day := range []time.Time{today.AddDate(0, 0, -2), today.AddDate(0, 0, -1), today} {
			if err := createCheckinStudySession(t, db, userID, subjectID, day); err != nil {
				t.Fatalf("create study session: %v", err)
			}
			if _, err := svc.CheckinToday(userID, day); err != nil {
				t.Fatalf("checkin %v: %v", day, err)
			}
		}

		streak, err := svc.GetStreak(userID)
		if err != nil {
			t.Fatalf("get streak: %v", err)
		}
		if streak != 3 {
			t.Fatalf("expected streak 3, got %d", streak)
		}
	})

	t.Run("gap on today returns zero", func(t *testing.T) {
		svc, db, subjectRepo := buildCheckinServiceForTest(t)
		today := time.Date(2026, 4, 5, 0, 0, 0, 0, time.Local)
		svc.now = func() time.Time { return today.Add(12 * time.Hour) }
		userID := uint64(1)
		subjectID := createCheckinSubjectForUser(t, subjectRepo, userID, "math")

		for _, day := range []time.Time{today.AddDate(0, 0, -2), today.AddDate(0, 0, -1)} {
			if err := createCheckinStudySession(t, db, userID, subjectID, day); err != nil {
				t.Fatalf("create study session: %v", err)
			}
			if _, err := svc.CheckinToday(userID, day); err != nil {
				t.Fatalf("checkin %v: %v", day, err)
			}
		}

		streak, err := svc.GetStreak(userID)
		if err != nil {
			t.Fatalf("get streak: %v", err)
		}
		if streak != 0 {
			t.Fatalf("expected streak 0, got %d", streak)
		}
	})
}

func buildCheckinServiceForTest(t *testing.T) (*Service, *gorm.DB, *subjects.Repository) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGOCheckin(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGOCheckin(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	subjectRepo := subjects.NewRepository(db)
	svc := NewService(NewRepository(db))
	svc.now = func() time.Time { return time.Date(2026, 4, 3, 12, 0, 0, 0, time.Local) }
	return svc, db, subjectRepo
}

func createCheckinSubjectForUser(t *testing.T, repo *subjects.Repository, userID uint64, name string) uint64 {
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

func createCheckinStudySession(t *testing.T, db *gorm.DB, userID, subjectID uint64, day time.Time) error {
	t.Helper()

	return db.Create(&models.StudySession{
		UserID:          userID,
		SubjectID:       subjectID,
		RecordType:      "MANUAL",
		StartAt:         day.Add(9 * time.Hour),
		EndAt:           day.Add(10 * time.Hour),
		DurationMinutes: 60,
	}).Error
}

func shouldSkipForCGOCheckin(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
