package controllers

import (
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/dto/responses"
	"api-mail/main/src/errors"
	"api-mail/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateSmtp func for creating a new SMTP.
func CreateSmtp(c *fiber.Ctx) error {
	// Create a new smtp struct for the request.
	req := &requests.CreateSmtp{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate smtp fields.
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(req.App); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Check if smtp exists.
	if available, err := services.IsSmtpAvailable(req.App, req.Mail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.SmtpAvailable, "Smtp mail already exist.")
	}

	// Create smtp.
	if err := services.CreateSmtp(req); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the smtp.
	smtp, err := services.GetSmtp(req.App, req.Mail)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.Smtp{}
	response.SetSmtp(smtp)

	return c.JSON(response)
}

// UpdateSmtp func for updating a SMTP record.
func UpdateSmtp(c *fiber.Ctx) error {
	// Create a new smtp struct for the request.
	req := &requests.UpdateSmtp{}

	// Get the app and mail from the URL.
	app := c.Params("app")
	mail := c.Params("mail")

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate smtp fields.
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(app, mail)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.AppName == "" && smtp.MailName == "" {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	// Check if the smtp data has been modified since it was last fetched.
	if req.UpdatedAt.Unix() < smtp.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update smtp.
	if err := services.UpdateSmtp(smtp, req); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the smtp.
	smtp, err = services.GetSmtp(app, mail)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.Smtp{}
	response.SetSmtp(smtp)

	return c.JSON(response)
}

// DeleteSmtp func for deleting a SMTP record.
func DeleteSmtp(c *fiber.Ctx) error {
	// Get the app and mail from the URL.
	app := c.Params("app")
	mail := c.Params("mail")

	// Check if app exists.
	if available, err := services.IsAppAvailable(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(app, mail)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.AppName == "" && smtp.MailName == "" {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	// Delete the SMTP.
	if err := services.DeleteSmtp(smtp); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
