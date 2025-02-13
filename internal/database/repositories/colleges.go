package repositories

import (
	"fmt"

	"github.com/cyberbrain-dev/na-meste-api/internal/database/entities"
	"github.com/cyberbrain-dev/na-meste-api/internal/models"
	"gorm.io/gorm"
)

// Represents a repository of colleges
type Colleges struct {
	db *gorm.DB
}

// Creates new college repo of the db passed
func NewColleges(db *gorm.DB) *Colleges {
	return &Colleges{db: db}
}

// Adds a college to the db
func (r *Colleges) Create(c *models.College) error {
	entity := entities.College{
		ID:   c.ID,
		Name: c.Name,
	}

	result := r.db.Create(&entity)

	return result.Error
}

// Returns college by its name
func (r *Colleges) GetByName(name string) (*models.College, error) {
	var entity entities.College

	result := r.db.Where("name = ?", name).First(&entity)

	if result.Error != nil {
		return nil, fmt.Errorf("cannot get the college: %w", result.Error)
	}

	college := models.College{
		ID:   entity.ID,
		Name: entity.Name,
	}

	return &college, nil
}
