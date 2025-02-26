package abstractions

import "github.com/cyberbrain-dev/na-meste-api/internal/models"

// Represents an abstract users repository
type UsersRepo interface {
	// Adds a new user record to the database
	Create(u *models.User) error

	// Returns a user with the ID passed if the one exists
	Get(id string) (*models.User, error)

	// Updates user with the ID passed
	Update(id uint, username *string, email *string) (uint, error)

	// Deletes user with the ID passed
	Delete(id uint) (uint, error)
}
