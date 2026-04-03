package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           uint64         `gorm:"primaryKey"`
	Username     string         `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Ext          datatypes.JSON `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
