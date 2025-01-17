package responses

import (
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"time"
)

// Azure struct for the Azure response.
type Azure struct {
	ID          uint      `json:"id"`
	AppMailID   uint      `json:"appMailId"`
	App         string    `json:"app"`
	Mail        string    `json:"mail"`
	ClientID    string    `json:"clientId"`
	TenantID    string    `json:"tenantId"`
	Secret      string    `json:"secret"`
	User        string    `json:"user"`
	Primary     bool      `json:"primary"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	AuthCodeURL *string   `json:"authCodeUrl"`
}

// SetAzure sets the Azure response.
func (response *Azure) SetAzure(azure *models.Azure, authCodeURL ...string) {
	response.ID = azure.ID
	response.AppMailID = azure.AppMailID
	response.App = azure.AppMail.AppName
	response.Mail = azure.AppMail.MailName
	response.ClientID = azure.ClientID
	response.TenantID = azure.TenantID
	response.Secret = azure.Secret
	response.User = azure.User
	response.CreatedAt = azure.CreatedAt
	response.UpdatedAt = azure.UpdatedAt

	azureType := enums.Azure
	if azure.AppMail.PrimaryType.Valid && azure.AppMail.PrimaryType.String == *azureType.ToString() {
		response.Primary = true
	}

	if len(authCodeURL) > 0 {
		response.AuthCodeURL = &authCodeURL[0]
	}
}
