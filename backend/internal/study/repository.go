package study

import (
	"context"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSession(ctx context.Context, session *models.StudySession) error {
	return r.db.WithContext(ctx).Create(session).Error
}
