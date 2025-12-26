package routes

import (
	"time"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/controllers"
	"go-fiber-boilerplate/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	auth := api.Group("/auth")

	keyGen := func(c *fiber.Ctx) string {
		return c.IP()
	}

	auth.Post("/register",
		middlewares.RateLimitMiddleware(5, time.Minute, keyGen),
		authController.Register)

	auth.Post("/login",
		middlewares.RateLimitMiddleware(10, time.Minute, keyGen),
		authController.Login)

	auth.Post("/forgot-password",
		middlewares.RateLimitMiddleware(5, time.Minute, keyGen),
		authController.ForgotPassword)

	auth.Post("/reset-password",
		middlewares.RateLimitMiddleware(5, time.Minute, keyGen),
		authController.ResetPassword)
}
