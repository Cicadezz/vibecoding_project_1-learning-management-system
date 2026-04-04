package stats

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Overview(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	overview, err := h.svc.Overview(userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStatsInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"overview": overview})
}

func (h *Handler) WeeklyTrend(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trend, err := h.svc.WeeklyTrend(userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStatsInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"weekly_trend": trend})
}

func (h *Handler) SubjectDistribution(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	distribution, err := h.svc.SubjectDistribution(userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStatsInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"subject_distribution": distribution})
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}
