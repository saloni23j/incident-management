package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the validator
func InitValidator() {
	validate = validator.New()
}

// ValidateAndGetErrors validates a struct and returns detailed error messages
func ValidateAndGetErrors(s interface{}) map[string]string {
	if validate == nil {
		InitValidator()
	}

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "min":
			errors[field] = field + " must be at least " + err.Param() + " characters"
		case "max":
			errors[field] = field + " must be at most " + err.Param() + " characters"
		case "oneof":
			errors[field] = field + " must be one of: " + err.Param()
		case "uuid4":
			errors[field] = field + " must be a valid UUID"
		default:
			errors[field] = field + " failed validation: " + err.Tag()
		}
	}

	return errors
}

// SanitizeString removes leading/trailing whitespace from a string
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeIncident sanitizes incident fields
func SanitizeIncident(incident interface{}) {
	// This would be implemented if we had reflection to modify struct fields
	// For now, we'll rely on the validation to catch issues
}
