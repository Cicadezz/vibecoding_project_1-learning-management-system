package tasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type taskCreateRequest struct {
	Title       string          `json:"title"`
	Priority    string          `json:"priority"`
	DueDate     *time.Time      `json:"due_date"`
	PlanDate    time.Time       `json:"plan_date"`
	Status      string          `json:"status"`
	CompletedAt *time.Time      `json:"completed_at"`
	Ext         json.RawMessage `json:"ext"`
}

type taskUpdateRequest struct {
	Title       *string         `json:"title"`
	Priority    *string         `json:"priority"`
	DueDate     *time.Time      `json:"due_date"`
	PlanDate    *time.Time      `json:"plan_date"`
	Status      *string         `json:"status"`
	CompletedAt *time.Time      `json:"completed_at"`
	Ext         json.RawMessage `json:"ext"`
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	date := time.Now()
	if dateParam := c.Query("date"); dateParam != "" {
		parsed, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		date = parsed
	}

	tasks, err := h.svc.ListByDate(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req taskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	task, err := h.svc.Create(CreateTaskInput{
		UserID:      userID,
		Title:       req.Title,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		PlanDate:    req.PlanDate,
		Status:      req.Status,
		CompletedAt: req.CompletedAt,
		Ext:         []byte(req.Ext),
	})
	if err != nil {
		if errors.Is(err, ErrInvalidTaskInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var req taskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	task, err := h.svc.Update(taskID, userID, UpdateTaskInput{
		Title:       req.Title,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		PlanDate:    req.PlanDate,
		Status:      req.Status,
		CompletedAt: req.CompletedAt,
		Ext:         []byte(req.Ext),
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidTaskInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case isTaskNotFound(err):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || taskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.Delete(taskID, userID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidTaskInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case isTaskNotFound(err):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}
