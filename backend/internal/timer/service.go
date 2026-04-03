package timer

import (
	"context"
	"errors"
	"time"

	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

var ErrInvalidTimerInput = errors.New("invalid timer input")
var ErrTimerNotRunning = errors.New("timer not running")
var ErrTimerAlreadyRunning = errors.New("timer already running")

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

func (s *Service) Start(userID, subjectID uint64) (*models.TimerState, error) {
	if userID == 0 || subjectID == 0 {
		return nil, ErrInvalidTimerInput
	}

	state, err := s.repo.GetStateByUser(context.Background(), userID)
	if err != nil {
		if !errors.Is(err, ErrTimerStateNotFound) {
			return nil, err
		}
		state = &models.TimerState{UserID: userID}
	}

	if state.Status == "RUNNING" {
		return nil, ErrTimerAlreadyRunning
	}

	now := s.now()
	state.Status = "RUNNING"
	subject := subjectID
	state.SubjectID = &subject
	state.StartedAt = &now
	state.LastResumedAt = &now
	state.PausedSeconds = 0
	state.DraftNote = nil

	if err := s.repo.SaveState(context.Background(), state); err != nil {
		return nil, err
	}
	return state, nil
}

func (s *Service) Stop(userID uint64, note *string) (*models.StudySession, error) {
	if userID == 0 {
		return nil, ErrInvalidTimerInput
	}

	state, err := s.repo.GetStateByUser(context.Background(), userID)
	if err != nil {
		if errors.Is(err, ErrTimerStateNotFound) {
			return nil, ErrTimerNotRunning
		}
		return nil, err
	}
	if state.Status != "RUNNING" || state.StartedAt == nil || state.SubjectID == nil {
		return nil, ErrTimerNotRunning
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

	session := &models.StudySession{
		UserID:          userID,
		SubjectID:       *state.SubjectID,
		RecordType:      "TIMER",
		StartAt:         startedAt,
		EndAt:           endAt,
		DurationMinutes: durationMinutes,
		Note:            sessionNote,
	}

	if err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(context.Background()).Create(session).Error; err != nil {
			return err
		}

		state.Status = "IDLE"
		state.SubjectID = nil
		state.StartedAt = nil
		state.LastResumedAt = nil
		state.PausedSeconds = 0
		state.DraftNote = nil
		return tx.WithContext(context.Background()).Save(state).Error
	}); err != nil {
		return nil, err
	}

	return session, nil
}
