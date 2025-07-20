package main

import (
	"bytes"
	"encoding/json"
	"incident-management/database"
	"incident-management/handlers"
	"incident-management/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestServer() *gin.Engine {
	// Initialize test database
	err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// Create Gin router
	r := gin.Default()

	// Create handler instance
	incidentHandler := handlers.NewIncidentHandler()

	// API routes
	api := r.Group("/api/v1")
	{
		incidents := api.Group("/incidents")
		{
			incidents.POST("/", incidentHandler.CreateIncident)
			incidents.GET("/", incidentHandler.GetAllIncidents)
		}
	}

	return r
}

func TestIntegration_CreateAndGetIncident(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test server
	router := setupTestServer()

	// Test data
	incident := model.Incident{
		Title:       "Integration Test Incident",
		Description: "This is a test incident for integration testing",
		Status:      "open",
		Priority:    "high",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(incident)
	if err != nil {
		t.Fatalf("Failed to marshal incident: %v", err)
	}

	// Create POST request
	postReq, err := http.NewRequest("POST", "/api/v1/incidents/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	postReq.Header.Set("Content-Type", "application/json")

	// Execute POST request
	postRecorder := httptest.NewRecorder()
	router.ServeHTTP(postRecorder, postReq)

	// Check POST response
	if postRecorder.Code != http.StatusCreated {
		t.Errorf("Expected POST status %d, got %d", http.StatusCreated, postRecorder.Code)
	}

	// Parse POST response
	var createdIncident model.Incident
	err = json.Unmarshal(postRecorder.Body.Bytes(), &createdIncident)
	if err != nil {
		t.Fatalf("Failed to unmarshal POST response: %v", err)
	}

	// Verify created incident
	if createdIncident.Title != "Integration Test Incident" {
		t.Errorf("Expected title 'Integration Test Incident', got '%s'", createdIncident.Title)
	}

	if createdIncident.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if createdIncident.AISeverity == "" {
		t.Error("Expected AI severity to be set")
	}

	if createdIncident.AICategory == "" {
		t.Error("Expected AI category to be set")
	}

	// Create GET request
	getReq, err := http.NewRequest("GET", "/api/v1/incidents/", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	// Execute GET request
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)

	// Check GET response
	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected GET status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	// Parse GET response
	var incidents []model.Incident
	err = json.Unmarshal(getRecorder.Body.Bytes(), &incidents)
	if err != nil {
		t.Fatalf("Failed to unmarshal GET response: %v", err)
	}

	// Verify the incident is in the list
	found := false
	for _, inc := range incidents {
		if inc.ID == createdIncident.ID {
			found = true
			// Verify all fields are preserved
			if inc.Title != createdIncident.Title {
				t.Errorf("Expected title to be preserved, got '%s'", inc.Title)
			}
			if inc.Description != createdIncident.Description {
				t.Errorf("Expected description to be preserved, got '%s'", inc.Description)
			}
			if inc.AISeverity != createdIncident.AISeverity {
				t.Errorf("Expected AI severity to be preserved, got '%s'", inc.AISeverity)
			}
			if inc.AICategory != createdIncident.AICategory {
				t.Errorf("Expected AI category to be preserved, got '%s'", inc.AICategory)
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find created incident in GET response")
	}
}

func TestIntegration_MultipleIncidents(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test server
	router := setupTestServer()

	// Test data for multiple incidents
	testIncidents := []model.Incident{
		{
			Title:       "First Integration Incident",
			Description: "First incident for multiple test",
			Status:      "open",
			Priority:    "low",
		},
		{
			Title:       "Second Integration Incident",
			Description: "Second incident for multiple test",
			Status:      "open",
			Priority:    "high",
		},
		{
			Title:       "Third Integration Incident",
			Description: "Third incident for multiple test",
			Status:      "open",
			Priority:    "medium",
		},
	}

	createdIDs := make([]string, 0)

	// Create all incidents
	for _, incident := range testIncidents {
		jsonData, err := json.Marshal(incident)
		if err != nil {
			t.Fatalf("Failed to marshal incident: %v", err)
		}

		req, err := http.NewRequest("POST", "/api/v1/incidents/", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
		}

		var createdIncident model.Incident
		err = json.Unmarshal(recorder.Body.Bytes(), &createdIncident)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		createdIDs = append(createdIDs, createdIncident.ID)
	}

	// Get all incidents
	getReq, err := http.NewRequest("GET", "/api/v1/incidents/", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getReq)

	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected GET status %d, got %d", http.StatusOK, getRecorder.Code)
	}

	var incidents []model.Incident
	err = json.Unmarshal(getRecorder.Body.Bytes(), &incidents)
	if err != nil {
		t.Fatalf("Failed to unmarshal GET response: %v", err)
	}

	// Verify all created incidents are present
	if len(incidents) < len(testIncidents) {
		t.Errorf("Expected at least %d incidents, got %d", len(testIncidents), len(incidents))
	}

	// Check that all our created incidents are in the response
	for _, createdID := range createdIDs {
		found := false
		for _, incident := range incidents {
			if incident.ID == createdID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find incident with ID %s in GET response", createdID)
		}
	}
}

func TestIntegration_AIAnalysis(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test server
	router := setupTestServer()

	// Test incident that should trigger specific AI analysis
	incident := model.Incident{
		Title:       "Critical Security Breach",
		Description: "Unauthorized access detected on production server. Multiple failed login attempts from suspicious IP addresses. Potential data breach.",
		Status:      "open",
		Priority:    "high",
	}

	jsonData, err := json.Marshal(incident)
	if err != nil {
		t.Fatalf("Failed to marshal incident: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/incidents/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var createdIncident model.Incident
	err = json.Unmarshal(recorder.Body.Bytes(), &createdIncident)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify AI analysis was performed
	if createdIncident.AISeverity == "" {
		t.Error("Expected AI severity to be set")
	}

	if createdIncident.AICategory == "" {
		t.Error("Expected AI category to be set")
	}

	// Verify AI fields are valid
	validSeverities := []string{"low", "medium", "high"}
	validCategories := []string{"network", "software", "hardware", "security"}

	severityValid := false
	for _, severity := range validSeverities {
		if createdIncident.AISeverity == severity {
			severityValid = true
			break
		}
	}
	if !severityValid {
		t.Errorf("Expected valid AI severity, got '%s'", createdIncident.AISeverity)
	}

	categoryValid := false
	for _, category := range validCategories {
		if createdIncident.AICategory == category {
			categoryValid = true
			break
		}
	}
	if !categoryValid {
		t.Errorf("Expected valid AI category, got '%s'", createdIncident.AICategory)
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
