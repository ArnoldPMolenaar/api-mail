package responses

import (
	"api-mail/main/src/models"
	"time"
)

// AzureCallback struct for the Oauth response.
type AzureCallback struct {
	ID           uint       `json:"id"`
	AppMailID    uint       `json:"appMailId"`
	ClientID     string     `json:"clientId"`
	TenantID     string     `json:"tenantId"`
	Secret       string     `json:"secret"`
	AccessToken  *string    `json:"accessToken"`
	RefreshToken *string    `json:"refreshToken"`
	TokenType    *string    `json:"tokenType"`
	Expiry       *time.Time `json:"expiry"`
	ExpiresIn    *int64     `json:"expiresIn"`
	User         string     `json:"user"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// SetAzureCallback sets the response.
func (response *AzureCallback) SetAzureCallback(azure *models.Azure) {
	response.ID = azure.ID
	response.AppMailID = azure.AppMailID
	response.ClientID = azure.ClientID
	response.TenantID = azure.TenantID
	response.Secret = azure.Secret

	if azure.AccessToken.Valid {
		response.AccessToken = &azure.AccessToken.String
	}
	if azure.RefreshToken.Valid {
		response.RefreshToken = &azure.RefreshToken.String
	}
	if azure.TokenType.Valid {
		response.TokenType = &azure.TokenType.String
	}
	if azure.Expiry.Valid {
		response.Expiry = &azure.Expiry.Time
	}
	if azure.ExpiresIn.Valid {
		response.ExpiresIn = &azure.ExpiresIn.Int64
	}

	response.User = azure.User
	response.CreatedAt = azure.CreatedAt
	response.UpdatedAt = azure.UpdatedAt
}
