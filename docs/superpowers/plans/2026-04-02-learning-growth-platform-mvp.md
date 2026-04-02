# Learning Growth Platform MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a local Web learning-growth tool that connects planning, execution, and review: tasks, study sessions, daily check-in, and stats.

**Architecture:** Backend uses Go (Gin) + Gorm + MySQL 8.0.27 for local APIs and persistence; frontend uses React + TypeScript + Vite + ECharts for responsive UI. Implement backend domain modules first, then wire frontend pages, and finish with end-to-end verification.

**Tech Stack:** Go 1.24, Gin, Gorm, MySQL 8.0.27, React, TypeScript, Vite, Vitest, ECharts, Docker Compose

---

## File Structure

- `backend/cmd/server/main.go`: application entrypoint
- `backend/internal/config/config.go`: env config loader
- `backend/internal/database/mysql.go`: MySQL connection
- `backend/internal/database/migrate.go`: migrations and index setup
- `backend/internal/database/models/user.go`: user entity
- `backend/internal/database/models/subject.go`: subject entity
- `backend/internal/database/models/task.go`: task entity
- `backend/internal/database/models/study_session.go`: study session entity
- `backend/internal/database/models/daily_checkin.go`: daily check-in entity
- `backend/internal/database/models/timer_state.go`: timer state entity
- `backend/internal/auth/handler.go`: auth endpoints
- `backend/internal/subjects/handler.go`: subject endpoints
- `backend/internal/tasks/handler.go`: task endpoints
- `backend/internal/study/handler.go`: study endpoints
- `backend/internal/timer/handler.go`: timer endpoints
- `backend/internal/checkin/handler.go`: check-in endpoints
- `backend/internal/stats/handler.go`: stats endpoints
- `frontend/src/pages/LoginPage.tsx`: login screen
- `frontend/src/pages/RegisterPage.tsx`: register screen
- `frontend/src/pages/DashboardPage.tsx`: overview metrics and charts
- `frontend/src/pages/TasksPage.tsx`: today's tasks module
- `frontend/src/pages/StudyPage.tsx`: manual/timer study records
- `frontend/src/pages/CheckinPage.tsx`: daily check-in and streak
- `frontend/src/pages/SettingsPage.tsx`: account settings
- `frontend/src/components/charts/WeeklyTrendChart.tsx`: 7-day trend chart
- `frontend/src/components/charts/SubjectPieChart.tsx`: subject distribution chart
- `docker-compose.yml`: local MySQL 8.0.27 runtime
- `README.md`: runbook and verification commands

### Task 1: Backend Skeleton and Health Endpoint

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/server/main.go`
- Create: `backend/internal/http/router/router.go`
- Test: `backend/internal/http/router/router_test.go`

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/http/router -run TestHealthRoute -v`  
Expected: FAIL with `undefined: NewRouter`

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/http/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}
```

```go
// backend/cmd/server/main.go
package main

import (
	"log"
	"learning-growth-platform/internal/http/router"
)

func main() {
	r := router.NewRouter(nil)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/http/router -run TestHealthRoute -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/go.mod backend/cmd/server/main.go backend/internal/http/router/router.go backend/internal/http/router/router_test.go
git commit -m "chore: bootstrap backend skeleton"
```

### Task 2: MySQL Config, Models, and Migration

**Files:**
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/database/mysql.go`
- Create: `backend/internal/database/migrate.go`
- Create: `backend/internal/database/models/user.go`
- Create: `backend/internal/database/models/subject.go`
- Create: `backend/internal/database/models/task.go`
- Create: `backend/internal/database/models/study_session.go`
- Create: `backend/internal/database/models/daily_checkin.go`
- Create: `backend/internal/database/models/timer_state.go`
- Create: `backend/.env.example`
- Test: `backend/internal/database/models/models_test.go`

- [ ] **Step 1: Write the failing test**

```go
package models

import (
	"testing"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCoreTablesExist(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil { t.Fatal(err) }

	err = db.AutoMigrate(&User{}, &Subject{}, &Task{}, &StudySession{}, &DailyCheckin{}, &TimerState{})
	if err != nil { t.Fatal(err) }

	for _, table := range []string{"users", "subjects", "tasks", "study_sessions", "daily_checkins", "timer_states"} {
		if !db.Migrator().HasTable(table) {
			t.Fatalf("expected table %s", table)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/database/models -run TestCoreTablesExist -v`  
Expected: FAIL with `undefined: User`

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/database/models/user.go
package models

import (
	"time"
	"gorm.io/datatypes"
)

type User struct {
	ID           uint64         `gorm:"primaryKey"`
	Username     string         `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Ext          datatypes.JSON `gorm:"type:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
```

```go
// backend/internal/database/models/task.go
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
	DueDate     *time.Time
	PlanDate    time.Time      `gorm:"type:date;not null;index:idx_tasks_user_plan_status,priority:2"`
	Status      string         `gorm:"type:enum('PENDING','DONE');default:'PENDING';not null;index:idx_tasks_user_plan_status,priority:3"`
	CompletedAt *time.Time
	CarryCount  int            `gorm:"default:0;not null"`
	Ext         datatypes.JSON `gorm:"type:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

```go
// backend/internal/database/migrate.go
package database

import (
	"learning-growth-platform/internal/database/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Subject{}, &models.Task{}, &models.StudySession{}, &models.DailyCheckin{}, &models.TimerState{})
}
```

```env
# backend/.env.example
APP_PORT=8080
MYSQL_DSN=root:010511@tcp(127.0.0.1:3306)/learning_growth?parseTime=true&loc=Local&charset=utf8mb4
JWT_SECRET=local-dev-secret
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/database/models -run TestCoreTablesExist -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/config backend/internal/database backend/.env.example
git commit -m "feat: add mysql config models and migrations"
```

### Task 3: Auth APIs (Register/Login/Me/Change Password)

**Files:**
- Create: `backend/internal/http/middleware/auth.go`
- Create: `backend/internal/auth/repository.go`
- Create: `backend/internal/auth/service.go`
- Create: `backend/internal/auth/handler.go`
- Modify: `backend/internal/http/router/router.go`
- Test: `backend/internal/auth/auth_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestRegisterLoginAndMe(t *testing.T) {
	r := buildAuthTestRouter(t)
	register := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username":"alice","password":"pass1234"}`))
	register.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, register)
	if w1.Code != http.StatusCreated { t.Fatalf("expected 201, got %d", w1.Code) }

	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"alice","password":"pass1234"}`))
	login.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, login)
	if w2.Code != http.StatusOK { t.Fatalf("expected 200, got %d", w2.Code) }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/auth -run TestRegisterLoginAndMe -v`  
Expected: FAIL with missing handler/service symbols

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/auth/handler.go
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.GET("/me", h.Me)
	rg.POST("/change-password", h.ChangePassword)
}
```

```go
// backend/internal/auth/service.go
func (s *Service) Register(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return err }
	return s.repo.Create(username, string(hash))
}
func (s *Service) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil { return "", err }
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredential
	}
	return s.tokenSigner.Sign(user.ID, user.Username)
}
func (s *Service) ChangePassword(userID uint64, oldPwd, newPwd string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil { return err }
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPwd)) != nil {
		return ErrInvalidCredential
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil { return err }
	return s.repo.UpdatePassword(userID, string(hash))
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/auth -run TestRegisterLoginAndMe -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/auth backend/internal/http/middleware/auth.go backend/internal/http/router/router.go
git commit -m "feat: implement local auth endpoints"
```

### Task 4: Subjects and Tasks Modules (CRUD + Carry-Over)

**Files:**
- Create: `backend/internal/subjects/repository.go`
- Create: `backend/internal/subjects/service.go`
- Create: `backend/internal/subjects/handler.go`
- Create: `backend/internal/tasks/repository.go`
- Create: `backend/internal/tasks/service.go`
- Create: `backend/internal/tasks/handler.go`
- Modify: `backend/internal/auth/service.go`
- Modify: `backend/internal/http/router/router.go`
- Test: `backend/internal/tasks/tasks_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestCarryOverPendingTasksIsIdempotent(t *testing.T) {
	svc := buildTaskServiceForTest(t)
	userID := uint64(1)
	yesterday := time.Date(2026, 4, 1, 0, 0, 0, 0, time.Local)
	today := time.Date(2026, 4, 2, 0, 0, 0, 0, time.Local)

	_, _ = svc.Create(CreateTaskInput{UserID: userID, Title: "刷题", PlanDate: yesterday})
	_, _ = svc.CarryOverPendingTasks(userID, today)
	_, _ = svc.CarryOverPendingTasks(userID, today)

	list, _ := svc.ListByDate(userID, today)
	if len(list) != 1 || list[0].CarryCount != 1 {
		t.Fatalf("expected exactly one carried task with carry_count=1")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/tasks -run TestCarryOverPendingTasksIsIdempotent -v`  
Expected: FAIL with missing `CarryOverPendingTasks`

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/tasks/service.go
func (s *Service) CarryOverPendingTasks(userID uint64, today time.Time) (int64, error) {
	yesterday := today.AddDate(0, 0, -1)
	return s.repo.CarryOverPending(userID, yesterday, today)
}
```

```go
// backend/internal/tasks/repository.go
func (r *Repository) CarryOverPending(userID uint64, fromDate, toDate time.Time) (int64, error) {
	res := r.db.Model(&models.Task{}).
		Where("user_id=? AND plan_date=? AND status='PENDING'", userID, fromDate.Format("2006-01-02")).
		Updates(map[string]any{"plan_date": toDate.Format("2006-01-02"), "carry_count": gorm.Expr("carry_count + 1")})
	return res.RowsAffected, res.Error
}
```

```go
// backend/internal/auth/service.go (after successful login)
_, _ = s.taskSvc.CarryOverPendingTasks(user.ID, time.Now())
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/tasks -run TestCarryOverPendingTasksIsIdempotent -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/subjects backend/internal/tasks backend/internal/auth/service.go backend/internal/http/router/router.go
git commit -m "feat: implement subjects and tasks with carry-over"
```

### Task 5: Study Sessions and Timer Module

**Files:**
- Create: `backend/internal/study/repository.go`
- Create: `backend/internal/study/service.go`
- Create: `backend/internal/study/handler.go`
- Create: `backend/internal/timer/repository.go`
- Create: `backend/internal/timer/service.go`
- Create: `backend/internal/timer/handler.go`
- Modify: `backend/internal/http/router/router.go`
- Test: `backend/internal/timer/timer_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestStopTimerCreatesSession(t *testing.T) {
	svc := buildTimerServiceForTest(t)
	_, _ = svc.Start(StartInput{UserID: 1, SubjectID: 1, StartedAt: time.Date(2026, 4, 2, 10, 0, 0, 0, time.Local)})
	session, err := svc.Stop(StopInput{UserID: 1, EndAt: time.Date(2026, 4, 2, 11, 0, 0, 0, time.Local), Note: "背单词"})
	if err != nil { t.Fatal(err) }
	if session.DurationMinutes != 60 { t.Fatalf("expected 60, got %d", session.DurationMinutes) }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/timer -run TestStopTimerCreatesSession -v`  
Expected: FAIL with missing timer methods

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/study/service.go
func (s *Service) CreateManual(input CreateManualInput) error {
	if !input.EndAt.After(input.StartAt) || input.DurationMinutes <= 0 { return ErrInvalidTimeRange }
	return s.repo.CreateManual(input)
}
```

```go
// backend/internal/timer/service.go
func (s *Service) Stop(input StopInput) (SessionDTO, error) {
	state, err := s.repo.GetByUser(input.UserID)
	if err != nil { return SessionDTO{}, err }
	duration := int(input.EndAt.Sub(state.StartedAt).Minutes()) - int(state.PausedSeconds/60)
	session, err := s.studySvc.CreateFromTimer(state.UserID, state.SubjectID, state.StartedAt, input.EndAt, duration, input.Note)
	if err != nil { return SessionDTO{}, err }
	_ = s.repo.Reset(input.UserID)
	return session, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/timer -run TestStopTimerCreatesSession -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/study backend/internal/timer backend/internal/http/router/router.go
git commit -m "feat: add study sessions and timer workflow"
```

### Task 6: Daily Check-In and Streak

**Files:**
- Create: `backend/internal/checkin/repository.go`
- Create: `backend/internal/checkin/service.go`
- Create: `backend/internal/checkin/handler.go`
- Modify: `backend/internal/http/router/router.go`
- Test: `backend/internal/checkin/checkin_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestCheckinRequiresStudyAndIsIdempotent(t *testing.T) {
	svc := buildCheckinServiceForTest(t)
	today := time.Date(2026, 4, 2, 0, 0, 0, 0, time.Local)
	if _, err := svc.CheckinToday(1, today); err == nil { t.Fatal("expected prerequisite error") }
	insertStudyForDate(t, svc.db, 1, today)
	if _, err := svc.CheckinToday(1, today); err != nil { t.Fatal(err) }
	if _, err := svc.CheckinToday(1, today); err != nil { t.Fatal(err) }
	streak, _ := svc.GetStreak(1, today)
	if streak != 1 { t.Fatalf("expected streak 1, got %d", streak) }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/checkin -run TestCheckinRequiresStudyAndIsIdempotent -v`  
Expected: FAIL with missing checkin service methods

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/checkin/service.go
func (s *Service) CheckinToday(userID uint64, day time.Time) (CheckinDTO, error) {
	has, err := s.repo.HasStudySessionOnDate(userID, day)
	if err != nil { return CheckinDTO{}, err }
	if !has { return CheckinDTO{}, ErrNoStudySession }
	return s.repo.UpsertCheckin(userID, day)
}
```

```go
// backend/internal/checkin/handler.go
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/today", h.CheckinToday)
	rg.GET("/streak", h.GetStreak)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/checkin -run TestCheckinRequiresStudyAndIsIdempotent -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/checkin backend/internal/http/router/router.go
git commit -m "feat: implement checkin and streak rules"
```

### Task 7: Stats APIs (Monday Week Start)

**Files:**
- Create: `backend/internal/stats/repository.go`
- Create: `backend/internal/stats/service.go`
- Create: `backend/internal/stats/handler.go`
- Modify: `backend/internal/http/router/router.go`
- Test: `backend/internal/stats/stats_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestOverviewUsesMondayWeekStart(t *testing.T) {
	svc := buildStatsServiceForTest(t)
	ref := time.Date(2026, 4, 2, 12, 0, 0, 0, time.Local)
	insertSession(t, svc.db, 1, time.Date(2026, 3, 30, 9, 0, 0, 0, time.Local), 60)
	insertSession(t, svc.db, 1, time.Date(2026, 3, 29, 9, 0, 0, 0, time.Local), 90)
	o, err := svc.Overview(1, ref)
	if err != nil { t.Fatal(err) }
	if o.WeekMinutes != 60 { t.Fatalf("expected 60, got %d", o.WeekMinutes) }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/stats -run TestOverviewUsesMondayWeekStart -v`  
Expected: FAIL with missing stats implementation

- [ ] **Step 3: Write minimal implementation**

```go
// backend/internal/stats/service.go
func weekRangeMonday(ref time.Time) (time.Time, time.Time) {
	wd := int(ref.Weekday())
	if wd == 0 { wd = 7 }
	start := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, ref.Location()).AddDate(0, 0, -(wd-1))
	return start, start.AddDate(0, 0, 7)
}

func (s *Service) Overview(userID uint64, ref time.Time) (OverviewDTO, error) {
	start, end := weekRangeMonday(ref)
	return s.repo.LoadOverview(userID, ref, start, end)
}
```

```go
// backend/internal/stats/handler.go
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/overview", h.Overview)
	rg.GET("/weekly-trend", h.WeeklyTrend)
	rg.GET("/subject-distribution", h.SubjectDistribution)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/stats -run TestOverviewUsesMondayWeekStart -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add backend/internal/stats backend/internal/http/router/router.go
git commit -m "feat: add monday-based stats apis"
```

### Task 8: Frontend Scaffold and Auth Pages

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/App.tsx`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/auth.ts`
- Create: `frontend/src/pages/LoginPage.tsx`
- Create: `frontend/src/pages/RegisterPage.tsx`
- Test: `frontend/src/__tests__/auth-pages.test.tsx`

- [ ] **Step 1: Write the failing test**

```tsx
test("login submits username and password", async () => {
  const onSubmit = vi.fn().mockResolvedValue(undefined);
  render(<LoginPage onSubmit={onSubmit} />);
  await userEvent.type(screen.getByLabelText("用户名"), "alice");
  await userEvent.type(screen.getByLabelText("密码"), "pass1234");
  await userEvent.click(screen.getByRole("button", { name: "登录" }));
  expect(onSubmit).toHaveBeenCalledWith({ username: "alice", password: "pass1234" });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/frontend && npm test -- --run auth-pages`  
Expected: FAIL with missing page component

- [ ] **Step 3: Write minimal implementation**

```tsx
// frontend/src/pages/LoginPage.tsx
export function LoginPage({ onSubmit }: { onSubmit?: (v: { username: string; password: string }) => Promise<void> }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  return (
    <form onSubmit={(e) => { e.preventDefault(); void onSubmit?.({ username, password }); }}>
      <label>用户名</label>
      <input aria-label="用户名" value={username} onChange={(e) => setUsername(e.target.value)} />
      <label>密码</label>
      <input aria-label="密码" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button type="submit">登录</button>
    </form>
  );
}
```

```ts
// frontend/src/lib/api.ts
export async function api<T>(url: string, init?: RequestInit): Promise<T> {
  const resp = await fetch(url, init);
  if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
  return resp.json() as Promise<T>;
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/frontend && npm test -- --run auth-pages`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add frontend
git commit -m "feat: scaffold frontend and auth pages"
```

### Task 9: Frontend Core Pages and Charts

**Files:**
- Create: `frontend/src/pages/DashboardPage.tsx`
- Create: `frontend/src/pages/TasksPage.tsx`
- Create: `frontend/src/pages/StudyPage.tsx`
- Create: `frontend/src/pages/CheckinPage.tsx`
- Create: `frontend/src/pages/SettingsPage.tsx`
- Create: `frontend/src/components/charts/WeeklyTrendChart.tsx`
- Create: `frontend/src/components/charts/SubjectPieChart.tsx`
- Modify: `frontend/src/App.tsx`
- Test: `frontend/src/__tests__/dashboard-page.test.tsx`

- [ ] **Step 1: Write the failing test**

```tsx
test("dashboard renders overview metrics", () => {
  render(<DashboardPage overview={{ todayMinutes: 120, weekMinutes: 480, doneTasks: 3, streak: 5 }} trend={[]} subjects={[]} />);
  expect(screen.getByText("今日学习总时长")).toBeInTheDocument();
  expect(screen.getByText("120 分钟")).toBeInTheDocument();
  expect(screen.getByText("连续打卡 5 天")).toBeInTheDocument();
});
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/frontend && npm test -- --run dashboard-page`  
Expected: FAIL with missing `DashboardPage`

- [ ] **Step 3: Write minimal implementation**

```tsx
// frontend/src/pages/DashboardPage.tsx
export function DashboardPage({ overview }: { overview: { todayMinutes: number; weekMinutes: number; doneTasks: number; streak: number } }) {
  return (
    <section>
      <h1>统计面板</h1>
      <p>今日学习总时长</p>
      <p>{overview.todayMinutes} 分钟</p>
      <p>本周学习总时长</p>
      <p>{overview.weekMinutes} 分钟</p>
      <p>今日完成任务数 {overview.doneTasks}</p>
      <p>连续打卡 {overview.streak} 天</p>
    </section>
  );
}
```

```tsx
// frontend/src/App.tsx
export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/dashboard" element={<DashboardPage overview={{ todayMinutes: 0, weekMinutes: 0, doneTasks: 0, streak: 0 }} trend={[]} subjects={[]} />} />
        <Route path="/tasks" element={<TasksPage />} />
        <Route path="/study" element={<StudyPage />} />
        <Route path="/checkin" element={<CheckinPage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Routes>
    </BrowserRouter>
  );
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/frontend && npm test -- --run dashboard-page`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add frontend/src/App.tsx frontend/src/pages frontend/src/components/charts
git commit -m "feat: add core pages and dashboard charts"
```

### Task 10: Runtime Setup, Integration Test, and Documentation

**Files:**
- Create: `docker-compose.yml`
- Create: `backend/internal/integration/mvp_flow_test.go`
- Create: `README.md`

- [ ] **Step 1: Write the failing test**

```go
func TestMVPFlow(t *testing.T) {
  c := newAPIClient(t)
  token := c.RegisterAndLogin("alice", "pass1234")
  c.CreateTask(token, "英语单词", "HIGH")
  c.CreateManualStudy(token, "英语", 50, "背单词")
  c.CheckinToday(token)
  overview := c.GetOverview(token)
  if overview.TodayMinutes <= 0 { t.Fatalf("expected positive today minutes") }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/integration -run TestMVPFlow -v`  
Expected: FAIL with missing helpers or endpoint gaps

- [ ] **Step 3: Write minimal implementation**

```yaml
# docker-compose.yml
services:
  mysql:
    image: mysql:8.0.27
    environment:
      MYSQL_ROOT_PASSWORD: 010511
      MYSQL_DATABASE: learning_growth
    ports:
      - "3306:3306"
```

```markdown
# README.md
## Start
1. `docker compose up -d mysql`
2. `cd backend && cp .env.example .env && go run ./cmd/server`
3. `cd frontend && npm install && npm run dev`

## Verify
1. `cd backend && go test ./...`
2. `cd frontend && npm test -- --run`
3. Manual flow: register -> task -> study -> checkin -> dashboard
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./internal/integration -run TestMVPFlow -v`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd .worktrees/codex-learning-growth-mvp
git add docker-compose.yml backend/internal/integration/mvp_flow_test.go README.md
git commit -m "test: add mvp flow integration and runtime docs"
```

## Global Verification Checklist

- [ ] Backend tests pass: `cd .worktrees/codex-learning-growth-mvp/backend && go test ./...`
- [ ] Frontend tests pass: `cd .worktrees/codex-learning-growth-mvp/frontend && npm test -- --run`
- [ ] Frontend build passes: `cd .worktrees/codex-learning-growth-mvp/frontend && npm run build`
- [ ] Health API passes: `curl http://localhost:8080/api/health`
- [ ] Manual acceptance passes: create task, create study session, check in, view dashboard metrics



