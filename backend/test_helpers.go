package main

import (
	"incident-management/database"
	"incident-management/model"
	"os"
)

// TestHelper provides helper functions for testing
type TestHelper struct{}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// SetupTestDatabase initializes a clean test database
func (h *TestHelper) SetupTestDatabase() error {
	// Remove existing test database
	os.Remove("incidents.db")

	// Initialize fresh database
	return database.InitDB()
}

// CleanupTestDatabase removes the test database
func (h *TestHelper) CleanupTestDatabase() {
	os.Remove("incidents.db")
}

// CreateTestIncident creates a test incident with default values
func (h *TestHelper) CreateTestIncident(title, description string) model.Incident {
	return model.Incident{
		Title:       title,
		Description: description,
		Status:      "open",
		Priority:    "medium",
	}
}

// CreateTestIncidentWithValues creates a test incident with specific values
func (h *TestHelper) CreateTestIncidentWithValues(title, description, status, priority string) model.Incident {
	return model.Incident{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
	}
}

// ValidateIncidentFields validates that an incident has all required fields
func (h *TestHelper) ValidateIncidentFields(incident model.Incident) []string {
	var errors []string

	if incident.ID == "" {
		errors = append(errors, "ID is empty")
	}

	if incident.Title == "" {
		errors = append(errors, "Title is empty")
	}

	if incident.Description == "" {
		errors = append(errors, "Description is empty")
	}

	if incident.Status == "" {
		errors = append(errors, "Status is empty")
	}

	if incident.Priority == "" {
		errors = append(errors, "Priority is empty")
	}

	if incident.AISeverity == "" {
		errors = append(errors, "AI Severity is empty")
	}

	if incident.AICategory == "" {
		errors = append(errors, "AI Category is empty")
	}

	return errors
}

// ValidateAIFields validates that AI-determined fields are valid
func (h *TestHelper) ValidateAIFields(incident model.Incident) []string {
	var errors []string

	validSeverities := []string{"low", "medium", "high"}
	validCategories := []string{"network", "software", "hardware", "security"}

	severityValid := false
	for _, severity := range validSeverities {
		if incident.AISeverity == severity {
			severityValid = true
			break
		}
	}
	if !severityValid {
		errors = append(errors, "Invalid AI severity: "+incident.AISeverity)
	}

	categoryValid := false
	for _, category := range validCategories {
		if incident.AICategory == category {
			categoryValid = true
			break
		}
	}
	if !categoryValid {
		errors = append(errors, "Invalid AI category: "+incident.AICategory)
	}

	return errors
}
