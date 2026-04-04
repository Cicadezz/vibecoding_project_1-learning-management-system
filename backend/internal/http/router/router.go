package router

import (
	"net/http"

	"learning-growth-platform/internal/auth"
	"learning-growth-platform/internal/checkin"
	"learning-growth-platform/internal/http/middleware"
	"learning-growth-platform/internal/stats"
	"learning-growth-platform/internal/study"
	"learning-growth-platform/internal/subjects"
	"learning-growth-platform/internal/tasks"
	"learning-growth-platform/internal/timer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if db != nil {
		authRepo := auth.NewRepository(db)
		authSvc := auth.NewService(authRepo, "")
		taskSvc := tasks.NewService(tasks.NewRepository(db))
		authSvc.SetTaskService(taskSvc)
		authHandler := auth.NewHandler(authSvc)
		authMW := middleware.NewAuthMiddleware(authSvc)
		subjectRepo := subjects.NewRepository(db)
		subjectHandler := subjects.NewHandler(subjects.NewService(subjectRepo))
		taskHandler := tasks.NewHandler(taskSvc)
		studyRepo := study.NewRepository(db)
		studyHandler := study.NewHandler(study.NewService(studyRepo, subjectRepo))
		checkinHandler := checkin.NewHandler(checkin.NewService(checkin.NewRepository(db)))
		statsHandler := stats.NewHandler(stats.NewService(stats.NewRepository(db)))
		timerHandler := timer.NewHandler(timer.NewService(timer.NewRepository(db), subjectRepo))

		authGroup := r.Group("/api/auth")
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		protected := authGroup.Group("")
		protected.Use(authMW.RequireAuth())
		protected.GET("/me", authHandler.Me)
		protected.POST("/change-password", authHandler.ChangePassword)

		apiGroup := r.Group("/api")
		apiGroup.Use(authMW.RequireAuth())
		apiGroup.GET("/subjects", subjectHandler.List)
		apiGroup.POST("/subjects", subjectHandler.Create)
		apiGroup.PUT("/subjects/:id", subjectHandler.Update)
		apiGroup.DELETE("/subjects/:id", subjectHandler.Delete)
		apiGroup.GET("/tasks", taskHandler.List)
		apiGroup.POST("/tasks", taskHandler.Create)
		apiGroup.PUT("/tasks/:id", taskHandler.Update)
		apiGroup.DELETE("/tasks/:id", taskHandler.Delete)
		apiGroup.POST("/study/sessions", studyHandler.Create)
		apiGroup.POST("/checkin/today", checkinHandler.Today)
		apiGroup.GET("/checkin/streak", checkinHandler.Streak)
		apiGroup.GET("/stats/overview", statsHandler.Overview)
		apiGroup.GET("/stats/weekly-trend", statsHandler.WeeklyTrend)
		apiGroup.GET("/stats/subject-distribution", statsHandler.SubjectDistribution)
		apiGroup.POST("/timer/start", timerHandler.Start)
		apiGroup.POST("/timer/stop", timerHandler.Stop)
	}

	return r
}
