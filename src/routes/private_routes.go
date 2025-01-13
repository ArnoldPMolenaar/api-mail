package routes

import (
	"api-mail/main/src/controllers"
	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register route for POST /v1/mail/send.
	route.Post("/mail/send", middleware.MachineProtected(), controllers.SendMail)

	// Register CRUD routes for /v1/smtp.
	// TODO: Add paginate combined for all mails.
	smtps := route.Group("/smtps", middleware.MachineProtected())
	smtps.Post("/", controllers.CreateSmtp)
	smtps.Get("/:id", controllers.GetSmtp)
	smtps.Put("/:id", controllers.UpdateSmtp)
	smtps.Delete("/:id", controllers.DeleteSmtp)
	smtps.Put("/:id/restore", controllers.RestoreSmtp)

	// Register CRUD routes for /v1/gmail.
	gmails := route.Group("/gmails", middleware.MachineProtected())
	gmails.Post("/", controllers.CreateGmail)
	gmails.Get("/:id", controllers.GetGmail)
	gmails.Put("/:id", controllers.UpdateGmail)
	gmails.Delete("/:id", controllers.DeleteGmail)
	gmails.Put("/:id/restore", controllers.RestoreGmail)

	// OAuth2 callbacks
	oauth := route.Group("/oauth2")
	oauth.Get("/gmails/callback", controllers.Oauth2GmailCallback)
}
