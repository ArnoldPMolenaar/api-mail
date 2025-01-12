package responses

import (
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"time"
)

// Gmail struct for the Gmail response.
type Gmail struct {
	ID        uint      `json:"id"`
	AppMailID uint      `json:"appMailId"`
	App       string    `json:"app"`
	Mail      string    `json:"mail"`
	ClientID  string    `json:"clientId"`
	Secret    string    `json:"secret"`
	User      string    `json:"user"`
	Primary   bool      `json:"primary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetGmail sets the Gmail response.
func (response *Gmail) SetGmail(gmail *models.Gmail) {
	response.ID = gmail.ID
	response.AppMailID = gmail.AppMailID
	response.App = gmail.AppMail.AppName
	response.Mail = gmail.AppMail.MailName
	response.ClientID = gmail.ClientID
	response.Secret = gmail.Secret
	response.User = gmail.User
	response.CreatedAt = gmail.CreatedAt
	response.UpdatedAt = gmail.UpdatedAt

	gmailType := enums.Gmail
	if gmail.AppMail.PrimaryType.Valid && gmail.AppMail.PrimaryType.String == *gmailType.ToString() {
		response.Primary = true
	}
}
