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
	if err := services.CreateSendMail(&appMail, sendMail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
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
			sendMail.Bccs); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, err.Error())
		}
	} else if sendGmailMail {
		// TODO: implement Gmail mail sending.
	} else if sendAzureMail {
		// TODO: implement Azure mail sending.
	} else {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.SendMail, "PrimaryType not found.")
	}

	return c.SendStatus(fiber.StatusCreated)
}
