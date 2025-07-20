package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
}

type AIAnalysisResult struct {
	Severity string `json:"severity"`
	Category string `json:"category"`
}

// NewAIService creates a new AI service instance
func NewAIService() *AIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// For development, you can set a default key or handle this differently
		fmt.Println("Warning: OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)
	return &AIService{
		client: client,
	}
}

// AnalyzeIncident analyzes an incident description to determine severity and category
func (s *AIService) AnalyzeIncident(title, description string) (*AIAnalysisResult, error) {
	// If no API key is set, return default values
	if os.Getenv("OPENAI_API_KEY") == "" {
		return &AIAnalysisResult{
			Severity: "medium",
			Category: "software",
		}, nil
	}

	prompt := fmt.Sprintf(`
Analyze the following incident and determine:
1. Severity: Choose from "low", "medium", or "high"
2. Category: Choose from "network", "software", "hardware", or "security"

Consider these guidelines:
- Severity: Based on potential impact, urgency, and scope
- Category: Based on the type of issue described

Incident Title: %s
Incident Description: %s

Respond with a JSON object in this exact format:
{
  "severity": "low|medium|high",
  "category": "network|software|hardware|security"
}
`, title, description)

	ctx := context.Background()
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // Low temperature for more consistent results
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	content = strings.TrimSpace(content)

	// Try to parse the JSON response
	var result AIAnalysisResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// If JSON parsing fails, try to extract values using string manipulation
		result = s.extractValuesFromText(content)
	}

	// Validate the results
	if !s.isValidSeverity(result.Severity) {
		result.Severity = "medium" // Default fallback
	}
	if !s.isValidCategory(result.Category) {
		result.Category = "software" // Default fallback
	}

	return &result, nil
}

// extractValuesFromText extracts severity and category from text if JSON parsing fails
func (s *AIService) extractValuesFromText(text string) AIAnalysisResult {
	result := AIAnalysisResult{
		Severity: "medium",
		Category: "software",
	}

	text = strings.ToLower(text)

	// Extract severity
	if strings.Contains(text, `"severity"`) || strings.Contains(text, "severity") {
		if strings.Contains(text, `"low"`) || strings.Contains(text, "low") {
			result.Severity = "low"
		} else if strings.Contains(text, `"high"`) || strings.Contains(text, "high") {
			result.Severity = "high"
		}
	}

	// Extract category
	if strings.Contains(text, `"category"`) || strings.Contains(text, "category") {
		if strings.Contains(text, `"network"`) || strings.Contains(text, "network") {
			result.Category = "network"
		} else if strings.Contains(text, `"hardware"`) || strings.Contains(text, "hardware") {
			result.Category = "hardware"
		} else if strings.Contains(text, `"security"`) || strings.Contains(text, "security") {
			result.Category = "security"
		}
	} else {
		// If no category keyword found, try to detect from content
		if strings.Contains(text, "hardware") {
			result.Category = "hardware"
		} else if strings.Contains(text, "network") {
			result.Category = "network"
		} else if strings.Contains(text, "security") {
			result.Category = "security"
		}
	}

	return result
}

// isValidSeverity checks if the severity is valid
func (s *AIService) isValidSeverity(severity string) bool {
	validSeverities := []string{"low", "medium", "high"}
	for _, valid := range validSeverities {
		if strings.ToLower(severity) == valid {
			return true
		}
	}
	return false
}

// isValidCategory checks if the category is valid
func (s *AIService) isValidCategory(category string) bool {
	validCategories := []string{"network", "software", "hardware", "security"}
	for _, valid := range validCategories {
		if strings.ToLower(category) == valid {
			return true
		}
	}
	return false
}
