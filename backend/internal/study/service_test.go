package study

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

func TestCreateManualRejectsMissingOrForeignSubject(t *testing.T) {
	svc, subjectRepo := buildStudyServiceForTest(t)
	userID := uint64(1)
	foreignSubjectID := createStudySubjectForUser(t, subjectRepo, 2, "history")
	startAt := time.Date(2026, 4, 3, 8, 0, 0, 0, time.Local)
	endAt := startAt.Add(45 * time.Minute)

	t.Run("missing subject", func(t *testing.T) {
		_, err := svc.CreateManual(CreateManualSessionInput{
			UserID:    userID,
			SubjectID: 9999,
			StartAt:   startAt,
			EndAt:     endAt,
		})
		if !errors.Is(err, ErrSubjectNotFound) {
			t.Fatalf("expected ErrSubjectNotFound, got %v", err)
		}
	})

	t.Run("foreign subject", func(t *testing.T) {
		_, err := svc.CreateManual(CreateManualSessionInput{
			UserID:    userID,
			SubjectID: foreignSubjectID,
			StartAt:   startAt,
			EndAt:     endAt,
		})
		if !errors.Is(err, ErrSubjectNotFound) {
			t.Fatalf("expected ErrSubjectNotFound, got %v", err)
		}
	})
}

func buildStudyServiceForTest(t *testing.T) (*Service, *subjects.Repository) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGOStudy(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGOStudy(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	subjectRepo := subjects.NewRepository(db)
	return NewService(NewRepository(db), subjectRepo), subjectRepo
}

func createStudySubjectForUser(t *testing.T, repo *subjects.Repository, userID uint64, name string) uint64 {
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

func shouldSkipForCGOStudy(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
