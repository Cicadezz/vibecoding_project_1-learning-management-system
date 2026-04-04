package checkin

import (
	"context"
	"errors"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

var ErrInvalidCheckinInput = errors.New("invalid checkin input")
var ErrStudySessionRequired = errors.New("study session required")

type Service struct {
	repo *Repository
	now  func() time.Time
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) CheckinToday(userID uint64, day time.Time) (*models.DailyCheckin, error) {
	if userID == 0 || day.IsZero() {
		return nil, ErrInvalidCheckinInput
	}

	hasStudy, err := s.repo.HasStudySessionOnDate(context.Background(), userID, day)
	if err != nil {
		return nil, err
	}
	if !hasStudy {
		return nil, ErrStudySessionRequired
	}

	checkin := &models.DailyCheckin{
		UserID:      userID,
		CheckinDate: normalizeDay(day),
		CheckedAt:   s.now(),
	}

	if err := s.repo.UpsertDailyCheckin(context.Background(), checkin); err != nil {
		return nil, err
	}

	persisted, err := s.repo.GetDailyCheckinByDate(context.Background(), userID, day)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return persisted, nil
}

func (s *Service) GetStreak(userID uint64) (int, error) {
	if userID == 0 {
		return 0, ErrInvalidCheckinInput
	}

	today := normalizeDay(s.now())
	hasToday, err := s.repo.HasCheckinOnDate(context.Background(), userID, today)
	if err != nil {
		return 0, err
	}
	if !hasToday {
		return 0, nil
	}

	streak := 1
	cursor := today.AddDate(0, 0, -1)
	for {
		hasCheckin, err := s.repo.HasCheckinOnDate(context.Background(), userID, cursor)
		if err != nil {
			return 0, err
		}
		if !hasCheckin {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}

	return streak, nil
}
