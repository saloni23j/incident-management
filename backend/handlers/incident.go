package handlers

import (
	"incident-management/model"
	"incident-management/services"
	"incident-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncidentHandler struct {
	service *services.IncidentService
}

// NewIncidentHandler creates a new incident handler
func NewIncidentHandler() *IncidentHandler {
	return &IncidentHandler{
		service: services.NewIncidentService(),
	}
}

// CreateIncident handles POST /incidents
func (h *IncidentHandler) CreateIncident(c *gin.Context) {
	var incident model.Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Validate the incident
	validationErrors := utils.ValidateAndGetErrors(&incident)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Set defaults if not provided
	if incident.Status == "" {
		incident.Status = "open"
	}
	if incident.Priority == "" {
		incident.Priority = "medium"
	}

	createdIncident, err := h.service.CreateIncident(incident)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create incident",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdIncident)
}

// GetAllIncidents handles GET /incidents
func (h *IncidentHandler) GetAllIncidents(c *gin.Context) {
	incidents, err := h.service.GetAllIncidents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve incidents",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

// HealthCheck handles GET /health
func (h *IncidentHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Incident Management API is running",
		"version": "1.0.0",
	})
}
