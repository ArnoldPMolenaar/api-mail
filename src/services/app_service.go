package services

import (
	"api-mail/main/src/database"
	"api-mail/main/src/models"
)

// IsAppAvailable method to check if an app is available.
func IsAppAvailable(app string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.App{}, "name = ?", app); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// CreateApp method to create an app.
func CreateApp(name string) (*models.App, error) {
	app := &models.App{Name: name}

	if err := database.Pg.FirstOrCreate(&models.App{}, app).Error; err != nil {
		return nil, err
	}

	return app, nil
}
