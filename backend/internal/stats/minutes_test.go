package stats

import (
	"testing"
	"time"
)

func TestMinutesInRangeUsesDurationMinutes(t *testing.T) {
	session := StudySessionRow{
		StartAt:         time.Date(2026, 4, 4, 10, 0, 0, 0, time.Local),
		EndAt:           time.Date(2026, 4, 4, 11, 0, 0, 0, time.Local),
		DurationMinutes: 40,
	}

	t.Run("full overlap uses session duration", func(t *testing.T) {
		got := minutesInRange(session, session.StartAt, session.EndAt)
		if got != 40 {
			t.Fatalf("expected 40 minutes, got %d", got)
		}
	})

	t.Run("partial overlap is proportional to elapsed time", func(t *testing.T) {
		start := time.Date(2026, 4, 4, 10, 30, 0, 0, time.Local)
		got := minutesInRange(session, start, session.EndAt)
		if got != 20 {
			t.Fatalf("expected 20 minutes, got %d", got)
		}
	})
}
