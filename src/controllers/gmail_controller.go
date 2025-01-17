package controllers

import (
	"api-mail/main/src/dto/requests"
	"api-mail/main/src/dto/responses"
	"api-mail/main/src/errors"
	"api-mail/main/src/services"
	"context"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"strconv"
)

// Oauth2GmailCallback func for handling the Gmail OAuth2 callback.
func Oauth2GmailCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	gmailID, err := utils.StringToUint(state)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get gmail.
	gmail, err := services.GetGmailByID(gmailID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Create OAuth2 config
	oauthConfig := services.CreateGmailOauthConfig(gmail.ClientID, gmail.Secret)

	// Exchange code for token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.OauthExchange, err.Error())
	}

	// Save the token into the database
	gmail, err = services.UpdateGmailToken(gmail, token)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.GmailCallback{}
	response.SetGmailCallback(gmail)

	return c.JSON(response)
}

// GetGmail func for getting a Gmail record.
func GetGmail(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Gmail.
	gmail, err := services.GetGmail(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if gmail.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.GmailExists, "Gmail does not exist.")
	}

	response := responses.Gmail{}
	response.SetGmail(gmail)

	return c.JSON(response)
}

// CreateGmail func for creating a new Gmail.
func CreateGmail(c *fiber.Ctx) error {
	// Create a new gmail struct for the request.
	req := &requests.CreateGmail{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate gmail fields.
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

	// Check if gmail exists.
	if available, err := services.IsGmailAvailable(req.App, req.Mail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.GmailAvailable, "Gmail mail already exist.")
	}

	// Create gmail.
	gmail, err := services.CreateGmail(req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateGmailOauthConfig(req.ClientID, req.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(gmail.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Gmail{}
	response.SetGmail(gmail, authCodeURL)

	return c.JSON(response)
}

// UpdateGmail func for updating a Gmail record.
func UpdateGmail(c *fiber.Ctx) error {
	// Create a new gmail struct for the request.
	req := &requests.UpdateGmail{}

	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate gmail fields.
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Find the Gmail.
	gmail, err := services.GetGmail(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if gmail.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.GmailExists, "Gmail does not exist.")
	}

	// Check if the gmail data has been modified since it was last fetched.
	if req.UpdatedAt.Unix() < gmail.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update gmail.
	gmail, err = services.UpdateGmail(gmail, req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateGmailOauthConfig(req.ClientID, req.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(gmail.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Gmail{}
	response.SetGmail(gmail, authCodeURL)

	return c.JSON(response)
}

// DeleteGmail func for deleting a Gmail record.
func DeleteGmail(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Gmail.
	gmail, err := services.GetGmail(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if gmail.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.GmailExists, "Gmail does not exist.")
	}

	// Delete the Gmail.
	if err := services.DeleteGmail(gmail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreGmail func for restoring a deleted Gmail record.
func RestoreGmail(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Gmail.
	gmail, err := services.GetGmail(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if gmail.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.GmailExists, "Gmail does not exist.")
	}

	// Restore the Gmail.
	if err := services.RestoreGmail(gmail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateGmailOauthConfig(gmail.ClientID, gmail.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(gmail.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Gmail{}
	response.SetGmail(gmail, authCodeURL)

	return c.JSON(response)
}
