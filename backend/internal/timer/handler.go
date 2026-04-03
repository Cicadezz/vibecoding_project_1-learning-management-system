package timer

import (
	"encoding/json"
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

type startRequest struct {
	SubjectID uint64          `json:"subject_id"`
	Note      *string         `json:"note"`
	Ext       json.RawMessage `json:"ext"`
}

type stopRequest struct {
	Note *string `json:"note"`
}

func (h *Handler) Start(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req startRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	state, err := h.svc.Start(userID, req.SubjectID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTimerInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrTimerAlreadyRunning):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"timer_state": state})
}

func (h *Handler) Stop(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req stopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	session, err := h.svc.Stop(userID, req.Note)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTimerInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrTimerNotRunning):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"study_session": session})
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}
