package routes

import (
	"api-mail/main/src/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// OAuth2 callbacks
	oauth := route.Group("/oauth2")
	oauth.Get("/gmails/callback", controllers.Oauth2GmailCallback)
	oauth.Get("/azures/callback", controllers.Oauth2AzureCallback)
}
