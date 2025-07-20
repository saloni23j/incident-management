package services

import (
	"incident-management/database"
	"incident-management/model"
	"os"
	"testing"
)

func TestNewIncidentService(t *testing.T) {
	service := NewIncidentService()
	if service == nil {
		t.Fatal("Expected incident service to be created, got nil")
	}
	if service.repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
	if service.ai == nil {
		t.Fatal("Expected AI service to be created, got nil")
	}
}

func TestCreateIncident_WithDefaults(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	service := NewIncidentService()

	incident := model.Incident{
		Title:       "Test Incident",
		Description: "This is a test incident",
	}

	createdIncident, err := service.CreateIncident(incident)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdIncident == nil {
		t.Fatal("Expected created incident, got nil")
	}

	// Check that defaults were set
	if createdIncident.Status != "open" {
		t.Errorf("Expected status 'open', got '%s'", createdIncident.Status)
	}

	if createdIncident.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", createdIncident.Priority)
	}

	// Check that AI fields were set (should be defaults when no API key)
	if createdIncident.AISeverity != "medium" {
		t.Errorf("Expected AI severity 'medium', got '%s'", createdIncident.AISeverity)
	}

	if createdIncident.AICategory != "software" {
		t.Errorf("Expected AI category 'software', got '%s'", createdIncident.AICategory)
	}

	// Check that ID was generated
	if createdIncident.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	// Check that timestamps were set
	if createdIncident.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set, got zero time")
	}

	if createdIncident.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set, got zero time")
	}
}

func TestCreateIncident_WithProvidedValues(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	service := NewIncidentService()

	incident := model.Incident{
		Title:       "Critical Server Issue",
		Description: "Production server is down",
		Status:      "open",
		Priority:    "high",
	}

	createdIncident, err := service.CreateIncident(incident)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdIncident == nil {
		t.Fatal("Expected created incident, got nil")
	}

	// Check that provided values were preserved
	if createdIncident.Title != "Critical Server Issue" {
		t.Errorf("Expected title 'Critical Server Issue', got '%s'", createdIncident.Title)
	}

	if createdIncident.Description != "Production server is down" {
		t.Errorf("Expected description 'Production server is down', got '%s'", createdIncident.Description)
	}

	if createdIncident.Status != "open" {
		t.Errorf("Expected status 'open', got '%s'", createdIncident.Status)
	}

	if createdIncident.Priority != "high" {
		t.Errorf("Expected priority 'high', got '%s'", createdIncident.Priority)
	}
}

func TestGetAllIncidents(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	service := NewIncidentService()

	// Create a test incident first
	incident := model.Incident{
		Title:       "Test Incident for GetAll",
		Description: "This incident should be retrieved by GetAll",
	}

	_, err = service.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create test incident: %v", err)
	}

	// Test GetAllIncidents
	incidents, err := service.GetAllIncidents()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if incidents == nil {
		t.Fatal("Expected incidents slice, got nil")
	}

	// Should have at least one incident
	if len(incidents) == 0 {
		t.Error("Expected at least one incident, got empty slice")
	}

	// Check that the test incident is in the results
	found := false
	for _, inc := range incidents {
		if inc.Title == "Test Incident for GetAll" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find test incident in GetAll results")
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
