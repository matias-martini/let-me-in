package cmd

import (
    "fmt"
    "let-me-in/controllers"
    "let-me-in/database"
    "let-me-in/modules/auth"

    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
)

var (
    port int
)

// serveCmd represents "let-me-in serve", command to run the web server
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the web server",
    Long:  `Starts the web server on a specified port (default: 8080).`,
    Run: func(cmd *cobra.Command, args []string) {
        startServer()
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)

    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}

func startServer() {
    database.Init()

    router := gin.Default()

    auth.RegisterAuthRoutes(router.Group("/auth"))

    // Session routes
    router.POST("/sessions/start", controllers.StartSession)
    router.GET("/sessions", controllers.ListSessions)

    // WebSocket route for terminal access
    router.GET("/ws/terminal", controllers.TerminalWebSocket)

    // Serve frontend
    router.Static("/static", "./static")

    // Run the server on the chosen port
    fmt.Printf("Starting server on port %d...\n", port)
    router.Run(fmt.Sprintf(":%d", port))
}
