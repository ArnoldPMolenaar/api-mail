package services

import (
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/models"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io"
	"strings"
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

	if result := database.Pg.Find(&appMail, "app_name = ? AND mail_name = ?", app, mail); result.Error != nil {
		return appMail, result.Error
	}

	return appMail, nil
}

// SendSmtpMail sends an email using SMTP.
func SendSmtpMail(appMail *models.AppMail, fromName, fromMail, to, subject, body, mimeType string, ccs []string, bccs []string, attachments []requests.SendMailAttachment) error {
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

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			email.Attach(&mail.File{
				Name:     attachment.FileName,
				MimeType: attachment.FileType,
				Data:     attachment.FileData,
				Inline:   false,
			})
		}
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

func SendGmailMail(appMail *models.AppMail, fromName, fromMail, to, subject, body, mimeType string, ccs []string, bccs []string) error {
	// Gmail record.
	var gmailRecord *models.Gmail
	var err error
	ctx := context.Background()

	if appMail.Gmail == nil {
		gmailID, err := GetGmailIDByAppMailID(appMail.ID)
		if err != nil {
			return errors.New("gmail ID not found")
		}

		if isInCache, err := IsGmailInCache(gmailID); err != nil {
			return err
		} else if isInCache {
			if gmailRecord, err = GetGmailFromCache(gmailID); err != nil {
				return err
			}
		} else {
			if gmailRecord, err = GetGmail(gmailID); err != nil {
				return err
			}
		}
	} else {
		gmailRecord = appMail.Gmail
	}

	if gmailRecord != nil {
		if err = SetGmailToCache(gmailRecord); err != nil {
			return err
		}
	} else {
		return errors.New("gmail not found")
	}

	// Create OAuth2 config.
	if !gmailRecord.AccessToken.Valid ||
		!gmailRecord.RefreshToken.Valid ||
		!gmailRecord.TokenType.Valid ||
		!gmailRecord.Expiry.Valid ||
		!gmailRecord.ExpiresIn.Valid {
		return errors.New("gmail record not authenticated")
	}

	oauthConfig := CreateGmailOauthConfig(gmailRecord.ClientID, gmailRecord.Secret)
	client := oauthConfig.Client(ctx, &oauth2.Token{
		AccessToken:  gmailRecord.AccessToken.String,
		TokenType:    gmailRecord.TokenType.String,
		RefreshToken: gmailRecord.RefreshToken.String,
		Expiry:       gmailRecord.Expiry.Time,
		ExpiresIn:    gmailRecord.ExpiresIn.Int64,
	})

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting gmail service: %s", err.Error()))
	}

	// Create the message.
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", fromName, fromMail)
	header["To"] = to
	header["Cc"] = strings.Join(ccs, ",")
	header["Bcc"] = strings.Join(bccs, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf(`%s; charset="utf-8"`, mimeType)
	header["Content-Transfer-Encoding"] = "base64"

	var msg string
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\n", k, v)
	}
	msg += "\n" + body

	gMsg := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msg)),
	}

	// Send the message
	_, err = gmailService.Users.Messages.Send("me", &gMsg).Do()
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending gmail message: %s", err.Error()))
	}

	return nil
}

func SendAzureMail(appMail *models.AppMail, to, subject, body, mimeType string, ccs []string, bccs []string) error {
	// Azure record.
	var azure *models.Azure
	var err error
	ctx := context.Background()

	if appMail.Azure == nil {
		azureID, err := GetAzureIDByAppMailID(appMail.ID)
		if err != nil {
			return errors.New("azure ID not found")
		}

		if isInCache, err := IsAzureInCache(azureID); err != nil {
			return err
		} else if isInCache {
			if azure, err = GetAzureFromCache(azureID); err != nil {
				return err
			}
		} else {
			if azure, err = GetAzure(azureID); err != nil {
				return err
			}
		}
	} else {
		azure = appMail.Azure
	}

	if azure != nil {
		if err = SetAzureToCache(azure); err != nil {
			return err
		}
	} else {
		return errors.New("azure not found")
	}

	// Create OAuth2 config.
	if !azure.AccessToken.Valid ||
		!azure.RefreshToken.Valid ||
		!azure.TokenType.Valid ||
		!azure.Expiry.Valid ||
		!azure.ExpiresIn.Valid {
		return errors.New("azure not authenticated")
	}

	oauthConfig := CreateAzureOauthConfig(azure.ClientID, azure.TenantID, azure.Secret)
	client := oauthConfig.Client(ctx, &oauth2.Token{
		AccessToken:  azure.AccessToken.String,
		TokenType:    azure.TokenType.String,
		RefreshToken: azure.RefreshToken.String,
		Expiry:       azure.Expiry.Time,
		ExpiresIn:    azure.ExpiresIn.Int64,
	})

	// Create the email.
	var contentType graphmodels.BodyType
	switch mimeType {
	case "text/plain":
		contentType = graphmodels.TEXT_BODYTYPE
	case "text/html":
		contentType = graphmodels.HTML_BODYTYPE
	default:
		contentType = graphmodels.TEXT_BODYTYPE
	}

	requestBody := graphusers.NewItemSendMailPostRequestBody()
	message := graphmodels.NewMessage()
	message.SetSubject(&subject)
	itemBody := graphmodels.NewItemBody()
	itemBody.SetContentType(&contentType)
	itemBody.SetContent(&body)
	message.SetBody(itemBody)

	// To recipients.
	recipient := graphmodels.NewRecipient()
	emailAddress := graphmodels.NewEmailAddress()
	emailAddress.SetAddress(&to)
	recipient.SetEmailAddress(emailAddress)
	toRecipients := []graphmodels.Recipientable{
		recipient,
	}
	message.SetToRecipients(toRecipients)

	// Cc recipients.
	for _, cc := range ccs {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		emailAddress.SetAddress(&cc)
		recipient.SetEmailAddress(emailAddress)
		ccRecipients := []graphmodels.Recipientable{
			recipient,
		}
		message.SetCcRecipients(ccRecipients)
	}

	// Bcc recipients.
	for _, bcc := range bccs {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		emailAddress.SetAddress(&bcc)
		recipient.SetEmailAddress(emailAddress)
		bccRecipients := []graphmodels.Recipientable{
			recipient,
		}
		message.SetBccRecipients(bccRecipients)
	}

	requestBody.SetMessage(message)
	saveToSentItems := false
	requestBody.SetSaveToSentItems(&saveToSentItems)

	writer := jsonserialization.NewJsonSerializationWriter()
	if err := requestBody.Serialize(writer); err != nil {
		return err
	}
	requestBodyJson, err := writer.GetSerializedContent()
	if err != nil {
		return errors.New(fmt.Sprintf("Error serializing request body: %s", err.Error()))
	}
	requestBodyJson = append([]byte("{"), append(requestBodyJson, '}')...)

	// Send the mail via microsoft graph
	resp, err := client.Post("https://graph.microsoft.com/v1.0/me/sendMail", "application/json", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return errors.New(fmt.Sprintf("Error while sending mail: %s", err.Error()))
	}

	// Expect a 202 or throw an error
	if resp.StatusCode != 202 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Error sending mail: %s, Body: %s", resp.Status, string(bodyBytes)))
	}

	if err = resp.Body.Close(); err != nil {
		return errors.New(fmt.Sprintf("Error closing response body: %s", err.Error()))
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
		Attachments: make([]models.SendMailAttachment, 0),
	}

	for _, cc := range req.Ccs {
		sendMail.Ccs = append(sendMail.Ccs, models.SendMailCc{Cc: cc})
	}

	for _, bcc := range req.Bccs {
		sendMail.Bccs = append(sendMail.Bccs, models.SendMailBcc{Bcc: bcc})
	}

	for _, attachment := range req.Attachments {
		sendMail.Attachments = append(sendMail.Attachments, models.SendMailAttachment{
			FileName: attachment.FileName,
			FileType: attachment.FileType,
			FileSize: attachment.FileSize,
			FileData: attachment.FileData,
		})
	}

	if result := database.Pg.Create(sendMail); result.Error != nil {
		return result.Error
	}

	return nil
}
