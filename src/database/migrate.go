package database

import (
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}
