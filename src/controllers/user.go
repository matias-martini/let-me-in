package controllers

import (
	"net/http"
	"let-me-in/database"
	"let-me-in/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration.
func RegisterUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

