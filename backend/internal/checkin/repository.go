package checkin

import (
	"context"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func normalizeDay(day time.Time) time.Time {
	if day.IsZero() {
		return time.Time{}
	}

	loc := day.Location()
	if loc == nil {
		loc = time.Local
	}

	d := day.In(loc)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
}

func (r *Repository) HasStudySessionOnDate(ctx context.Context, userID uint64, day time.Time) (bool, error) {
	start := normalizeDay(day)
	end := start.AddDate(0, 0, 1)

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.StudySession{}).
		Where("user_id = ? AND start_at >= ? AND start_at < ?", userID, start, end).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UpsertDailyCheckin(ctx context.Context, checkin *models.DailyCheckin) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "checkin_date"}},
			DoNothing: true,
		}).
		Create(checkin).Error
}

func (r *Repository) GetDailyCheckinByDate(ctx context.Context, userID uint64, day time.Time) (*models.DailyCheckin, error) {
	start := normalizeDay(day)
	end := start.AddDate(0, 0, 1)

	var checkin models.DailyCheckin
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userID, start, end).
		First(&checkin).Error; err != nil {
		return nil, err
	}
	return &checkin, nil
}

func (r *Repository) ListCheckinDates(ctx context.Context, userID uint64, until time.Time) ([]time.Time, error) {
	end := normalizeDay(until).AddDate(0, 0, 1)

	var dates []time.Time
	if err := r.db.WithContext(ctx).
		Model(&models.DailyCheckin{}).
		Where("user_id = ? AND checkin_date < ?", userID, end).
		Order("checkin_date DESC").
		Pluck("checkin_date", &dates).Error; err != nil {
		return nil, err
	}
	return dates, nil
}
