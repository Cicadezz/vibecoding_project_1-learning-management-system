package router

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthRoute(t *testing.T) {
    r := NewRouter(nil)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
}
