package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"strings"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodHead,
				fiber.MethodOptions,
			}, ","),
			AllowHeaders: "Accept,Content-Type",
		}),

		// Add simple logger.
		logger.New(),

		// Catch a panic and return a 500 response.
		recover.New(),
	)
}
