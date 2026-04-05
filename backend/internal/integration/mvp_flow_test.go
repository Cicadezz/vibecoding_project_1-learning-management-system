package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/http/router"

	"gorm.io/gorm"
)

const defaultMySQLDSN = "root:010511@tcp(127.0.0.1:3306)/learning_growth?charset=utf8mb4&parseTime=True&loc=Local&timeout=2s&readTimeout=2s&writeTimeout=2s"

func TestMVPFlow(t *testing.T) {
	db := openIntegrationDB(t)
	if db == nil {
		t.Skip("mysql integration dependency unavailable")
	}
	defer closeDB(t, db)

	r := router.NewRouter(db)
	now := time.Now().In(time.Local)
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	unique := time.Now().UnixNano()
	username := fmt.Sprintf("mvp-%d", unique)
	password := "pass1234"

	registerResp := authFlow(t, r, http.MethodPost, "/api/auth/register", username, password)
	if registerResp.Token == "" || registerResp.User.ID == 0 {
		t.Fatalf("unexpected register response: %+v", registerResp)
	}

	loginResp := authFlow(t, r, http.MethodPost, "/api/auth/login", username, password)
	if loginResp.Token == "" || loginResp.User.ID == 0 {
		t.Fatalf("unexpected login response: %+v", loginResp)
	}
	if loginResp.User.ID != registerResp.User.ID {
		t.Fatalf("login user id %d did not match register user id %d", loginResp.User.ID, registerResp.User.ID)
	}

	subjectID := createSubject(t, r, loginResp.Token, fmt.Sprintf("math-%d", unique))
	createTask(t, r, loginResp.Token, day)
	createStudySession(t, r, loginResp.Token, subjectID, day)
	checkinToday(t, r, loginResp.Token)

	overview := fetchOverview(t, r, loginResp.Token)
	if overview.TodayMinutes != 45 {
		t.Fatalf("expected today_minutes 45, got %d", overview.TodayMinutes)
	}
	if overview.WeekMinutes != 45 {
		t.Fatalf("expected week_minutes 45, got %d", overview.WeekMinutes)
	}
	if overview.DoneTasks != 1 {
		t.Fatalf("expected done_tasks 1, got %d", overview.DoneTasks)
	}
	if overview.Streak != 1 {
		t.Fatalf("expected streak 1, got %d", overview.Streak)
	}
}

func openIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := strings.TrimSpace(os.Getenv("MYSQL_DSN"))
	if dsn == "" {
		dsn = defaultMySQLDSN
	}

	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		db, err := database.OpenMySQL(dsn)
		if err == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				pingErr = sqlDB.Ping()
			}
			if pingErr == nil {
				if err = database.Migrate(db); err == nil {
					return db
				}
			}
			if err == nil {
				err = pingErr
			}
		}

		lastErr = err
		if !isMySQLUnavailable(err) {
			t.Fatalf("open mysql integration database: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Skipf("mysql integration dependency unavailable: %v", lastErr)
	return nil
}

func closeDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("close db handle: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
}

type credentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authUserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type authResponse struct {
	Token string           `json:"token"`
	User  authUserResponse `json:"user"`
}

type taskCreateRequest struct {
	Title       string     `json:"title"`
	Priority    string     `json:"priority"`
	PlanDate    time.Time  `json:"plan_date"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at"`
}

type studySessionRequest struct {
	SubjectID uint64    `json:"subject_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
}

type overviewResponse struct {
	TodayMinutes int `json:"today_minutes"`
	WeekMinutes  int `json:"week_minutes"`
	DoneTasks    int `json:"done_tasks"`
	Streak       int `json:"streak"`
}

func authFlow(t *testing.T, r http.Handler, method, path, username, password string) authResponse {
	t.Helper()

	var resp authResponse
	postJSON(t, r, method, path, credentialRequest{Username: username, Password: password}, &resp, "")
	return resp
}

func createSubject(t *testing.T, r http.Handler, token, name string) uint64 {
	t.Helper()

	var resp struct {
		Subject struct {
			ID uint64 `json:"id"`
		} `json:"subject"`
	}
	postJSON(t, r, http.MethodPost, "/api/subjects", map[string]string{"name": name}, &resp, token)
	if resp.Subject.ID == 0 {
		t.Fatal("expected subject id")
	}
	return resp.Subject.ID
}

func createTask(t *testing.T, r http.Handler, token string, day time.Time) {
	t.Helper()

	completedAt := day.Add(11 * time.Hour)
	var resp struct {
		Task struct {
			ID uint64 `json:"id"`
		} `json:"task"`
	}
	postJSON(t, r, http.MethodPost, "/api/tasks", taskCreateRequest{
		Title:       "MVP task",
		Priority:    "MEDIUM",
		PlanDate:    day,
		Status:      "DONE",
		CompletedAt: &completedAt,
	}, &resp, token)
	if resp.Task.ID == 0 {
		t.Fatal("expected task id")
	}
}

func createStudySession(t *testing.T, r http.Handler, token string, subjectID uint64, day time.Time) {
	t.Helper()

	startAt := day.Add(9 * time.Hour)
	endAt := startAt.Add(45 * time.Minute)
	var resp struct {
		StudySession struct {
			ID uint64 `json:"id"`
		} `json:"study_session"`
	}
	postJSON(t, r, http.MethodPost, "/api/study/sessions", studySessionRequest{
		SubjectID: subjectID,
		StartAt:   startAt,
		EndAt:     endAt,
	}, &resp, token)
	if resp.StudySession.ID == 0 {
		t.Fatal("expected study session id")
	}
}

func checkinToday(t *testing.T, r http.Handler, token string) {
	t.Helper()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/checkin/today", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/checkin/today expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func fetchOverview(t *testing.T, r http.Handler, token string) overviewResponse {
	t.Helper()

	var resp struct {
		Overview overviewResponse `json:"overview"`
	}
	postJSON(t, r, http.MethodGet, "/api/stats/overview", nil, &resp, token)
	return resp.Overview
}

func postJSON(t *testing.T, r http.Handler, method, path string, body any, out any, token string) {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal %s request: %v", path, err)
		}
	}

	w := httptest.NewRecorder()
	var req *http.Request
	if body == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	r.ServeHTTP(w, req)

	if w.Code < 200 || w.Code >= 300 {
		t.Fatalf("%s expected 2xx, got %d, body=%s", path, w.Code, w.Body.String())
	}
	if out == nil {
		return
	}
	if err := json.Unmarshal(w.Body.Bytes(), out); err != nil {
		t.Fatalf("decode %s response: %v", path, err)
	}
}

func isMySQLUnavailable(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "cannot assign requested address") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "connectex") ||
		strings.Contains(msg, "server has gone away") ||
		strings.Contains(msg, "unknown database")
}
