package controllers

import (
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/enums"
	"api-mail/main/src/errors"
	"api-mail/main/src/services"
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
	if sendMail.Mail != "" {
		if available, err := services.IsMailAvailable(sendMail.Mail); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.MailExists, "MailName does not exist.")
		}
	}

	// Check if type exists.
	if sendMail.Type != "" {
		if available, err := services.IsMailTypeAvailable(sendMail.Type); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.MailTypeExists, "MailType does not exist.")
		}
	}

	// Get app mail.
	appMail, err := services.GetAppMail(sendMail.App, sendMail.Mail, sendMail.Type)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Send mail.
	mailType, err := enums.ToMailType(appMail.MailType)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.MailTypeExists, "MailType does not exist.")
	}

	switch mailType {
	case enums.SMTP:
		// TODO: implement SMTP mail sending.
	case enums.Gmail:
		// TODO: implement Gmail mail sending.
	case enums.Azure:
		// TODO: implement Azure mail sending.
	}

	return c.SendStatus(fiber.StatusCreated)
}
