package responses

import (
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"time"
)

// AzurePagination struct for the Azure GET all response.
type AzurePagination struct {
	ID        uint      `json:"id"`
	ClientID  string    `json:"clientId"`
	TenantID  string    `json:"tenantId"`
	User      string    `json:"user"`
	Primary   bool      `json:"primary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetAzurePagination sets the response.
func (response *AzurePagination) SetAzurePagination(azure *models.Azure) {
	response.ID = azure.ID
	response.ClientID = azure.ClientID
	response.TenantID = azure.TenantID
	response.User = azure.User
	response.CreatedAt = azure.CreatedAt
	response.UpdatedAt = azure.UpdatedAt

	azureType := enums.Azure
	if azure.AppMail.PrimaryType.Valid && azure.AppMail.PrimaryType.String == *azureType.ToString() {
		response.Primary = true
	}
}
