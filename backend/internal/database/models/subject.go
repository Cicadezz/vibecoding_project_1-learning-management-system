package models

import (
	"time"

	"gorm.io/datatypes"
)

type Subject struct {
	ID        uint64         `gorm:"primaryKey"`
	UserID    uint64         `gorm:"not null;uniqueIndex:uniq_user_subject_name,priority:1"`
	Name      string         `gorm:"size:64;not null;uniqueIndex:uniq_user_subject_name,priority:2"`
	Color     *string        `gorm:"size:16"`
	Ext       datatypes.JSON `gorm:"type:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
