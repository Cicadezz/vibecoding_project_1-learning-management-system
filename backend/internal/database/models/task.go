package models

import (
	"time"

	"gorm.io/datatypes"
)

type Task struct {
	ID          uint64         `gorm:"primaryKey"`
	UserID      uint64         `gorm:"not null;index:idx_tasks_user_plan_status,priority:1"`
	Title       string         `gorm:"size:255;not null"`
	Priority    string         `gorm:"type:enum('HIGH','MEDIUM','LOW');default:'MEDIUM';not null"`
	DueDate     *time.Time     `gorm:"type:date;index:idx_tasks_user_due_date"`
	PlanDate    time.Time      `gorm:"type:date;not null;index:idx_tasks_user_plan_status,priority:2"`
	Status      string         `gorm:"type:enum('PENDING','DONE');default:'PENDING';not null;index:idx_tasks_user_plan_status,priority:3"`
	CompletedAt *time.Time
	CarryCount  int            `gorm:"not null;default:0"`
	Ext         datatypes.JSON `gorm:"type:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
