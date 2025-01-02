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

func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRefreshTokenHandler(t *testing.T) {
	// Set up the test environment
	database.InitTestDB()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	// Register a user first
	registerBody := map[string]string{
		"display_name": "testuser",
		"email":        "testuser3@example.com",
		"password":     "testpassword",
	}
	wRegister := performRequest(router, "POST", "/auth/register", registerBody)
	assert.Equal(t, http.StatusOK, wRegister.Code)

	// Login to get the refresh token
	loginBody := map[string]string{
		"email":    "testuser3@example.com",
		"password": "testpassword",
	}
	wLogin := performRequest(router, "POST", "/auth/login", loginBody)
	assert.Equal(t, http.StatusOK, wLogin.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	refreshToken := loginResponse["refresh_token"].(string)

	// Now test the refresh token
	refreshBody := map[string]string{
		"refresh_token": refreshToken,
	}
	wRefresh := performRequest(router, "POST", "/auth/refresh", refreshBody)

	// Assert the refresh response
	assert.Equal(t, http.StatusOK, wRefresh.Code)

	// Parse the refresh response body
	var refreshResponse map[string]interface{}
	json.Unmarshal(wRefresh.Body.Bytes(), &refreshResponse)

	// Check if the refresh response contains the expected fields
	assert.Contains(t, refreshResponse, "access_token")
	assert.Contains(t, refreshResponse, "refresh_token")

	database.ResetTestDB()
}

func TestRefreshTokenCannotBeUsedMoreThanOnce(t *testing.T) {
	// Set up the test environment
	database.InitTestDB()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	// Register a user
	registerBody := map[string]string{
		"display_name": "testuser",
		"email":        "testuser4@example.com",
		"password":     "testpassword",
	}
	wRegister := performRequest(router, "POST", "/auth/register", registerBody)
	assert.Equal(t, http.StatusOK, wRegister.Code)

	// Login to get the refresh token
	loginBody := map[string]string{
		"email":    "testuser4@example.com",
		"password": "testpassword",
	}
	wLogin := performRequest(router, "POST", "/auth/login", loginBody)
	assert.Equal(t, http.StatusOK, wLogin.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	refreshToken := loginResponse["refresh_token"].(string)

	// Use the refresh token for the first time
	refreshBody := map[string]string{
		"refresh_token": refreshToken,
	}
	wRefresh := performRequest(router, "POST", "/auth/refresh", refreshBody)
	assert.Equal(t, http.StatusOK, wRefresh.Code)

	// Try to use the same refresh token again
	wRefresh2 := performRequest(router, "POST", "/auth/refresh", refreshBody)

	// Assert that the second attempt fails
	assert.Equal(t, http.StatusUnauthorized, wRefresh2.Code)

	var errorResponse map[string]string
	json.Unmarshal(wRefresh2.Body.Bytes(), &errorResponse)
	assert.Equal(t, "Invalid refresh token", errorResponse["error"])

	database.ResetTestDB()
}
