package routes

import (
	"api-mail/main/src/controllers"
	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1", middleware.MachineProtected())

	// Register route for POST /v1/mail/send.
	route.Post("/mail/send", controllers.SendMail)

	// Register CRUD routes for /v1/smtp.
	// TODO: Add paginate combined for all mails.
	route.Post("/smtps", controllers.CreateSmtp)
	route.Get("/smtps/:id", controllers.GetSmtp)
	route.Put("/smtps/:id", controllers.UpdateSmtp)
	route.Delete("/smtps/:id", controllers.DeleteSmtp)
	route.Put("/smtps/:id/restore", controllers.RestoreSmtp)

	// Register CRUD routes for /v1/gmail.
	route.Post("/gmails", controllers.CreateGmail)
	route.Get("/gmails/:id", controllers.GetGmail)
	route.Put("/gmails/:id", controllers.UpdateGmail)
	route.Delete("/gmails/:id", controllers.DeleteGmail)
	route.Put("/gmails/:id/restore", controllers.RestoreGmail)
}
