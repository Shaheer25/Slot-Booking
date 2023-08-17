package Testcontrollers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func mockSignup() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func mockLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func TestMain(m *testing.M) {

	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func Signup(t *testing.T) {

	router := gin.Default()

	router.POST("/signup", mockSignup())

	payload := []byte(`{
		"First_name": "John",
		"Last_name": "Doe",
		"Email": "john.doe@example.com",
		"Password": "password123",
		"Phone": "1234567890",
	}`)
	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	} else {

		t.Log("Success")

	}
}

func Login(t *testing.T) {

	router := gin.Default()

	router.POST("/login", mockLogin())

	payload := []byte(`{
		"Email": "john.doe@example.com",
		"Password": "password123",
	}`)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	} else {
		t.Log("Success")
	}

}
