package models

import (
	"time"

	"gorm.io/datatypes"
)

type DailyCheckin struct {
	ID         uint64         `gorm:"primaryKey"`
	UserID     uint64         `gorm:"not null;uniqueIndex:uniq_user_checkin_date,priority:1"`
	CheckinDate time.Time     `gorm:"type:date;not null;uniqueIndex:uniq_user_checkin_date,priority:2"`
	CheckedAt  time.Time      `gorm:"not null"`
	Ext        datatypes.JSON `gorm:"type:json"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
