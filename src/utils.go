package main

import (
    "let-me-in/database"

    "let-me-in/modules/auth"
    "let-me-in/models"
)


func MigrateModels() {
    database.ConnectDatabase()

    database.DB.AutoMigrate(&auth.User{})
    database.DB.AutoMigrate(&models.Session{})
}
