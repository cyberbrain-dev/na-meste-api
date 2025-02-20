package abstractions

import "github.com/cyberbrain-dev/na-meste-api/internal/models"

type AttendancesRepo interface {
	// Adds a new record to the db
	Create(a *models.Attendance) error

	// Returns an attendance by an ID
	Get(id uint) (*models.Attendance, error)

	// Deletes an attendance by an ID
	Delete(id uint) (uint, error)
}
