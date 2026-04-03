package study

import (
	"context"
	"errors"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/datatypes"
)

var ErrInvalidStudyInput = errors.New("invalid study input")

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.StudySession) error
}

type Service struct {
	repo SessionRepository
}

type CreateManualSessionInput struct {
	UserID    uint64
	SubjectID uint64
	StartAt   time.Time
	EndAt     time.Time
	Note      *string
	Ext       []byte
}

func NewService(repo SessionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateManual(input CreateManualSessionInput) (*models.StudySession, error) {
	if input.UserID == 0 || input.SubjectID == 0 || input.StartAt.IsZero() || input.EndAt.IsZero() || input.EndAt.Before(input.StartAt) {
		return nil, ErrInvalidStudyInput
	}

	durationMinutes := int(input.EndAt.Sub(input.StartAt).Seconds() / 60)
	session := &models.StudySession{
		UserID:          input.UserID,
		SubjectID:       input.SubjectID,
		RecordType:      "MANUAL",
		StartAt:         input.StartAt,
		EndAt:           input.EndAt,
		DurationMinutes: durationMinutes,
		Note:            input.Note,
		Ext:             datatypes.JSON(input.Ext),
	}

	if err := s.repo.CreateSession(context.Background(), session); err != nil {
		return nil, err
	}
	return session, nil
}
