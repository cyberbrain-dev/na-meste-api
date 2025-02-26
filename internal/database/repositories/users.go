package repositories

import (
	"fmt"

	"github.com/cyberbrain-dev/na-meste-api/internal/database/entities"
	"github.com/cyberbrain-dev/na-meste-api/internal/models"
	"gorm.io/gorm"
)

// Represents a Postgres users repository
type Users struct {
	db *gorm.DB
}

// Creates a new users repository on a databse passed
func NewUsers(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (r *Users) Create(u *models.User) error {
	entity := entities.User{
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		CollegeID:    u.CollegeID,
	}

	result := r.db.Create(&entity)

	return result.Error
}

func (r *Users) Get(email string) (*models.User, error) {
	var entities []entities.User

	result := r.db.Where("email = ?", email).Find(&entities)
	if result.Error != nil {
		return nil, fmt.Errorf("cannot get the user: %w", result.Error)
	}

	// If user has not been found
	if len(entities) == 0 {
		return nil, nil
	}

	e := entities[0]

	user := models.User{
		ID:           e.ID,
		Username:     e.Username,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		Role:         e.Role,
		CollegeID:    e.CollegeID,
	}

	return &user, nil
}

func (r *Users) Update(id uint, username *string, email *string) (uint, error) {
	if username == nil && email == nil {
		return id, nil
	}

	if username == nil && email != nil {
		result := r.db.Model(&entities.User{}).Where("id = ?", id).Update("email", *email)
		if result.Error != nil {
			return 0, fmt.Errorf("cannot update the email: %w", result.Error)
		}
	} else if username != nil && email == nil {
		result := r.db.Model(&entities.User{}).Where("id = ?", id).Update("username", *username)
		if result.Error != nil {
			return 0, fmt.Errorf("cannot update the username: %w", result.Error)
		}
	} else if username != nil && email != nil {
		result := r.db.Model(&entities.User{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"username": *username,
				"email":    *email,
			})

		if result.Error != nil {
			return 0, fmt.Errorf("cannot update user: %w", result.Error)
		}
	}

	return id, nil
}

func (r *Users) Delete(id uint) (uint, error) {
	result := r.db.Where("id = ?", id).Delete(&entities.User{})
	if result.Error != nil {
		return 0, fmt.Errorf(`not able to delete the user: %w`, result.Error)
	}

	return id, nil
}
