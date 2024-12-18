package database

import (
	"github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := database.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}
