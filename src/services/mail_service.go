package services

import (
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"errors"
	"fmt"
	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

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

// IsMailInSendMails method to check if a mail is in send-mails.
func IsMailInSendMails(mail string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.SendMail{}, "mail_name = ?", mail); result.Error != nil {
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

	if result := query.Order("\"primary\" DESC").First(&appMail); result.Error != nil {
		return appMail, result.Error
	}

	return appMail, nil
}

// GetAppMailsByMail finds all active mails by mail.
func GetAppMailsByMail(mail string) ([]models.AppMail, error) {
	var appMails []models.AppMail
	if result := database.Pg.Find(&appMails, "mail_name = ?", mail); result.Error != nil {
		return appMails, result.Error
	}

	return appMails, nil
}

// SendSmtpMail sends an email using SMTP.
func SendSmtpMail(appName, mailName, fromName, fromMail, to, subject, body, mimeType string, ccs []string, bccs []string) error {
	// SMTP record.
	var smtp *models.Smtp
	var err error

	if isInCache, err := IsSmtpInCache(appName, mailName); err != nil {
		return err
	} else if isInCache {
		if smtp, err = GetSmtpFromCache(appName, mailName); err != nil {
			return err
		}
	} else {
		if smtp, err = GetSmtp(appName, mailName); err != nil {
			return err
		}
	}

	if smtp != nil {
		if err = SetSmtpToCache(smtp); err != nil {
			return err
		}
	} else {
		return errors.New("smtp not found")
	}

	// SMTP server configuration.
	server := mail.NewSMTPClient()
	server.Host = smtp.Host
	server.Port = smtp.Port
	server.Username = smtp.Username

	// Decrypt password.
	password, err := smtp.DecryptPassword()
	if err != nil {
		return err
	}

	server.Password = password

	// SMTP Encryption.
	if smtp.Port == 465 {
		server.Encryption = mail.EncryptionSSLTLS
	} else if smtp.Port == 587 {
		server.Encryption = mail.EncryptionSTARTTLS
	} else {
		server.Encryption = mail.EncryptionNone
	}

	server.Authentication = mail.AuthLogin
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// SMTP client.
	client, err := server.Connect()
	if err != nil {
		return errors.New(fmt.Sprintf("smtp client error: %s", err.Error()))
	}

	// Email.
	var contentType mail.ContentType
	switch mimeType {
	case "text/plain":
		contentType = mail.TextPlain
	case "text/html":
		contentType = mail.TextHTML
	default:
		contentType = mail.TextHTML
	}

	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("%s <%s>", fromName, fromMail)).
		AddTo(to).
		SetSubject(subject).
		SetBody(contentType, body)

	if len(ccs) > 0 {
		email.AddCc(ccs...)
	}

	if len(bccs) > 0 {
		email.AddBcc(bccs...)
	}

	// Dkim.
	if smtp.DkimPrivateKey != nil && smtp.DkimDomain != nil && smtp.DkimCanonicalizationName != nil {
		var canonicalization string
		switch enums.ToDkimCanonicalization(*smtp.DkimCanonicalizationName) {
		case enums.Simple:
			canonicalization = "simple/simple"
		case enums.Relaxed:
			canonicalization = "relaxed/relaxed"
		}

		options := dkim.NewSigOptions()
		options.PrivateKey = []byte(*smtp.DkimPrivateKey)
		options.Domain = *smtp.DkimDomain
		options.Selector = "default"
		options.SignatureExpireIn = 3600
		options.Headers = []string{"from", "date", "mime-version", "received", "received"}
		options.AddSignatureTimestamp = true
		options.Canonicalization = canonicalization

		email.SetDkim(options)
	}

	if email.Error != nil {
		return errors.New(fmt.Sprintf("creating email error: %s", email.Error))
	}

	// Send email.
	if err := email.Send(client); err != nil {
		return errors.New(fmt.Sprintf("sending email error: %s", err.Error()))
	}

	return nil
}

// CreateSendMail creates a new send-mail.
func CreateSendMail(req *requests.SendMail) error {
	sendMail := &models.SendMail{
		AppName:  req.App,
		MailName: req.Mail,
		MailType: req.Type,
		FromName: req.FromName,
		FromMail: req.FromMail,
		To:       req.To,
		Subject:  req.Subject,
		Body:     req.Body,
		MimeType: req.MimeType,
		Ccs:      make([]models.SendMailCc, 0),
		Bccs:     make([]models.SendMailBcc, 0),
	}

	for _, cc := range req.Ccs {
		sendMail.Ccs = append(sendMail.Ccs, models.SendMailCc{Cc: cc})
	}

	for _, bcc := range req.Bccs {
		sendMail.Bccs = append(sendMail.Bccs, models.SendMailBcc{Bcc: bcc})
	}

	if result := database.Pg.Create(sendMail); result.Error != nil {
		return result.Error
	}

	return nil
}
