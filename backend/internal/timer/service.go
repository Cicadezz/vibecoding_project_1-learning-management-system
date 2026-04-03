package timer

import (
	"context"
	"errors"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidTimerInput = errors.New("invalid timer input")
var ErrTimerNotRunning = errors.New("timer not running")
var ErrTimerAlreadyRunning = errors.New("timer already running")
var ErrSubjectNotFound = errors.New("subject not found")

type SubjectRepository interface {
	GetSubjectByID(ctx context.Context, subjectID, userID uint64) (*models.Subject, error)
}

type Service struct {
	repo        *Repository
	subjectRepo SubjectRepository
	now         func() time.Time
}

func NewService(repo *Repository, subjectRepo SubjectRepository) *Service {
	return &Service{
		repo:        repo,
		subjectRepo: subjectRepo,
		now:         time.Now,
	}
}

func (s *Service) Start(userID, subjectID uint64) (*models.TimerState, error) {
	if userID == 0 || subjectID == 0 {
		return nil, ErrInvalidTimerInput
	}
	if _, err := s.subjectRepo.GetSubjectByID(context.Background(), subjectID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubjectNotFound
		}
		return nil, err
	}

	now := s.now()
	var state models.TimerState

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		placeholder := &models.TimerState{UserID: userID, Status: "IDLE"}
		if err := tx.WithContext(context.Background()).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).Create(placeholder).Error; err != nil {
			return err
		}

		result := tx.WithContext(context.Background()).Model(&models.TimerState{}).
			Where("user_id = ? AND status <> ?", userID, "RUNNING").
			Updates(map[string]any{
				"status":          "RUNNING",
				"subject_id":      subjectID,
				"started_at":      now,
				"last_resumed_at": now,
				"paused_seconds":  0,
				"draft_note":      nil,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrTimerAlreadyRunning
		}

		return tx.WithContext(context.Background()).Where("user_id = ?", userID).First(&state).Error
	}); err != nil {
		return nil, err
	}

	return &state, nil
}

func (s *Service) Stop(userID uint64, note *string) (*models.StudySession, error) {
	if userID == 0 {
		return nil, ErrInvalidTimerInput
	}

	var session *models.StudySession
	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		var state models.TimerState
		if err := tx.WithContext(context.Background()).Where("user_id = ?", userID).First(&state).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTimerNotRunning
			}
			return err
		}
		if state.Status != "RUNNING" || state.StartedAt == nil || state.SubjectID == nil {
			return ErrTimerNotRunning
		}

		endAt := s.now()
		startedAt := *state.StartedAt
		activeSeconds := int(endAt.Sub(startedAt).Seconds()) - state.PausedSeconds
		if activeSeconds < 0 {
			activeSeconds = 0
		}
		durationMinutes := activeSeconds / 60

		sessionNote := note
		if sessionNote == nil {
			sessionNote = state.DraftNote
		}

		claimed := tx.WithContext(context.Background()).Model(&models.TimerState{}).
			Where("id = ? AND user_id = ? AND status = ?", state.ID, userID, "RUNNING").
			Updates(map[string]any{
				"status":          "IDLE",
				"subject_id":      nil,
				"started_at":      nil,
				"last_resumed_at": nil,
				"paused_seconds":  0,
				"draft_note":      nil,
			})
		if claimed.Error != nil {
			return claimed.Error
		}
		if claimed.RowsAffected == 0 {
			return ErrTimerNotRunning
		}

		session = &models.StudySession{
			UserID:          userID,
			SubjectID:       *state.SubjectID,
			RecordType:      "TIMER",
			StartAt:         startedAt,
			EndAt:           endAt,
			DurationMinutes: durationMinutes,
			Note:            sessionNote,
		}
		return tx.WithContext(context.Background()).Create(session).Error
	}); err != nil {
		return nil, err
	}

	return session, nil
}
