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

// Oauth2AzureCallback func for handling the Azure OAuth2 callback.
func Oauth2AzureCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	azureID, err := utils.StringToUint(state)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get azure.
	azure, err := services.GetAzureByID(azureID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Create OAuth2 config
	oauthConfig := services.CreateAzureOauthConfig(azure.ClientID, azure.TenantID, azure.Secret)

	// Exchange code for token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.OauthExchange, err.Error())
	}

	// Save the token into the database
	azure, err = services.UpdateAzureToken(azure, token)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.AzureCallback{}
	response.SetAzureCallback(azure)

	return c.JSON(response)
}

// GetAzure func for getting a Azure record.
func GetAzure(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Azure.
	azure, err := services.GetAzure(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if azure.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AzureExists, "Azure does not exist.")
	}

	response := responses.Azure{}
	response.SetAzure(azure)

	return c.JSON(response)
}

// CreateAzure func for creating a new Azure.
func CreateAzure(c *fiber.Ctx) error {
	// Create a new azure struct for the request.
	req := &requests.CreateAzure{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate azure fields.
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

	// Check if azure exists.
	if available, err := services.IsAzureAvailable(req.App, req.Mail); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AzureAvailable, "Azure mail already exist.")
	}

	// Create azure.
	azure, err := services.CreateAzure(req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateAzureOauthConfig(req.ClientID, req.TenantID, req.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(azure.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Azure{}
	response.SetAzure(azure, authCodeURL)

	return c.JSON(response)
}

// UpdateAzure func for updating a Azure record.
func UpdateAzure(c *fiber.Ctx) error {
	// Create a new azure struct for the request.
	req := &requests.UpdateAzure{}

	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate azure fields.
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Find the Azure.
	azure, err := services.GetAzure(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if azure.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AzureExists, "Azure does not exist.")
	}

	// Check if the azure data has been modified since it was last fetched.
	if req.UpdatedAt.Unix() < azure.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update azure.
	azure, err = services.UpdateAzure(azure, req)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateAzureOauthConfig(req.ClientID, req.TenantID, req.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(azure.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Azure{}
	response.SetAzure(azure, authCodeURL)

	return c.JSON(response)
}

// DeleteAzure func for deleting a azure record.
func DeleteAzure(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Azure.
	azure, err := services.GetAzure(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if azure.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AzureExists, "Azure does not exist.")
	}

	// Delete the Azure.
	if err := services.DeleteAzure(azure); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreAzure func for restoring a deleted Azure record.
func RestoreAzure(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Azure.
	azure, err := services.GetAzure(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if azure.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AzureExists, "Azure does not exist.")
	}

	// Restore the Azure.
	if err := services.RestoreAzure(azure); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	oauthConfig := services.CreateAzureOauthConfig(azure.ClientID, azure.TenantID, azure.Secret)
	authCodeURL := oauthConfig.AuthCodeURL(strconv.Itoa(int(azure.ID)), oauth2.AccessTypeOffline)

	// Return the url to request the token.
	response := responses.Azure{}
	response.SetAzure(azure, authCodeURL)

	return c.JSON(response)
}
