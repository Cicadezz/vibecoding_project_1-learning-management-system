package router

import (
	"net/http"

	"learning-growth-platform/internal/auth"
	"learning-growth-platform/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if db != nil {
		repo := auth.NewRepository(db)
		svc := auth.NewService(repo, "")
		h := auth.NewHandler(svc)
		authMW := middleware.NewAuthMiddleware(svc)

		authGroup := r.Group("/api/auth")
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)

		protected := authGroup.Group("")
		protected.Use(authMW.RequireAuth())
		protected.GET("/me", h.Me)
		protected.POST("/change-password", h.ChangePassword)
	}

	return r
}
