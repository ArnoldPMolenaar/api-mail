package services

import (
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
)

// IsSmtpAvailable checks if the smtp exists.
func IsSmtpAvailable(app, mail string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Smtp{}, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetSmtp gets the smtp.
func GetSmtp(app, mail string) (*models.Smtp, error) {
	smtp := &models.Smtp{}
	if result := database.Pg.Find(smtp, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return nil, result.Error
	}

	return smtp, nil
}

// CreateSmtp creates a new smtp.
func CreateSmtp(req *requests.CreateSmtp) error {
	smtpType := enums.SMTP
	smtp := &models.Smtp{
		AppName:                  req.App,
		MailName:                 req.Mail,
		Username:                 req.Username,
		Password:                 req.Password,
		Host:                     req.Host,
		Port:                     req.Port,
		DkimPrivateKey:           req.DkimPrivateKey,
		DkimDomain:               req.DkimDomain,
		DkimCanonicalizationName: req.DkimCanonicalization,
	}

	if err := smtp.EncryptPassword(); err != nil {
		return err
	}

	if result := database.Pg.FirstOrCreate(&models.Mail{
		Name: req.Mail,
	}); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Create(smtp); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Create(&models.AppMail{
		AppName:  req.App,
		MailName: req.Mail,
		MailType: smtpType.ToString(),
		Primary:  req.Primary,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateSmtp updates a existing smtp.
func UpdateSmtp(oldSmtp *models.Smtp, req *requests.UpdateSmtp) error {
	smtpType := enums.SMTP
	smtp := &models.Smtp{
		AppName:                  oldSmtp.AppName,
		MailName:                 oldSmtp.MailName,
		Username:                 req.Username,
		Password:                 oldSmtp.Password,
		Host:                     req.Host,
		Port:                     req.Port,
		DkimPrivateKey:           req.DkimPrivateKey,
		DkimDomain:               req.DkimDomain,
		DkimCanonicalizationName: req.DkimCanonicalization,
		CreatedAt:                oldSmtp.CreatedAt,
	}

	if req.Password != "" {
		smtp.Password = req.Password
		if err := smtp.EncryptPassword(); err != nil {
			return err
		}
	}

	if result := database.Pg.Save(smtp); result.Error != nil {
		return result.Error
	}

	if result := database.Pg.Save(&models.AppMail{
		AppName:  oldSmtp.AppName,
		MailName: oldSmtp.MailName,
		MailType: smtpType.ToString(),
		Primary:  req.Primary,
	}); result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteSmtp deletes a existing smtp.
func DeleteSmtp(smtp *models.Smtp) error {
	if result := database.Pg.Delete(smtp); result.Error != nil {
		return result.Error
	}

	return nil
}
