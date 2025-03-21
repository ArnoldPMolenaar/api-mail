package requests

// CreateApp request DTO to create a App.
type CreateApp struct {
	Name string `json:"name" validate:"required"`
}
