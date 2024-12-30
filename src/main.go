package main

import (
	"let-me-in/controllers"
	"let-me-in/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	router := gin.Default()

	// User routes
	router.POST("/users/register", controllers.RegisterUser)
	router.POST("/auth/login", controllers.Login)

	// Session routes
	router.POST("/sessions/start", controllers.StartSession)
	router.GET("/sessions", controllers.ListSessions)

	// WebSocket route for terminal access
	router.GET("/ws/terminal", controllers.TerminalWebSocket)

	// Serve frontend
	router.Static("/static", "./static")

	router.Run(":8080")
}
