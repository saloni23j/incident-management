package services

import (
	"os"
	"testing"
)

func TestNewAIService(t *testing.T) {
	// Test with no API key
	aiService := NewAIService()
	if aiService == nil {
		t.Fatal("Expected AI service to be created, got nil")
	}
	if aiService.client == nil {
		t.Fatal("Expected OpenAI client to be created, got nil")
	}
}

func TestAnalyzeIncident_NoAPIKey(t *testing.T) {
	// Temporarily unset API key
	originalKey := os.Getenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", originalKey)

	aiService := NewAIService()

	result, err := aiService.AnalyzeIncident("Test Incident", "This is a test incident")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Should return default values
	if result.Severity != "medium" {
		t.Errorf("Expected severity 'medium', got '%s'", result.Severity)
	}

	if result.Category != "software" {
		t.Errorf("Expected category 'software', got '%s'", result.Category)
	}
}

func TestExtractValuesFromText(t *testing.T) {
	aiService := NewAIService()

	tests := []struct {
		name             string
		text             string
		expectedSeverity string
		expectedCategory string
	}{
		{
			name:             "Valid JSON response",
			text:             `{"severity": "high", "category": "security"}`,
			expectedSeverity: "high",
			expectedCategory: "security",
		},
		{
			name:             "Text with severity and category",
			text:             "The severity is low and the category is network",
			expectedSeverity: "low",
			expectedCategory: "network",
		},
		{
			name:             "Text with high severity and hardware category",
			text:             "This is a high severity hardware issue",
			expectedSeverity: "high",
			expectedCategory: "hardware",
		},
		{
			name:             "Invalid text - should use defaults",
			text:             "Random text without severity or category",
			expectedSeverity: "medium",
			expectedCategory: "software",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aiService.extractValuesFromText(tt.text)

			if result.Severity != tt.expectedSeverity {
				t.Errorf("Expected severity '%s', got '%s'", tt.expectedSeverity, result.Severity)
			}

			if result.Category != tt.expectedCategory {
				t.Errorf("Expected category '%s', got '%s'", tt.expectedCategory, result.Category)
			}
		})
	}
}

func TestIsValidSeverity(t *testing.T) {
	aiService := NewAIService()

	tests := []struct {
		severity string
		expected bool
	}{
		{"low", true},
		{"medium", true},
		{"high", true},
		{"LOW", true},
		{"Medium", true},
		{"HIGH", true},
		{"critical", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			result := aiService.isValidSeverity(tt.severity)
			if result != tt.expected {
				t.Errorf("Expected %v for severity '%s', got %v", tt.expected, tt.severity, result)
			}
		})
	}
}

func TestIsValidCategory(t *testing.T) {
	aiService := NewAIService()

	tests := []struct {
		category string
		expected bool
	}{
		{"network", true},
		{"software", true},
		{"hardware", true},
		{"security", true},
		{"NETWORK", true},
		{"Software", true},
		{"Hardware", true},
		{"Security", true},
		{"database", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := aiService.isValidCategory(tt.category)
			if result != tt.expected {
				t.Errorf("Expected %v for category '%s', got %v", tt.expected, tt.category, result)
			}
		})
	}
}
