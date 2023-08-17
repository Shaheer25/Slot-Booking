package Testcontrollers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockUserCollection struct {
	mock.Mock
}

func (m *mockUserCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {

	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)

}

type mockTicketCollection struct {
	mock.Mock
}

func (m *mockTicketCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {

	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

type MockTicketCollection struct {
	mock.Mock
}

func (m *MockTicketCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

type MockReservationCollection struct {
	mock.Mock
}

func (m *MockReservationCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func TestGetUser(t *testing.T) {

	router := gin.Default()
	router.GET("/users/:user_id", GetUser())

	
	req, err := http.NewRequest(http.MethodGet, "/users/user_id", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	} else {
		t.Log("Success")
	}
}

func TestGenerateTickets(t *testing.T) {

	router := gin.Default()
	router.POST("/generate-tickets", GenerateTickets())

	payload := map[string]interface{}{
		"date":       "2023-07-21T00:00:00Z",
		"start_time": "2023-07-21T10:00:00Z",
		"end_time":   "2023-07-21T12:00:00Z",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "/generate-tickets", bytes.NewBuffer(jsonPayload))
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

func TestDeleteTicket(t *testing.T) {

	router := gin.Default()
	router.DELETE("/delete-ticket", DeleteTicket())

	payload := map[string]interface{}{
		"ticket_ids": []int{1, 2, 3},
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodDelete, "/delete-ticket", bytes.NewBuffer(jsonPayload))
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

func TestGetAllReservations(t *testing.T) {

	router := gin.Default()
	router.GET("/reservations", GetAllReservations())

	req, err := http.NewRequest(http.MethodGet, "/reservations", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	} else {
		t.Log("Success")
	}
}

type MockGetUserHandler struct {
	UserID       string
	ReturnedUser interface{}
	Err          error
}

func (m *MockGetUserHandler) Handle(c *gin.Context) {

	if m.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": m.Err.Error()})
		return
	}

	user := map[string]interface{}{
		"User_ID":   m.UserID,
		"First_Name": "John",
		"Last_Name":  "Doe",
		"Email":      "john.doe@example.com",
	}
	c.JSON(http.StatusOK, user)
}
type MockGenerateTicketsHandler struct {
	Success bool
}

func (m *MockGenerateTicketsHandler) Handle(c *gin.Context) {

	if m.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Tickets generated successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tickets"})
	}
}


type MockDeleteTicketHandler struct {
	Success bool
}

func (m *MockDeleteTicketHandler) Handle(c *gin.Context) {

	if m.Success {
		c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
	}
}


type MockGetAllReservationsHandler struct {
	Reservations []interface{}
	Err          error
}

func (m *MockGetAllReservationsHandler) Handle(c *gin.Context) {

	if m.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": m.Err.Error()})
		return
	}


	reservations := m.Reservations
	c.JSON(http.StatusOK, reservations)
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GenerateTickets() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func DeleteTicket() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetAllReservations() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}