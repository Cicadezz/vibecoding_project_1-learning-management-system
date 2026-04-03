package subjects

import (
	"context"
	"errors"
	"strings"

	"learning-growth-platform/internal/database/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var ErrInvalidSubjectInput = errors.New("invalid subject input")

// Service keeps subject operations thin and repository-backed.
type Service struct {
	repo *Repository
}

type CreateSubjectInput struct {
	UserID uint64
	Name   string
	Color  *string
	Ext    []byte
}

type UpdateSubjectInput struct {
	Name  *string
	Color *string
	Ext   []byte
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateSubjectInput) (*models.Subject, error) {
	name := strings.TrimSpace(input.Name)
	if input.UserID == 0 || name == "" {
		return nil, ErrInvalidSubjectInput
	}

	subject := &models.Subject{
		UserID: input.UserID,
		Name:   name,
		Color:  input.Color,
		Ext:    datatypes.JSON(input.Ext),
	}
	if err := s.repo.CreateSubject(context.Background(), subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *Service) ListByUser(userID uint64) ([]models.Subject, error) {
	return s.repo.ListSubjects(context.Background(), userID)
}

func (s *Service) Update(subjectID, userID uint64, input UpdateSubjectInput) (*models.Subject, error) {
	if subjectID == 0 || userID == 0 {
		return nil, ErrInvalidSubjectInput
	}

	subject, err := s.repo.GetSubjectByID(context.Background(), subjectID, userID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidSubjectInput
		}
		subject.Name = name
	}
	if input.Color != nil {
		subject.Color = input.Color
	}
	if input.Ext != nil {
		subject.Ext = datatypes.JSON(input.Ext)
	}

	if err := s.repo.UpdateSubject(context.Background(), subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *Service) Delete(subjectID, userID uint64) error {
	if subjectID == 0 || userID == 0 {
		return ErrInvalidSubjectInput
	}
	return s.repo.DeleteSubject(context.Background(), subjectID, userID)
}

func isSubjectNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
