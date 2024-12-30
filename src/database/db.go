package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"let-me-in/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("terminal.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models
	err = DB.AutoMigrate(&models.User{}, &models.Session{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connected and migrated")
}

