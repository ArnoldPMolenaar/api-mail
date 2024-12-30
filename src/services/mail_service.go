package services

import (
	"api-mail/main/src/database"
	"api-mail/main/src/models"
	// mail "github.com/xhit/go-simple-mail/v2"
)

// TODO: Put the email settings into valkey.

// IsMailAvailable method to check if a mail is available.
func IsMailAvailable(mail string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Mail{}, "name = ?", mail); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsMailTypeAvailable method to check if a mail-type is available.
func IsMailTypeAvailable(mailType string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.MailType{}, "name = ?", mailType); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetAppMail finds a mail by app, mail and mailType.
// If mail or mailType is empty, it will return the primary mail.
func GetAppMail(app, mail, mailType string) (models.AppMail, error) {
	var appMail models.AppMail
	query := database.Pg.Where("app_name = ?", app)

	if mail != "" {
		query = query.Where("mail_name = ?", mail)
	}

	if mailType != "" {
		query = query.Where("mail_type = ?", mailType)
	}

	if result := query.Order("primary DESC").First(&appMail); result.Error != nil {
		return appMail, result.Error
	}

	return appMail, nil
}

func SendSMTPMail(fromName, fromMail, to, subject, body, mimeType string, ccs []string, bccs []string) error {
	// server := mail.NewSMTPClient()

	// SMTP server configuration.
	//server.Host =

	return nil
}
