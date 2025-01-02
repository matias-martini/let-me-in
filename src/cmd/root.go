package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for your application.
// It doesn’t do anything by itself, but it can hold persistent flags,
// global config, etc.
var rootCmd = &cobra.Command{
	Use:   "let-me-in",
	Short: "A secure and user-friendly terminal access management service.",
	Long: `Let Me In is a robust service designed to provide secure and persistent access 
to terminal sessions through a web interface. 

It supports multiple terminal sessions, session persistence after WebSocket disconnections, 
and configurable session expiration timeouts. Built with scalability and security in mind, 
Let Me In simplifies terminal session management while ensuring a seamless experience 
for both administrators and end users.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
