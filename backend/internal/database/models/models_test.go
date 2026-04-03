package models

import (
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrateCreatesAllTables(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	err = db.AutoMigrate(
		&User{},
		&Subject{},
		&Task{},
		&StudySession{},
		&DailyCheckin{},
		&TimerState{},
	)
	if shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	expectedTables := []string{
		"users",
		"subjects",
		"tasks",
		"study_sessions",
		"daily_checkins",
		"timer_states",
	}

	for _, table := range expectedTables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("expected table %q to exist", table)
		}
	}
}

func shouldSkipForCGO(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
