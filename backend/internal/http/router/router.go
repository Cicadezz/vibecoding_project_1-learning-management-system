package router

import (
	"net/http"

	"learning-growth-platform/internal/auth"
	"learning-growth-platform/internal/http/middleware"
	"learning-growth-platform/internal/subjects"
	"learning-growth-platform/internal/tasks"

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
		subjectHandler := subjects.NewHandler(subjects.NewService(subjects.NewRepository(db)))
		taskHandler := tasks.NewHandler(taskSvc)

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
	}

	return r
}
