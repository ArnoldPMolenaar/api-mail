package requests

// CreateGmail struct for creating a new Gmail.
type CreateGmail struct {
	App      string `json:"app" validate:"required"`
	Mail     string `json:"mail" validate:"required,email"`
	ClientID string `json:"clientId" validate:"required"`
	Secret   string `json:"secret" validate:"required"`
	User     string `json:"user" validate:"required"`
	Primary  bool   `json:"primary"`
}
