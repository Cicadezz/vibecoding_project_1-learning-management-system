package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/http/router"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterAndLoginFlow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	}
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := database.Migrate(db); shouldSkipForCGO(err) {
		t.Skipf("sqlite unavailable in this environment: %v", err)
	} else if err != nil {
		t.Fatalf("migrate sqlite: %v", err)
	}

	r := router.NewRouter(db)

	registerBody := map[string]string{
		"username": "alice",
		"password": "pass1234",
	}
	registerPayload, _ := json.Marshal(registerBody)

	registerW := httptest.NewRecorder()
	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerPayload))
	registerReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(registerW, registerReq)

	if registerW.Code != http.StatusCreated {
		t.Fatalf("register expected status %d, got %d, body=%s", http.StatusCreated, registerW.Code, registerW.Body.String())
	}

	loginBody := map[string]string{
		"username": "alice",
		"password": "pass1234",
	}
	loginPayload, _ := json.Marshal(loginBody)

	loginW := httptest.NewRecorder()
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("login expected status %d, got %d, body=%s", http.StatusOK, loginW.Code, loginW.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(loginW.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if strings.TrimSpace(resp.Token) == "" {
		t.Fatal("expected non-empty token in login response")
	}
}

func shouldSkipForCGO(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "requires cgo") || strings.Contains(msg, "cgo_enabled=0")
}
