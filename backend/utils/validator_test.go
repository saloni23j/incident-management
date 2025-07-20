package utils

import (
	"strings"
	"testing"
)

func TestInitValidator(t *testing.T) {
	InitValidator()
	if validate == nil {
		t.Fatal("Expected validator to be initialized, got nil")
	}
}

func TestValidateAndGetErrors_ValidStruct(t *testing.T) {
	InitValidator()

	type TestStruct struct {
		Name string `validate:"required,min=1,max=50"`
		Age  int    `validate:"min=0,max=150"`
	}

	validStruct := TestStruct{
		Name: "John Doe",
		Age:  30,
	}

	errors := ValidateAndGetErrors(&validStruct)
	if errors != nil {
		t.Errorf("Expected no errors for valid struct, got: %v", errors)
	}
}

func TestValidateAndGetErrors_InvalidStruct(t *testing.T) {
	InitValidator()

	type TestStruct struct {
		Name string `validate:"required,min=1,max=50"`
		Age  int    `validate:"min=0,max=150"`
	}

	invalidStruct := TestStruct{
		Name: "",  // Empty name should fail required validation
		Age:  200, // Age 200 should fail max validation
	}

	errors := ValidateAndGetErrors(&invalidStruct)
	if errors == nil {
		t.Fatal("Expected errors for invalid struct, got nil")
	}

	// Check for specific error messages
	if _, exists := errors["name"]; !exists {
		t.Error("Expected name validation error")
	}

	if _, exists := errors["age"]; !exists {
		t.Error("Expected age validation error")
	}
}

func TestValidateAndGetErrors_OneOfValidation(t *testing.T) {
	InitValidator()

	type TestStruct struct {
		Status string `validate:"oneof=open closed pending"`
	}

	invalidStruct := TestStruct{
		Status: "invalid_status",
	}

	errors := ValidateAndGetErrors(&invalidStruct)
	if errors == nil {
		t.Fatal("Expected errors for invalid status, got nil")
	}

	if _, exists := errors["status"]; !exists {
		t.Error("Expected status validation error")
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello world  ", "hello world"},
		{"no spaces", "no spaces"},
		{"", ""},
		{"   ", ""},
		{"\t\n\r", ""},
	}

	for _, test := range tests {
		result := SanitizeString(test.input)
		if result != test.expected {
			t.Errorf("SanitizeString(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestValidateAndGetErrors_ErrorMessages(t *testing.T) {
	InitValidator()

	type TestStruct struct {
		Title       string `validate:"required,min=1,max=200"`
		Description string `validate:"required,min=1,max=1000"`
		Status      string `validate:"oneof=open in_progress resolved closed"`
	}

	invalidStruct := TestStruct{
		Title:       "",
		Description: "",
		Status:      "invalid",
	}

	errors := ValidateAndGetErrors(&invalidStruct)

	// Check error messages are descriptive
	if msg, exists := errors["title"]; exists {
		if !strings.Contains(msg, "required") {
			t.Errorf("Expected title error to mention 'required', got: %s", msg)
		}
	}

	if msg, exists := errors["description"]; exists {
		if !strings.Contains(msg, "required") {
			t.Errorf("Expected description error to mention 'required', got: %s", msg)
		}
	}

	if msg, exists := errors["status"]; exists {
		if !strings.Contains(msg, "one of:") {
			t.Errorf("Expected status error to mention 'one of:', got: %s", msg)
		}
	}
}
