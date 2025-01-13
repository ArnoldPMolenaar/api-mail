package responses

import (
	"api-mail/main/src/models"
	"time"
)

// GmailCallback struct for the Oauth response.
type GmailCallback struct {
	ID           uint      `json:"id"`
	AppMailID    uint      `json:"appMailId"`
	ClientID     string    `json:"clientId"`
	Secret       string    `json:"secret"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    string    `json:"tokenType"`
	Expiry       time.Time `json:"expiry"`
	ExpiresIn    int64     `json:"expiresIn"`
	User         string    `json:"user"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// SetGmailCallback sets the response.
func (response *GmailCallback) SetGmailCallback(gmail *models.Gmail) {
	response.ID = gmail.ID
	response.AppMailID = gmail.AppMailID
	response.ClientID = gmail.ClientID
	response.Secret = gmail.Secret
	response.AccessToken = *gmail.AccessToken
	response.RefreshToken = *gmail.RefreshToken
	response.TokenType = *gmail.TokenType
	response.Expiry = *gmail.Expiry
	response.ExpiresIn = *gmail.ExpiresIn
	response.User = gmail.User
	response.CreatedAt = gmail.CreatedAt
	response.UpdatedAt = gmail.UpdatedAt
}
