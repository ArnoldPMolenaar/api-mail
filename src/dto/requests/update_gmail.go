package requests

import "time"

// UpdateGmail struct for updating a Gmail record.
type UpdateGmail struct {
	ClientID  string    `json:"clientId" validate:"required"`
	Secret    string    `json:"secret" validate:"required"`
	User      string    `json:"user" validate:"required"`
	Primary   bool      `json:"primary"`
	UpdatedAt time.Time `json:"updatedAt"`
}
