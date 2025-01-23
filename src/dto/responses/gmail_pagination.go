package responses

import (
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"time"
)

// GmailPagination struct for the Gmail GET all response.
type GmailPagination struct {
	ID        uint      `json:"id"`
	ClientID  string    `json:"clientId"`
	User      string    `json:"user"`
	Primary   bool      `json:"primary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetGmailPagination sets the response.
func (response *GmailPagination) SetGmailPagination(gmail *models.Gmail) {
	response.ID = gmail.ID
	response.ClientID = gmail.ClientID
	response.User = gmail.User
	response.CreatedAt = gmail.CreatedAt
	response.UpdatedAt = gmail.UpdatedAt

	gmailType := enums.Gmail
	if gmail.AppMail.PrimaryType.Valid && gmail.AppMail.PrimaryType.String == *gmailType.ToString() {
		response.Primary = true
	}
}
