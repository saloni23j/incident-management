package repository

import (
	"incident-management/database"
	"incident-management/model"

	"gorm.io/gorm"
)

type IncidentRepository struct {
	db *gorm.DB
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository() *IncidentRepository {
	return &IncidentRepository{
		db: database.GetDB(),
	}
}

// Create creates a new incident
func (r *IncidentRepository) Create(incident *model.Incident) error {
	return r.db.Create(incident).Error
}

// GetAll retrieves all incidents
func (r *IncidentRepository) GetAll() ([]model.Incident, error) {
	var incidents []model.Incident
	err := r.db.Find(&incidents).Error
	return incidents, err
}
