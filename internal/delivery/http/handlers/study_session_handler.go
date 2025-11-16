package handlers

import (
	"net/http"
	"strconv"

	"github.com/dmehra2102/StudySync/internal/service"
	"github.com/gin-gonic/gin"
)

type StudySessionHandler struct {
	studySessionSvc *service.StudySessionService
}

func NewStudySessionHandler(studySessionSvc *service.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{
		studySessionSvc: studySessionSvc,
	}
}

func (h *StudySessionHandler) CreateSession(c *gin.Context) {
	userID := c.GetUint("userID")

	var req service.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.studySessionSvc.CreateSession(userID, req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Study session created successfully",
		"session": session,
	})
}

func (h *StudySessionHandler) GetSessions(c *gin.Context) {
	userID := c.GetUint("userID")

	sessions, err := h.studySessionSvc.GetUserSessions(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *StudySessionHandler) GetUpcomingSessions(c *gin.Context) {
	userID := c.GetUint("userID")

	sessions, err := h.studySessionSvc.GetUpcomingSessions(userID)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *StudySessionHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("userID")

	stats, err := h.studySessionSvc.GetUserStats(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *StudySessionHandler) GetSession(c *gin.Context) {
	userID := c.GetUint("userID")

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.studySessionSvc.GetSession(userID, uint(sessionID))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *StudySessionHandler) UpdateSession(c *gin.Context) {
	userID := c.GetUint("userID")

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req service.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	session, err := h.studySessionSvc.UpdateSession(userID, uint(sessionID), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"session": session,
	})
}

func (h *StudySessionHandler) DeleteSession(c *gin.Context) {
	userID := c.GetUint("userID")

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := h.studySessionSvc.DeleteSession(userID, uint(sessionID)); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}
