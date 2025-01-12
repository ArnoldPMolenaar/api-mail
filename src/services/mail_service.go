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

// IsPrimaryTypeAvailable method to check if a mail-type is available.
func IsPrimaryTypeAvailable(primaryType string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.AppMailPrimaryType{}, "name = ?", primaryType); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetAppMail finds a mail by app, mail.
func GetAppMail(app, mail string, preload ...bool) (models.AppMail, error) {
	var appMail models.AppMail
	query := database.Pg

	if len(preload) > 0 && preload[0] {
		query = query.Preload("Smtp").Preload("Gmail").Preload("Azure")
	}

	if result := database.Pg.First(&appMail, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return appMail, result.Error
	}

	return appMail, nil
}

// SendSmtpMail sends an email using SMTP.
func SendSmtpMail(appMail *models.AppMail, fromName, fromMail, to, subject, body, mimeType string, ccs []string, bccs []string) error {
	// SMTP record.
	var smtp *models.Smtp
	var err error

	if appMail.Smtp == nil {
		smtpID, err := GetSmtpIDByAppMailID(appMail.ID)
		if err != nil {
			return errors.New("smtp ID not found")
		}

		if isInCache, err := IsSmtpInCache(smtpID); err != nil {
			return err
		} else if isInCache {
			if smtp, err = GetSmtpFromCache(smtpID); err != nil {
				return err
			}
		} else {
			if smtp, err = GetSmtp(smtpID); err != nil {
				return err
			}
		}
	} else {
		smtp = appMail.Smtp
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
func CreateSendMail(appMail *models.AppMail, req *requests.SendMail) error {
	smtpType := enums.SMTP
	primaryType := smtpType.ToString()

	if appMail.PrimaryType.Valid {
		primaryType = &appMail.PrimaryType.String
	}

	sendMail := &models.SendMail{
		AppMailID:   appMail.ID,
		PrimaryType: *primaryType,
		FromName:    req.FromName,
		FromMail:    req.FromMail,
		To:          req.To,
		Subject:     req.Subject,
		Body:        req.Body,
		MimeType:    req.MimeType,
		Ccs:         make([]models.SendMailCc, 0),
		Bccs:        make([]models.SendMailBcc, 0),
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
