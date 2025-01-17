package requests

// CreateAzure struct for creating a new Azure.
type CreateAzure struct {
	App      string `json:"app" validate:"required"`
	Mail     string `json:"mail" validate:"required,email"`
	ClientID string `json:"clientId" validate:"required"`
	TenantID string `json:"tenantId" validate:"required"`
	Secret   string `json:"secret" validate:"required"`
	User     string `json:"user" validate:"required"`
	Primary  bool   `json:"primary"`
}
