package responses

import (
	"api-mail/main/src/models"
	"time"
)

// Smtp struct for the SMTP response.
type Smtp struct {
	App                  string    `json:"app"`
	Mail                 string    `json:"mail"`
	Username             string    `json:"username"`
	Host                 string    `json:"host"`
	Port                 int       `json:"port"`
	DkimPrivateKey       string    `json:"dkimPrivateKey"`
	DkimDomain           string    `json:"dkimDomain"`
	DkimCanonicalization string    `json:"dkimCanonicalization"`
	Primary              bool      `json:"primary"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// SetSmtp sets the SMTP response.
func (response *Smtp) SetSmtp(smtp *models.Smtp) {
	response.App = smtp.AppName
	response.Mail = smtp.MailName
	response.Username = smtp.Username
	response.Host = smtp.Host
	response.Port = smtp.Port
	response.DkimPrivateKey = smtp.DkimPrivateKey
	response.DkimDomain = smtp.DkimDomain
	response.DkimCanonicalization = smtp.DkimCanonicalizationName
	response.CreatedAt = smtp.CreatedAt
	response.UpdatedAt = smtp.UpdatedAt

	if appMail := smtp.GetAppMail(); appMail != nil {
		response.Primary = appMail.Primary
	}
}
