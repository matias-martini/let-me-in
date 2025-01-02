package auth

import (
	"bytes"
	"encoding/json"
	"let-me-in/database"
	"let-me-in/modules/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginUserHandler(t *testing.T) {
	// Set up the test environment
	database.InitTestDB()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	// Register a user first
	registerBody := map[string]string{
		"display_name": "testuser",
		"email":        "testuser2@example.com",
		"password":     "testpassword",
	}
	jsonRegisterBody, _ := json.Marshal(registerBody)

	reqRegister, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonRegisterBody))
	reqRegister.Header.Set("Content-Type", "application/json")

	wRegister := httptest.NewRecorder()
	router.ServeHTTP(wRegister, reqRegister)

	assert.Equal(t, http.StatusOK, wRegister.Code)
	// Now test the login
	loginBody := map[string]string{
		"email":    "testuser2@example.com",
		"password": "testpassword",
	}
	jsonLoginBody, _ := json.Marshal(loginBody)

	reqLogin, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonLoginBody))
	reqLogin.Header.Set("Content-Type", "application/json")

	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)

	// Assert the login response
	assert.Equal(t, http.StatusOK, wLogin.Code)

	// Parse the login response body
	var loginResponse map[string]interface{}
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)

	// Check if the login response contains the expected fields
	assert.Contains(t, loginResponse, "access_token")
	assert.Contains(t, loginResponse, "refresh_token")
	database.ResetTestDB()
}
