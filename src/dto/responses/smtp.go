package responses

import (
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"time"
)

// Smtp struct for the SMTP response.
type Smtp struct {
	ID                   uint      `json:"id"`
	AppMailID            uint      `json:"appMailId"`
	App                  string    `json:"app"`
	Mail                 string    `json:"mail"`
	Username             string    `json:"username"`
	Host                 string    `json:"host"`
	Port                 int       `json:"port"`
	DkimPrivateKey       *string   `json:"dkimPrivateKey"`
	DkimDomain           *string   `json:"dkimDomain"`
	DkimCanonicalization *string   `json:"dkimCanonicalization"`
	Primary              bool      `json:"primary"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// SetSmtp sets the SMTP response.
func (response *Smtp) SetSmtp(smtp *models.Smtp) {
	response.ID = smtp.ID
	response.AppMailID = smtp.AppMailID
	response.App = smtp.AppMail.AppName
	response.Mail = smtp.AppMail.MailName
	response.Username = smtp.Username
	response.Host = smtp.Host
	response.Port = smtp.Port
	response.DkimPrivateKey = smtp.DkimPrivateKey
	response.DkimDomain = smtp.DkimDomain
	response.DkimCanonicalization = smtp.DkimCanonicalizationName
	response.CreatedAt = smtp.CreatedAt
	response.UpdatedAt = smtp.UpdatedAt

	smtpType := enums.SMTP
	if smtp.AppMail.PrimaryType.Valid && smtp.AppMail.PrimaryType.String == *smtpType.ToString() {
		response.Primary = true
	}
}
