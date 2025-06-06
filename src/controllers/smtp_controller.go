package controllers

import (
	"api-mail/main/src/database"
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/dto/responses"
	"api-mail/main/src/errors"
	"api-mail/main/src/models"
	"api-mail/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetSmtps func for getting all SMTP records.
func GetSmtps(c *fiber.Ctx) error {
	smtps := make([]models.Smtp, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":           true,
		"username":     true,
		"host":         true,
		"port":         true,
		"created_at":   true,
		"updated_at":   true,
		"primary_type": true,
		"app_name":     true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)

	db := database.Pg.Scopes(queryFunc, sortFunc).
		Limit(limit).
		Offset(offset).
		Preload("AppMail").
		Joins("JOIN \"app_mails\" ON \"app_mails\".\"id\" = \"app_mail_id\"").
		Find(&smtps)
	if db.Error != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, db.Error.Error())
	}

	total := int64(0)
	database.Pg.Scopes(queryFunc).
		Model(&models.Smtp{}).
		Joins("JOIN \"app_mails\" ON \"app_mails\".\"id\" = \"app_mail_id\"").
		Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), toSmtpPagination(smtps))

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetSmtp func for getting an SMTP record.
func GetSmtp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	response := responses.Smtp{}
	response.SetSmtp(smtp)

	return c.JSON(response)
}

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
	smtp, err := services.CreateSmtp(req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the smtp.
	if smtp != nil {
		response := responses.Smtp{}
		response.SetSmtp(smtp)

		return c.JSON(response)
	} else {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, "Failed to create smtp.")
	}
}

// UpdateSmtp func for updating a SMTP record.
func UpdateSmtp(c *fiber.Ctx) error {
	// Create a new smtp struct for the request.
	req := &requests.UpdateSmtp{}

	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate smtp fields.
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	// Check if the smtp data has been modified since it was last fetched.
	if req.UpdatedAt.Unix() < smtp.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update smtp.
	smtp, err = services.UpdateSmtp(smtp, req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the smtp.
	if smtp != nil {
		response := responses.Smtp{}
		response.SetSmtp(smtp)

		return c.JSON(response)
	} else {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, "Failed to update smtp.")
	}
}

// DeleteSmtp func for deleting a SMTP record.
func DeleteSmtp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	// Delete the SMTP.
	if err := services.DeleteSmtp(smtp); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreSmtp func for restoring a deleted SMTP record.
func RestoreSmtp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the SMTP.
	smtp, err := services.GetSmtp(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if smtp.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.SmtpExists, "Smtp does not exist.")
	}

	// Restore the SMTP.
	if err := services.RestoreSmtp(smtp); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// toSmtpPagination func for converting SMTPs to SMTP responses.
func toSmtpPagination(smtps []models.Smtp) []responses.Smtp {
	smtpResponses := make([]responses.Smtp, len(smtps))

	for i := range smtps {
		response := responses.Smtp{}
		response.SetSmtp(&smtps[i])
		smtpResponses[i] = response
	}

	return smtpResponses
}
