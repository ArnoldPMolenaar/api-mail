package responses

import (
	"api-mail/main/src/models"
	"time"
)

// GmailCallback struct for the Oauth response.
type GmailCallback struct {
	ID           uint       `json:"id"`
	AppMailID    uint       `json:"appMailId"`
	ClientID     string     `json:"clientId"`
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

// SetGmailCallback sets the response.
func (response *GmailCallback) SetGmailCallback(gmail *models.Gmail) {
	response.ID = gmail.ID
	response.AppMailID = gmail.AppMailID
	response.ClientID = gmail.ClientID
	response.Secret = gmail.Secret

	if gmail.AccessToken.Valid {
		response.AccessToken = &gmail.AccessToken.String
	}
	if gmail.RefreshToken.Valid {
		response.RefreshToken = &gmail.RefreshToken.String
	}
	if gmail.TokenType.Valid {
		response.TokenType = &gmail.TokenType.String
	}
	if gmail.Expiry.Valid {
		response.Expiry = &gmail.Expiry.Time
	}
	if gmail.ExpiresIn.Valid {
		response.ExpiresIn = &gmail.ExpiresIn.Int64
	}

	response.User = gmail.User
	response.CreatedAt = gmail.CreatedAt
	response.UpdatedAt = gmail.UpdatedAt
}
