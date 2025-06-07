// Contains utils for connecting to databases
package database

import (
	"fmt"

	"github.com/cyberbrain-dev/na-meste-api/internal/config"
	"github.com/cyberbrain-dev/na-meste-api/internal/database/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Tries to connect to the db and returns a gorm db
func ConnectPostgres(cfg config.PostgresConnection) (*gorm.DB, error) {

	// declaring a connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
	)

	// opening the db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database (check the config)")
	}

	return db, nil
}

// Closes the connection to the postgreSQL database
func DisconnectPostgres(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("cannot get the sql.DB from gorm db")
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("unable to close the db")
	}

	return nil
}

// Migrates the entities to Postgres database
func MigrateEntities(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.College{},
		&entities.User{},
		&entities.Attendance{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate the entities: %w", err)
	}

	db.Exec(`
		ALTER TABLE users DROP CONSTRAINT fk_colleges_users;

		ALTER TABLE users 
    	ADD CONSTRAINT fk_colleges_users 
    	FOREIGN KEY (college_id) REFERENCES colleges(id) 
		ON UPDATE CASCADE
		ON DELETE SET NULL;
	`)

	db.Exec(`
		ALTER TABLE attendances DROP CONSTRAINT fk_colleges_attendances;

		ALTER TABLE attendances
    	ADD CONSTRAINT fk_colleges_attendances
    	FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE;
	`)

	db.Exec(`
		ALTER TABLE attendances DROP CONSTRAINT fk_users_attendances;

		ALTER TABLE attendances
    	ADD CONSTRAINT fk_users_attendances
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
	`)

	return nil
}
