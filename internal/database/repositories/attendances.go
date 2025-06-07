package repositories

import (
	"fmt"
	"time"

	"github.com/cyberbrain-dev/na-meste-api/internal/database/entities"
	"github.com/cyberbrain-dev/na-meste-api/internal/models"
	"gorm.io/gorm"
)

// Represents a repository of attendances
type Attendances struct {
	db *gorm.DB
}

// Creates new attendances repo of the db passed
func NewAttendances(db *gorm.DB) *Attendances {
	return &Attendances{db: db}
}

// Creates a new attendance record
func (r *Attendances) Create(a *models.Attendance) error {
	entity := entities.Attendance{
		ID:        a.ID,
		UserID:    a.UserID,
		CollegeID: a.CollegeID,
		Date:      a.Date,
	}

	result := r.db.Create(&entity)

	return result.Error
}

// Returns attendance by its ID
func (r *Attendances) Get(id uint) (*models.Attendance, error) {
	var entities []entities.Attendance

	result := r.db.Where("id = ?", id).Find(&entities)

	if result.Error != nil {
		return nil, fmt.Errorf("cannot get the college: %w", result.Error)
	}

	// if nothing has been found
	if len(entities) == 0 {
		return nil, nil
	}

	attendance := models.Attendance{
		ID:        entities[0].ID,
		UserID:    entities[0].UserID,
		CollegeID: entities[0].CollegeID,
		Date:      entities[0].Date,
	}

	return &attendance, nil
}

// Returns the attendances of the user and date span
func (r *Attendances) GetByStudentAndDatespan(id uint, start time.Time, end time.Time) ([]*models.Attendance, error) {
	var entities []entities.Attendance

	result := r.db.Where("(user_id = ?) AND (date BETWEEN ? AND ?)", id, start, end).Find(&entities)

	if result.Error != nil {
		return nil, fmt.Errorf("cannot get the attendances: %w", result.Error)
	}

	if len(entities) == 0 {
		return nil, nil
	}

	var attmodels []*models.Attendance

	for _, entity := range entities {
		attmodels = append(
			attmodels,
			&models.Attendance{
				ID:        entity.ID,
				UserID:    entity.UserID,
				CollegeID: entity.CollegeID,
				Date:      entity.Date,
			},
		)
	}

	return attmodels, nil
}

// Deletes an attendance by an ID
func (r *Attendances) Delete(id uint) (uint, error) {
	result := r.db.Where("id = ?", id).Delete(&entities.Attendance{})
	if result.Error != nil {
		return 0, fmt.Errorf(`not able to delete the attendance №%d: %w`, id, result.Error)
	}

	return id, nil
}
