package routes

import (
    "go-fiber-boilerplate/config"
    "go-fiber-boilerplate/internal/controllers"
    "github.com/gofiber/fiber/v2"
)

func SetupSampleRoutes(api fiber.Router, cfg *config.Config) {
	sampleController := controllers.NewSampleController(cfg)
	
	samples := api.Group("/samples")
	
	samples.Get("/", sampleController.GetSamples)
	samples.Get("/:id", sampleController.GetSample)
	samples.Post("/", sampleController.CreateSample)
	samples.Put("/:id", sampleController.UpdateSample)
	samples.Delete("/:id", sampleController.DeleteSample)
}