package main

import (
	"fmt"
	"log"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/middlewares"
	"go-fiber-boilerplate/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	database.ConnectDB(cfg)

	app := fiber.New(fiber.Config{
		ErrorHandler:          middlewares.ErrorHandler,
		DisableStartupMessage: true,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middlewares.CORSMiddleware(cfg))

	routes.SetupRoutes(app, cfg)

	fmt.Printf("  ➜  [API] Local:   http://localhost:%s\n", cfg.Port)
	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}
