package main

import (
	"fmt"
	"log"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database.ConnectDB(cfg)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		DisableStartupMessage: true, // Menonaktifkan pesan startup default Fiber
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.CORSMiddleware(cfg))

	// Setup routes
	routes.SetupRoutes(app, cfg)

	// Start server
	fmt.Printf("  ➜  [API] Local:   http://localhost:%s\n", cfg.Port) // Mencetak URL lokal saja
	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}

