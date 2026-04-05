package stats

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"
)

var ErrInvalidStatsInput = errors.New("invalid stats input")

type Service struct {
	repo *Repository
	now  func() time.Time
}

type Overview struct {
	TodayMinutes int `json:"today_minutes"`
	WeekMinutes  int `json:"week_minutes"`
	DoneTasks    int `json:"done_tasks"`
	Streak       int `json:"streak"`
}

type WeeklyTrendItem struct {
	Date    string `json:"date"`
	Minutes int    `json:"minutes"`
}

type SubjectDistributionItem struct {
	SubjectID   uint64 `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	Minutes     int    `json:"minutes"`
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) Overview(userID uint64) (*Overview, error) {
	if userID == 0 {
		return nil, ErrInvalidStatsInput
	}

	now := s.now()
	todayStart := normalizeDay(now)
	weekStart := mondayStart(now)
	weekEnd := weekStart.AddDate(0, 0, 7)

	rows, err := s.repo.ListStudySessionRowsBetween(context.Background(), userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	overview := &Overview{
		TodayMinutes: sumRowsForDay(rows, todayStart),
		WeekMinutes:  sumRowsInRange(rows, weekStart, weekEnd),
	}

	doneTasks, err := s.repo.CountDoneTasks(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	overview.DoneTasks = doneTasks

	streak, err := s.streak(userID, now)
	if err != nil {
		return nil, err
	}
	overview.Streak = streak

	return overview, nil
}

func (s *Service) WeeklyTrend(userID uint64) ([]WeeklyTrendItem, error) {
	if userID == 0 {
		return nil, ErrInvalidStatsInput
	}

	now := s.now()
	weekStart := mondayStart(now)
	weekEnd := weekStart.AddDate(0, 0, 7)

	rows, err := s.repo.ListStudySessionRowsBetween(context.Background(), userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	trend := make([]WeeklyTrendItem, 0, 7)
	for i := 0; i < 7; i++ {
		day := weekStart.AddDate(0, 0, i)
		dayEnd := day.AddDate(0, 0, 1)
		trend = append(trend, WeeklyTrendItem{
			Date:    day.Format("2006-01-02"),
			Minutes: sumRowsInRange(rows, day, dayEnd),
		})
	}
	return trend, nil
}

func (s *Service) SubjectDistribution(userID uint64) ([]SubjectDistributionItem, error) {
	if userID == 0 {
		return nil, ErrInvalidStatsInput
	}

	now := s.now()
	weekStart := mondayStart(now)
	weekEnd := weekStart.AddDate(0, 0, 7)

	rows, err := s.repo.ListStudySessionRowsBetween(context.Background(), userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	type aggregate struct {
		SubjectID   uint64
		SubjectName string
		Minutes     int
	}

	bySubject := map[uint64]*aggregate{}
	for _, row := range rows {
		item, ok := bySubject[row.SubjectID]
		if !ok {
			item = &aggregate{SubjectID: row.SubjectID, SubjectName: row.SubjectName}
			bySubject[row.SubjectID] = item
		}
		item.Minutes += minutesInRange(row, weekStart, weekEnd)
	}

	distribution := make([]SubjectDistributionItem, 0, len(bySubject))
	for _, item := range bySubject {
		distribution = append(distribution, SubjectDistributionItem{
			SubjectID:   item.SubjectID,
			SubjectName: item.SubjectName,
			Minutes:     item.Minutes,
		})
	}

	sort.Slice(distribution, func(i, j int) bool {
		if distribution[i].Minutes != distribution[j].Minutes {
			return distribution[i].Minutes > distribution[j].Minutes
		}
		if distribution[i].SubjectName != distribution[j].SubjectName {
			return distribution[i].SubjectName < distribution[j].SubjectName
		}
		return distribution[i].SubjectID < distribution[j].SubjectID
	})

	return distribution, nil
}

func (s *Service) streak(userID uint64, now time.Time) (int, error) {
	today := normalizeDay(now)
	hasToday, err := s.repo.HasCheckinOnDate(context.Background(), userID, today)
	if err != nil {
		return 0, err
	}
	if !hasToday {
		return 0, nil
	}

	streak := 1
	cursor := today.AddDate(0, 0, -1)
	for {
		hasCheckin, err := s.repo.HasCheckinOnDate(context.Background(), userID, cursor)
		if err != nil {
			return 0, err
		}
		if !hasCheckin {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}

	return streak, nil
}

func sumRows(rows []StudySessionRow) int {
	total := 0
	for _, row := range rows {
		total += row.DurationMinutes
	}
	return total
}

func sumRowsForDay(rows []StudySessionRow, day time.Time) int {
	return sumRowsInRange(rows, day, day.AddDate(0, 0, 1))
}

func sumRowsInRange(rows []StudySessionRow, start, end time.Time) int {
	total := 0
	for _, row := range rows {
		total += minutesInRange(row, start, end)
	}
	return total
}

func minutesInRange(row StudySessionRow, start, end time.Time) int {
	if row.EndAt.IsZero() || row.StartAt.IsZero() {
		return 0
	}

	intersectStart := row.StartAt
	if start.After(intersectStart) {
		intersectStart = start
	}
	intersectEnd := row.EndAt
	if end.Before(intersectEnd) {
		intersectEnd = end
	}
	if !intersectEnd.After(intersectStart) {
		return 0
	}

	elapsed := row.EndAt.Sub(row.StartAt)
	if elapsed <= 0 {
		return 0
	}

	overlap := intersectEnd.Sub(intersectStart)
	if overlap >= elapsed {
		return row.DurationMinutes
	}

	return int(math.Floor(float64(row.DurationMinutes) * overlap.Seconds() / elapsed.Seconds()))
}

