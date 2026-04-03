package models

import (
	"time"

	"gorm.io/datatypes"
)

type TimerState struct {
	ID           uint64         `gorm:"primaryKey"`
	UserID       uint64         `gorm:"not null;uniqueIndex"`
	Status       string         `gorm:"type:enum('IDLE','RUNNING','PAUSED');not null;default:'IDLE'"`
	SubjectID    *uint64
	StartedAt    *time.Time
	LastResumedAt *time.Time
	PausedSeconds int           `gorm:"not null;default:0"`
	DraftNote    *string        `gorm:"size:1000"`
	Ext          datatypes.JSON `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
