package controllers

import (
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/dto/responses"
	"api-mail/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateApp method to create an app.
func CreateApp(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateApp{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate document fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Create the app.
	app, err := services.CreateApp(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the document.
	response := responses.App{}
	response.SetApp(app)

	return c.JSON(response)
}
