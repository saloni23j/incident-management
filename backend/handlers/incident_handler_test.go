package handlers

import (
	"bytes"
	"encoding/json"
	"incident-management/database"
	"incident-management/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewIncidentHandler(t *testing.T) {
	handler := NewIncidentHandler()
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	if handler.service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

func TestCreateIncident(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	handler := NewIncidentHandler()

	// Create test request
	incident := model.Incident{
		Title:       "Test Handler Incident",
		Description: "This is a test incident from handler",
		Status:      "open",
		Priority:    "medium",
	}

	jsonData, err := json.Marshal(incident)
	if err != nil {
		t.Fatalf("Failed to marshal incident: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/incidents", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler.CreateIncident(c)

	// Check response status
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response body
	var response model.Incident
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check response fields
	if response.Title != "Test Handler Incident" {
		t.Errorf("Expected title 'Test Handler Incident', got '%s'", response.Title)
	}

	if response.Description != "This is a test incident from handler" {
		t.Errorf("Expected description 'This is a test incident from handler', got '%s'", response.Description)
	}

	if response.Status != "open" {
		t.Errorf("Expected status 'open', got '%s'", response.Status)
	}

	if response.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", response.Priority)
	}

	// Check that AI fields were set
	if response.AISeverity == "" {
		t.Error("Expected AI severity to be set")
	}

	if response.AICategory == "" {
		t.Error("Expected AI category to be set")
	}

	// Check that ID was generated
	if response.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestCreateIncident_InvalidJSON(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	handler := NewIncidentHandler()

	// Create test request with invalid JSON
	invalidJSON := `{"title": "Test", "description": "Test", invalid json}`

	req, err := http.NewRequest("POST", "/api/v1/incidents", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler.CreateIncident(c)

	// Check response status - should be bad request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateIncident_ValidationError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	handler := NewIncidentHandler()

	// Create test request with invalid data (empty title)
	incident := model.Incident{
		Title:       "", // Empty title should fail validation
		Description: "This is a test incident",
		Status:      "open",
		Priority:    "medium",
	}

	jsonData, err := json.Marshal(incident)
	if err != nil {
		t.Fatalf("Failed to marshal incident: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/incidents", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler.CreateIncident(c)

	// Check response status - should be bad request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check error message
	if response["error"] != "Validation failed" {
		t.Errorf("Expected error 'Validation failed', got '%v'", response["error"])
	}

	// Check that details contain validation errors
	details, exists := response["details"].(map[string]interface{})
	if !exists {
		t.Fatal("Expected details to contain validation errors")
	}

	if _, exists := details["title"]; !exists {
		t.Error("Expected title validation error in details")
	}
}

func TestGetAllIncidents(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	handler := NewIncidentHandler()

	// Create test request
	req, err := http.NewRequest("GET", "/api/v1/incidents", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler.GetAllIncidents(c)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var response []model.Incident
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Response should be an array (even if empty)
	if response == nil {
		t.Error("Expected incidents array, got nil")
	}
}

func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	handler := NewIncidentHandler()

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	handler.HealthCheck(c)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check response fields
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", response["status"])
	}

	if response["message"] != "Incident Management API is running" {
		t.Errorf("Expected message 'Incident Management API is running', got '%v'", response["message"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", response["version"])
	}
}

func TestMain(m *testing.M) {
	// Clean up test database before running tests
	os.Remove("incidents.db")

	// Run tests
	code := m.Run()

	// Clean up test database after running tests
	os.Remove("incidents.db")

	os.Exit(code)
}
