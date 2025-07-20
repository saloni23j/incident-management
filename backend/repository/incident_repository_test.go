package repository

import (
	"incident-management/database"
	"incident-management/model"
	"os"
	"testing"
)

func TestNewIncidentRepository(t *testing.T) {
	// Initialize database first
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	repo := NewIncidentRepository()
	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
	if repo.db == nil {
		t.Fatal("Expected database connection, got nil")
	}
}

func TestCreateAndGetAll(t *testing.T) {
	// Initialize test database
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	repo := NewIncidentRepository()

	// Test Create
	incident := &model.Incident{
		Title:       "Test Repository Incident",
		Description: "This is a test incident for repository",
		Status:      "open",
		Priority:    "medium",
		AISeverity:  "medium",
		AICategory:  "software",
	}

	err = repo.Create(incident)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	// Check that ID was generated
	if incident.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	// Test GetAll
	incidents, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all incidents: %v", err)
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
		if inc.Title == "Test Repository Incident" {
			found = true
			// Check that all fields are preserved
			if inc.Description != "This is a test incident for repository" {
				t.Errorf("Expected description to be preserved, got '%s'", inc.Description)
			}
			if inc.Status != "open" {
				t.Errorf("Expected status 'open', got '%s'", inc.Status)
			}
			if inc.Priority != "medium" {
				t.Errorf("Expected priority 'medium', got '%s'", inc.Priority)
			}
			if inc.AISeverity != "medium" {
				t.Errorf("Expected AI severity 'medium', got '%s'", inc.AISeverity)
			}
			if inc.AICategory != "software" {
				t.Errorf("Expected AI category 'software', got '%s'", inc.AICategory)
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find test incident in GetAll results")
	}
}

func TestCreateMultipleIncidents(t *testing.T) {
	// Initialize test database
	err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	repo := NewIncidentRepository()

	// Create multiple incidents
	incidents := []*model.Incident{
		{
			Title:       "First Test Incident",
			Description: "First incident description",
			Status:      "open",
			Priority:    "low",
			AISeverity:  "low",
			AICategory:  "network",
		},
		{
			Title:       "Second Test Incident",
			Description: "Second incident description",
			Status:      "open",
			Priority:    "high",
			AISeverity:  "high",
			AICategory:  "security",
		},
		{
			Title:       "Third Test Incident",
			Description: "Third incident description",
			Status:      "open",
			Priority:    "medium",
			AISeverity:  "medium",
			AICategory:  "hardware",
		},
	}

	// Create all incidents
	for _, incident := range incidents {
		err = repo.Create(incident)
		if err != nil {
			t.Fatalf("Failed to create incident '%s': %v", incident.Title, err)
		}
	}

	// Get all incidents
	allIncidents, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all incidents: %v", err)
	}

	// Should have at least the number of incidents we created
	if len(allIncidents) < len(incidents) {
		t.Errorf("Expected at least %d incidents, got %d", len(incidents), len(allIncidents))
	}

	// Check that all our test incidents are present
	expectedTitles := map[string]bool{
		"First Test Incident":  false,
		"Second Test Incident": false,
		"Third Test Incident":  false,
	}

	for _, incident := range allIncidents {
		if _, exists := expectedTitles[incident.Title]; exists {
			expectedTitles[incident.Title] = true
		}
	}

	for title, found := range expectedTitles {
		if !found {
			t.Errorf("Expected to find incident with title '%s'", title)
		}
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
