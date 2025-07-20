package services

import (
	"incident-management/model"
	"incident-management/repository"
	"log"
)

type IncidentService struct {
	repo *repository.IncidentRepository
	ai   *AIService
}

// NewIncidentService creates a new incident service
func NewIncidentService() *IncidentService {
	return &IncidentService{
		repo: repository.NewIncidentRepository(),
		ai:   NewAIService(),
	}
}

// CreateIncident creates a new incident with AI analysis
func (s *IncidentService) CreateIncident(incident model.Incident) (*model.Incident, error) {
	// Set default values if not provided
	if incident.Status == "" {
		incident.Status = "open"
	}
	if incident.Priority == "" {
		incident.Priority = "medium"
	}

	// Use AI to analyze the incident and determine severity and category
	aiResult, err := s.ai.AnalyzeIncident(incident.Title, incident.Description)
	if err != nil {
		// If AI analysis fails, use default value
		log.Println("AI analysis failed, using default values", err)
		incident.AISeverity = "medium"
		incident.AICategory = "software"
	} else {
		incident.AISeverity = aiResult.Severity
		incident.AICategory = aiResult.Category
	}

	err = s.repo.Create(&incident)
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

// GetAllIncidents retrieves all incidents
func (s *IncidentService) GetAllIncidents() ([]model.Incident, error) {
	return s.repo.GetAll()
}
