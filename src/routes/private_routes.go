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
	route.Post("/smtps", controllers.CreateSmtp)
	route.Get("/smtps/:app/:mail", controllers.GetSmtp)
	route.Put("/smtps/:app/:mail", controllers.UpdateSmtp)
	route.Delete("/smtps/:app/:mail", controllers.DeleteSmtp)
	route.Put("/smtps/:app/:mail/restore", controllers.RestoreSmtp)
}
