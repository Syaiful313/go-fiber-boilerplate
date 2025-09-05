package routes

import (
    "go-fiber-boilerplate/config"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
    app.Get("/api/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok", 
            "message": "Server is running",
        })
    })

    api := app.Group("/")

    SetupAuthRoutes(api, cfg)
	SetupSampleRoutes(api, cfg)
}