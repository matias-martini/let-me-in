package auth 

import (
	"net/http"
	"regexp"
    "time"
    "net/mail"

	"let-me-in/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUserHandler handles user registration with salt and pepper for password security
func RegisterUserHandler(c *gin.Context) {
	var input struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	db := database.DB

	// Validate request payload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate required fields with specific messages
	if input.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display name cannot be empty"})
		return
	}

	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if input.Password == "" || len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Validate password length
	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Check if email is valid using mail.ParseAddress
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Check if email already exists
	var existingUser UserCredentials
	if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email existence"})
		return
	}

	// Validate display name
	if input.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display name cannot be empty"})
		return
	}

	// Generate salt
	salt, err := generateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}

	// Hash password using salt and pepper
	hashedPassword, err := hashPassword(input.Password, salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create User entry
	user := User{DisplayName: input.DisplayName}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Create UserCredentials entry
	userCredentials := UserCredentials{
		Email:    input.Email,
		Password: hashedPassword,
		Salt:     salt,
		UserID:   user.ID,
	}

	if err := db.Create(&userCredentials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}


func LoginUserHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	db := database.DB

	// Validate input payload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user credentials by email
	var userCredentials UserCredentials
	if err := db.Where("email = ?", input.Email).First(&userCredentials).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user credentials: " + err.Error()})
		}
		return
	}

	// Verify the password with salt and pepper
	if !verifyPassword(input.Password, userCredentials.Salt, userCredentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate Access Token (JWT)
	accessToken, err := GenerateJWT(userCredentials.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token: " + err.Error()})
		return
	}

	// Generate Refresh Token
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token: " + err.Error()})
		return
	}

	expiresAt := time.Now().Add(1 * 24 * time.Hour).Unix() // 1 day expiration

	// Store Refresh Token
	refreshTokenModel := RefreshToken{
		Token:     refreshToken,
		UserID:    userCredentials.UserID,
		ExpiresAt: expiresAt,
		Active:    true,
	}

	if err := db.Create(&refreshTokenModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshTokenHandler handles refreshing access and ID tokens using a valid refresh token
func RefreshTokenHandler(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	db := database.DB

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the refresh token in the database
	var refreshTokenModel RefreshToken
	if err := db.Where("token = ?", input.RefreshToken).First(&refreshTokenModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query refresh token"})
		}
		return
	}

	// Check if the refresh token is active and not expired
	if !refreshTokenModel.Active || time.Now().Unix() > refreshTokenModel.ExpiresAt {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is expired or inactive"})
		return
	}

	// Fetch user credentials
	var userCredentials UserCredentials
	if err := db.Where("user_id = ?", refreshTokenModel.UserID).First(&userCredentials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
		return
	}

	// Generate new Access Token
	accessToken, err := GenerateJWT(userCredentials.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate new ID Token (same as Access Token in this case, but can be expanded if needed)
	idToken, err := GenerateJWT(userCredentials.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID token"})
		return
	}

	// (Optional) Rotate refresh token
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rotate refresh token"})
		return
	}

	// Deactivate old refresh token and store new one
	refreshTokenModel.Token = newRefreshToken
	refreshTokenModel.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix() // Extend expiry
	if err := db.Save(&refreshTokenModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	// Return the new tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"id_token":      idToken,
		"refresh_token": newRefreshToken,
	})
}

