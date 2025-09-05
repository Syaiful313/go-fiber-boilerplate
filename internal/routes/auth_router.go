package routes

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	auth := api.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
}
