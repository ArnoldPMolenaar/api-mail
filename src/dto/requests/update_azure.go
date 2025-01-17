package requests

import "time"

// UpdateAzure struct for updating a Azure record.
type UpdateAzure struct {
	ClientID  string    `json:"clientId" validate:"required"`
	TenantID  string    `json:"tenantId" validate:"required"`
	Secret    string    `json:"secret" validate:"required"`
	User      string    `json:"user" validate:"required"`
	Primary   bool      `json:"primary"`
	UpdatedAt time.Time `json:"updatedAt"`
}
