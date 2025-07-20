package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Incident struct {
	ID          string `json:"id" gorm:"primaryKey;type:varchar(36)" validate:"omitempty,uuid4"`
	Title       string `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
	Description string `json:"description" gorm:"type:text" validate:"required,min=1,max=1000"`
	Status      string `json:"status" gorm:"default:'open'" validate:"omitempty,oneof=open in_progress resolved closed"`
	Priority    string `json:"priority" gorm:"default:'medium'" validate:"omitempty,oneof=low medium high critical"`
	// AI-determined fields
	AISeverity string    `json:"ai_severity" gorm:"default:'medium'" validate:"omitempty,oneof=low medium high"`
	AICategory string    `json:"ai_category" gorm:"default:'software'" validate:"omitempty,oneof=network software hardware security"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (incident *Incident) BeforeCreate(tx *gorm.DB) error {
	if incident.ID == "" {
		incident.ID = uuid.New().String()
	}
	return nil
}

// ValidateSeverity checks if the severity is valid
func (incident *Incident) ValidateSeverity() bool {
	validSeverities := []string{"low", "medium", "high"}
	for _, severity := range validSeverities {
		if incident.AISeverity == severity {
			return true
		}
	}
	return false
}

// ValidateCategory checks if the category is valid
func (incident *Incident) ValidateCategory() bool {
	validCategories := []string{"network", "software", "hardware", "security"}
	for _, category := range validCategories {
		if incident.AICategory == category {
			return true
		}
	}
	return false
}
