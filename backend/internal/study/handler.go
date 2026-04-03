package study

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createSessionRequest struct {
	SubjectID uint64          `json:"subject_id"`
	StartAt   time.Time       `json:"start_at"`
	EndAt     time.Time       `json:"end_at"`
	Note      *string         `json:"note"`
	Ext       json.RawMessage `json:"ext"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	session, err := h.svc.CreateManual(CreateManualSessionInput{
		UserID:    userID,
		SubjectID: req.SubjectID,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		Note:      req.Note,
		Ext:       []byte(req.Ext),
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStudyInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrSubjectNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
