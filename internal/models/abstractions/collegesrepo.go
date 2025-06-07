package abstractions

import "github.com/cyberbrain-dev/na-meste-api/internal/models"

// Represents an abstract colleges repository
type CollegesRepo interface {
	// Adds a college to the db
	Create(c *models.College) error

	// Returns college by its name
	Get(name string) (*models.College, error)

	// Deletes the college and returns its ID
	Delete(id uint) (uint, error)
}
