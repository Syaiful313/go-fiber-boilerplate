package middleware

import (
	"go-fiber-boilerplate/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware(cfg *config.Config) fiber.Handler {
	// If credentials are allowed, origins cannot be wildcard
	if cfg.AllowCredentials && cfg.AllowedOrigins == "*" {
		// For development, disable credentials when using wildcard
		return cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: false,
		})
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: cfg.AllowCredentials,
	})
}

