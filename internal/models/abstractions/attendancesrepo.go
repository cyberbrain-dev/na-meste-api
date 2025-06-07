package abstractions

import (
	"time"

	"github.com/cyberbrain-dev/na-meste-api/internal/models"
)

type AttendancesRepo interface {
	// Adds a new record to the db
	Create(a *models.Attendance) error

	// Returns an attendance by an ID
	Get(id uint) (*models.Attendance, error)

	// Returns the attendances of the user and date span
	GetByStudentAndDatespan(id uint, from time.Time, to time.Time) ([]*models.Attendance, error)

	// Deletes an attendance by an ID
	Delete(id uint) (uint, error)
}
