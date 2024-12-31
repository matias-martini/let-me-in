
package cmd

import (
    "fmt"
    "let-me-in/database"
    "github.com/spf13/cobra"

    "let-me-in/modules/auth"
    "let-me-in/models"
)



// dbCmd is the parent command: "let-me-in db"
var dbCmd = &cobra.Command{
    Use:   "db",
    Short: "Database utilities",
    Long:  `Perform various database operations like migrate, seed, drop, etc.`,
    // We won’t specify a Run function here,
    // so running "let-me-in db" alone will just show help.
}

// dbMigrateCmd represents "let-me-in db migrate"
var dbMigrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Run database migrations",
    Long:  `Executes all pending database migrations.`,
    Run: func(cmd *cobra.Command, args []string) {
        migrateModels()
    },
}

func init() {
    rootCmd.AddCommand(dbCmd)
    dbCmd.AddCommand(dbMigrateCmd)
}

func migrateModels() {
    fmt.Println("Running migrations...")
    database.ConnectDatabase()

    if err := database.DB.AutoMigrate(&auth.User{}); err != nil {
        fmt.Printf("Error migrating User model: %v\n", err)
        return
    }

    if err := database.DB.AutoMigrate(&models.Session{}); err != nil {
        fmt.Printf("Error migrating Session model: %v\n", err)
        return
    }

    fmt.Println("Migrations completed successfully!")
}
