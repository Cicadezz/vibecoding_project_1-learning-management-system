package stats

import (
	"context"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type StudySessionRow struct {
	SubjectID       uint64    `gorm:"column:subject_id"`
	SubjectName     string    `gorm:"column:subject_name"`
	StartAt         time.Time `gorm:"column:start_at"`
	EndAt           time.Time `gorm:"column:end_at"`
	DurationMinutes int       `gorm:"column:duration_minutes"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListStudySessionRowsBetween(ctx context.Context, userID uint64, startAt, endAt time.Time) ([]StudySessionRow, error) {
	var rows []StudySessionRow
	if err := r.db.WithContext(ctx).
		Table("study_sessions").
		Select("study_sessions.subject_id, COALESCE(subjects.name, '') AS subject_name, study_sessions.start_at, study_sessions.end_at, study_sessions.duration_minutes").
		Joins("LEFT JOIN subjects ON subjects.id = study_sessions.subject_id").
		Where("study_sessions.user_id = ? AND study_sessions.start_at < ? AND study_sessions.end_at > ?", userID, endAt, startAt).
		Order("study_sessions.start_at ASC, study_sessions.id ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repository) CountDoneTasks(ctx context.Context, userID uint64) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("user_id = ? AND status = ?", userID, "DONE").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) HasCheckinOnDate(ctx context.Context, userID uint64, day time.Time) (bool, error) {
	start := normalizeDay(day)
	end := start.AddDate(0, 0, 1)

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.DailyCheckin{}).
		Where("user_id = ? AND checkin_date >= ? AND checkin_date < ?", userID, start, end).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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

func mondayStart(day time.Time) time.Time {
	start := normalizeDay(day)
	if start.IsZero() {
		return time.Time{}
	}

	offset := (int(start.Weekday()) + 6) % 7
	return start.AddDate(0, 0, -offset)
}
