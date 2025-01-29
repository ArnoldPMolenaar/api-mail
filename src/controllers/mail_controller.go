package controllers

import (
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/errors"
	"api-mail/main/src/services"
	"database/sql"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// SendMail func for sending mail.
func SendMail(c *fiber.Ctx) error {
	// Create a new mail struct for the request.
	sendMail := &requests.SendMail{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(sendMail); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate sendMail fields.
	validate := utils.NewValidator()
	if err := validate.Struct(sendMail); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Validate each attachment.
	for _, attachment := range sendMail.Attachments {
		// Validate FileType.
		if !isValidMimeType(attachment.FileType) {
			return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, "Invalid file type")
		}

		// Validate FileData.
		if len(attachment.FileData) == 0 {
			return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, "File data is empty")
		}
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(sendMail.App); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Check if mail exists.
	if available, err := services.IsMailAvailable(sendMail.Mail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.MailExists, "MailName does not exist.")
	}

	// Check if type exists.
	if sendMail.Type != nil {
		if available, err := services.IsPrimaryTypeAvailable(*sendMail.Type); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.MailTypeExists, "MailType does not exist.")
		}
	}

	// Get app mail.
	appMail, err := services.GetAppMail(sendMail.App, sendMail.Mail)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Check if primary type is set.
	sendSmtpMail := false
	sendGmailMail := false
	sendAzureMail := false

	primaryType := enums.ToAppMailPrimaryType(&appMail.PrimaryType.String)
	if sendMail.Type != nil {
		primaryType = enums.ToAppMailPrimaryType(sendMail.Type)
	}
	switch *primaryType {
	case enums.SMTP:
		sendSmtpMail = true
	case enums.Gmail:
		sendGmailMail = true
	case enums.Azure:
		sendAzureMail = true
	default:
		appMail, err := services.GetAppMail(sendMail.App, sendMail.Mail, true)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		}

		if appMail.Azure != nil {
			sendAzureMail = true
		} else if appMail.Gmail != nil {
			sendGmailMail = true
		} else {
			sendSmtpMail = true
		}
	}

	if sendAzureMail {
		azureType := enums.Azure
		appMail.PrimaryType = sql.NullString{String: *azureType.ToString(), Valid: true}
	} else if sendGmailMail {
		gmailType := enums.Gmail
		appMail.PrimaryType = sql.NullString{String: *gmailType.ToString(), Valid: true}
	} else {
		smtpType := enums.SMTP
		appMail.PrimaryType = sql.NullString{String: *smtpType.ToString(), Valid: true}
	}

	// Create mail.
	if !sendMail.DisableSave {
		if err := services.CreateSendMail(&appMail, sendMail); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		}
	}

	// Send mail.
	if sendSmtpMail {
		if err := services.SendSmtpMail(
			&appMail,
			sendMail.FromName,
			sendMail.FromMail,
			sendMail.To,
			sendMail.Subject,
			sendMail.Body,
			sendMail.MimeType,
			sendMail.Ccs,
			sendMail.Bccs,
			sendMail.Attachments); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, err.Error())
		}
	} else if sendGmailMail {
		if err := services.SendGmailMail(
			&appMail,
			sendMail.FromName,
			sendMail.FromMail,
			sendMail.To,
			sendMail.Subject,
			sendMail.Body,
			sendMail.MimeType,
			sendMail.Ccs,
			sendMail.Bccs); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, err.Error())
		}
	} else if sendAzureMail {
		if err := services.SendAzureMail(
			&appMail,
			sendMail.To,
			sendMail.Subject,
			sendMail.Body,
			sendMail.MimeType,
			sendMail.Ccs,
			sendMail.Bccs); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, err.Error())
		}
	} else {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, "PrimaryType not found.")
	}

	return c.SendStatus(fiber.StatusCreated)
}

// isValidMimeType checks if the provided MIME type is valid.
func isValidMimeType(mimeType string) bool {
	// Add your valid MIME types here.
	validMimeTypes := map[string]bool{
		"application/pdf":    true,
		"image/jpeg":         true,
		"image/png":          true,
		"text/plain":         true,
		"text/html":          true,
		"application/msword": true, // .doc
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
		"application/vnd.ms-excel": true, // .xls
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true, // .xlsx
		"application/vnd.ms-powerpoint":                                             true, // .ppt
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
		"application/zip":              true,
		"application/x-rar-compressed": true,
		"application/x-7z-compressed":  true,
		"application/x-tar":            true,
	}

	return validMimeTypes[mimeType]
}
