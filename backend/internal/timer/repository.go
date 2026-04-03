package timer

import (
	"context"
	"errors"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

var ErrTimerStateNotFound = errors.New("timer state not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStateByUser(ctx context.Context, userID uint64) (*models.TimerState, error) {
	var state models.TimerState
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimerStateNotFound
		}
		return nil, err
	}
	return &state, nil
}

func (r *Repository) SaveState(ctx context.Context, state *models.TimerState) error {
	return r.db.WithContext(ctx).Save(state).Error
}
