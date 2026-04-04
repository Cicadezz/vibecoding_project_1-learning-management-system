package checkin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/database/models"
	"learning-growth-platform/internal/http/router"
	"learning-growth-platform/internal/subjects"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckinRoutes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGORouteCheckin(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGORouteCheckin(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	r := router.NewRouter(db)
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	registerBody := map[string]string{"username": "checkin-user", "password": "pass1234"}
	registerPayload, _ := json.Marshal(registerBody)

	registerW := httptest.NewRecorder()
	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerPayload))
	registerReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(registerW, registerReq)
	if registerW.Code != http.StatusCreated {
		t.Fatalf("register expected 201, got %d, body=%s", registerW.Code, registerW.Body.String())
	}

	loginBody := map[string]string{"username": "checkin-user", "password": "pass1234"}
	loginPayload, _ := json.Marshal(loginBody)

	loginW := httptest.NewRecorder()
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("login expected 200, got %d, body=%s", loginW.Code, loginW.Body.String())
	}

	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID uint64 `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(loginW.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if strings.TrimSpace(loginResp.Token) == "" || loginResp.User.ID == 0 {
		t.Fatalf("unexpected login response: %+v", loginResp)
	}

	subjectRepo := subjects.NewRepository(db)
	subjectID := createRouteSubjectForUser(t, subjectRepo, loginResp.User.ID, "math")
	if err := createRouteStudySession(t, db, loginResp.User.ID, subjectID, day); err != nil {
		t.Fatalf("create study session: %v", err)
	}

	todayW := httptest.NewRecorder()
	todayReq := httptest.NewRequest(http.MethodPost, "/api/checkin/today", nil)
	todayReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	r.ServeHTTP(todayW, todayReq)
	if todayW.Code != http.StatusOK {
		t.Fatalf("POST /api/checkin/today expected 200, got %d, body=%s", todayW.Code, todayW.Body.String())
	}

	streakW := httptest.NewRecorder()
	streakReq := httptest.NewRequest(http.MethodGet, "/api/checkin/streak", nil)
	streakReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	r.ServeHTTP(streakW, streakReq)
	if streakW.Code != http.StatusOK {
		t.Fatalf("GET /api/checkin/streak expected 200, got %d, body=%s", streakW.Code, streakW.Body.String())
	}

	var streakResp struct {
		Streak int `json:"streak"`
	}
	if err := json.Unmarshal(streakW.Body.Bytes(), &streakResp); err != nil {
		t.Fatalf("decode streak response: %v", err)
	}
	if streakResp.Streak != 1 {
		t.Fatalf("expected streak 1, got %d", streakResp.Streak)
	}
}

func createRouteSubjectForUser(t *testing.T, repo *subjects.Repository, userID uint64, name string) uint64 {
	t.Helper()

	subject := &models.Subject{
		UserID: userID,
		Name:   name,
	}
	if err := repo.CreateSubject(context.Background(), subject); err != nil {
		t.Fatalf("create subject: %v", err)
	}
	return subject.ID
}

func createRouteStudySession(t *testing.T, db *gorm.DB, userID, subjectID uint64, day time.Time) error {
	t.Helper()

	return db.Create(&models.StudySession{
		UserID:          userID,
		SubjectID:       subjectID,
		RecordType:      "MANUAL",
		StartAt:         day.Add(9 * time.Hour),
		EndAt:           day.Add(10 * time.Hour),
		DurationMinutes: 60,
	}).Error
}

func shouldSkipForCGORouteCheckin(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
