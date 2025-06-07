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
func (r *Colleges) Get(name string) (*models.College, error) {
	var entities []entities.College

	result := r.db.Where("name = ?", name).Find(&entities)

	if result.Error != nil {
		return nil, fmt.Errorf("cannot get the college: %w", result.Error)
	}

	// if college has not been found
	if len(entities) == 0 {
		return nil, nil
	}

	college := models.College{
		ID:   entities[0].ID,
		Name: entities[0].Name,
	}

	return &college, nil
}

// Deletes the college and returns its ID
func (r *Colleges) Delete(id uint) (uint, error) {
	result := r.db.Where("id = ?", id).Delete(&entities.College{})
	if result.Error != nil {
		return 0, fmt.Errorf(`not able to delete the college №%d: %w`, id, result.Error)
	}

	return id, nil
}
