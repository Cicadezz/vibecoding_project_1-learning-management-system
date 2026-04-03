package database

import (
	"learning-growth-platform/internal/database/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Subject{},
		&models.Task{},
		&models.StudySession{},
		&models.DailyCheckin{},
		&models.TimerState{},
	)
}
