package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(_ *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
