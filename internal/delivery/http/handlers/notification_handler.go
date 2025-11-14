package handlers

import (
	"net/http"
	"strconv"

	"github.com/dmehra2102/StudySync/internal/service"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationSvc *service.NotificationService
}

func NewNotificationHandler(notificationSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationSvc: notificationSvc,
	}
}

func (h *NotificationHandler) GetNotificaitons(c *gin.Context) {
	userID := c.GetUint("userID")
	unreadOnly, _ := strconv.ParseBool(c.Query("unread"))

	notifications, err := h.notificationSvc.GetUserNotifications(userID, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notificaitons": notifications,
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notificaiton ID"})
		return
	}

	if err := h.notificationSvc.MarkAsRead(uint(notificationID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.notificationSvc.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.notificationSvc.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
