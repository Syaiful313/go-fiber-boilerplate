package routes

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/handlers"
	"go-fiber-boilerplate/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	sampleHandler := handlers.NewSampleHandler()

	// Health check
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware(cfg))
	protected.Get("/profile", authHandler.GetProfile)

	// Sample routes (protected)
	samples := protected.Group("/samples")
	samples.Get("/", sampleHandler.GetSamples)
	samples.Post("/", sampleHandler.CreateSample)
	samples.Get("/:id", sampleHandler.GetSample)
	samples.Put("/:id", sampleHandler.UpdateSample)
	samples.Delete("/:id", sampleHandler.DeleteSample)
}

