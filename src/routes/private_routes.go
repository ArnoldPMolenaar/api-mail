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

	// Register CRUD routes for /v1/smtps.
	smtps := route.Group("/smtps", middleware.MachineProtected())
	smtps.Get("/", controllers.GetSmtps)
	smtps.Post("/", controllers.CreateSmtp)
	smtps.Get("/:id", controllers.GetSmtp)
	smtps.Put("/:id", controllers.UpdateSmtp)
	smtps.Delete("/:id", controllers.DeleteSmtp)
	smtps.Put("/:id/restore", controllers.RestoreSmtp)

	// Register CRUD routes for /v1/gmails.
	gmails := route.Group("/gmails", middleware.MachineProtected())
	gmails.Get("/", controllers.GetGmails)
	gmails.Post("/", controllers.CreateGmail)
	gmails.Get("/:id", controllers.GetGmail)
	gmails.Put("/:id", controllers.UpdateGmail)
	gmails.Delete("/:id", controllers.DeleteGmail)
	gmails.Put("/:id/restore", controllers.RestoreGmail)

	// Register CRUD routes for /v1/azures.
	azures := route.Group("/azures", middleware.MachineProtected())
	azures.Post("/", controllers.CreateAzure)
	azures.Get("/:id", controllers.GetAzure)
	azures.Put("/:id", controllers.UpdateAzure)
	azures.Delete("/:id", controllers.DeleteAzure)
	azures.Put("/:id/restore", controllers.RestoreAzure)
}
