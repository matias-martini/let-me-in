package terminal

import (
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

// StartTerminalSession starts a PTY terminal session and links it with WebSocket
func StartTerminalSession(conn *websocket.Conn, command string) {
	// Start a shell or any command (e.g., Django or Rails shell)
	cmd := exec.Command("bash") // You can replace with "python manage.py shell" or "rails console"

	// Start PTY session
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Failed to start PTY:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to start PTY session."))
		return
	}
	defer ptmx.Close()

	// Goroutine to read from PTY and send to WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				log.Println("PTY read error:", err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, buf[:n])
			if err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		}
	}()

	// Read from WebSocket and send to PTY
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}
		_, err = ptmx.Write(msg)
		if err != nil {
			log.Println("PTY write error:", err)
			return
		}
	}
}

