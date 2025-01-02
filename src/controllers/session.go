package controllers

import (
	"let-me-in/database"
	"let-me-in/models"
	"let-me-in/modules/auth"
	"let-me-in/terminal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartSession starts a new terminal session.
func StartSession(c *gin.Context) {
	var input struct {
		UserID      uint   `json:"user_id" binding:"required"`
		ContainerID string `json:"container_id" binding:"required"`
		IPAddress   string `json:"ip_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user auth.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	session := models.Session{
		UserID:       input.UserID,
		ContainerID:  input.ContainerID,
		IPAddress:    input.IPAddress,
		Status:       "active",
		LastActivity: time.Now(),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Session started successfully", "session_id": session.ID})
}

// ListSessions lists all active sessions.
func ListSessions(c *gin.Context) {
	var sessions []models.Session
	if err := database.DB.Preload("User").Where("status = ?", "active").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust for production)
	},
}

func TerminalWebSocket(c *gin.Context) {
	token := c.Query("token") // Get token from query parameter

	_, err := auth.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Start terminal session for the authenticated user
	terminal.StartTerminalSession(conn, "bash")
}
