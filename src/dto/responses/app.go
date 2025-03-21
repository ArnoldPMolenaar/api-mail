package responses

import "api-mail/main/src/models"

// App struct to map the app.
type App struct {
	Name string `json:"name"`
}

// SetApp method to set the app.
func (a *App) SetApp(app *models.App) {
	a.Name = app.Name
}
