package auth

import (
	"bytes"
	"encoding/json"
	"let-me-in/database"
	"net/http"
	"net/http/httptest"
	"testing"

	"let-me-in/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserHandlerSuccess(t *testing.T) {
	// Set up the test environment
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	// Create a test request body
	requestBody := map[string]string{
		"display_name": "testuser",
		"email":        "testuser@example.com",
		"password":     "testpassword",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	database.ResetTestDB()
}

func TestRegisterUserHandlerShortPassword(t *testing.T) {
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	requestBody := map[string]string{
		"display_name": "testuser",
		"email":        "testuser@example.com",
		"password":     "ab", // Only 2 characters
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "Password must be at least")

	database.ResetTestDB()
}

func TestRegisterUserHandlerInvalidEmail(t *testing.T) {
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	requestBody := map[string]string{
		"display_name": "tito_pruebas",
		"email":        "not_an_email",
		"password":     "testpassword",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "Invalid email format")
}

func TestRegisterUserHandlerEmptyDisplayName(t *testing.T) {
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	requestBody := map[string]string{
		"display_name": "",
		"email":        "valid@example.com",
		"password":     "validpassword",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "Display name cannot be empty")

	database.ResetTestDB()
}

func TestRegisterUserHandlerDuplicateEmail(t *testing.T) {
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	// First registration
	requestBody := map[string]string{
		"display_name": "test_user",
		"email":        "duplicate@example.com",
		"password":     "validpassword",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Second registration with the same email
	req, _ = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["error"], "Email already exists")

	database.ResetTestDB()
}

func TestRegisterUserHandlerMissingData(t *testing.T) {
	database.InitTestDB()
	router := gin.Default()
	auth.RegisterAuthRoutes(router.Group("/auth"))

	testCases := []struct {
		name         string
		requestBody  map[string]string
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Missing Display Name",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "validpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Display name cannot be empty",
		},
		{
			name: "Missing Email",
			requestBody: map[string]string{
				"display_name": "Test User",
				"password":     "validpassword",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Invalid email format",
		},
		{
			name: "Missing Password",
			requestBody: map[string]string{
				"display_name": "Test User",
				"email":        "test@example.com",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "Password must be at least 8 characters long",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Contains(t, response["error"], tc.expectedErr)
		})
	}

	database.ResetTestDB()
}
