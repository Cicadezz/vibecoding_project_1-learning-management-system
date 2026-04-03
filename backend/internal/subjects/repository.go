package subjects

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

func (r *Repository) CreateSubject(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *Repository) ListSubjects(ctx context.Context, userID uint64) ([]models.Subject, error) {
	var subjects []models.Subject
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (r *Repository) GetSubjectByID(ctx context.Context, subjectID, userID uint64) (*models.Subject, error) {
	var subject models.Subject
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", subjectID, userID).
		First(&subject).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *Repository) UpdateSubject(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).Save(subject).Error
}

func (r *Repository) DeleteSubject(ctx context.Context, subjectID, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", subjectID, userID).
		Delete(&models.Subject{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
