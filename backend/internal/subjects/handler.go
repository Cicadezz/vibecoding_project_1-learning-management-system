package subjects

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type subjectCreateRequest struct {
	Name  string          `json:"name"`
	Color *string         `json:"color"`
	Ext   json.RawMessage `json:"ext"`
}

type subjectUpdateRequest struct {
	Name  *string         `json:"name"`
	Color *string         `json:"color"`
	Ext   json.RawMessage `json:"ext"`
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subjects, err := h.svc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subjects": subjects})
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req subjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	subject, err := h.svc.Create(CreateSubjectInput{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Ext:    []byte(req.Ext),
	})
	if err != nil {
		if errors.Is(err, ErrInvalidSubjectInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"subject": subject})
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subjectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || subjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var req subjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	subject, err := h.svc.Update(subjectID, userID, UpdateSubjectInput{
		Name:  req.Name,
		Color: req.Color,
		Ext:   []byte(req.Ext),
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidSubjectInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case isSubjectNotFound(err):
			c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"subject": subject})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subjectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || subjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.Delete(subjectID, userID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidSubjectInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case isSubjectNotFound(err):
			c.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
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
