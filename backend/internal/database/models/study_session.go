package models

import (
	"time"

	"gorm.io/datatypes"
)

type StudySession struct {
	ID              uint64         `gorm:"primaryKey"`
	UserID          uint64         `gorm:"not null;index:idx_study_user_start"`
	SubjectID       uint64         `gorm:"not null;index:idx_study_user_subject_start,priority:2"`
	RecordType      string         `gorm:"type:enum('MANUAL','TIMER');not null"`
	StartAt         time.Time      `gorm:"not null;index:idx_study_user_start;index:idx_study_user_subject_start,priority:3"`
	EndAt           time.Time      `gorm:"not null"`
	DurationMinutes int            `gorm:"not null"`
	Note            *string        `gorm:"size:1000"`
	Ext             datatypes.JSON `gorm:"type:json"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
