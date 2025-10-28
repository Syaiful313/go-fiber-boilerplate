package routes

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/controllers"
	"go-fiber-boilerplate/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupSampleRoutes(api fiber.Router, cfg *config.Config) {
	sampleController := controllers.NewSampleController(cfg)

	samples := api.Group("/samples")

	samples.Get("/", sampleController.GetSamples)
	samples.Get("/:id", middlewares.AuthMiddleware(cfg), sampleController.GetSampleById)
	samples.Post("/",
		middlewares.AuthMiddleware(cfg),
		middlewares.NewUploaderMiddleware().ImageUpload(2, []string{"image/jpeg", "image/png"}),
		sampleController.CreateSample)
	samples.Patch("/:id",
		middlewares.AuthMiddleware(cfg),
		middlewares.NewUploaderMiddleware().ImageUpload(2, []string{"image/jpeg", "image/png"}),
		sampleController.UpdateSample)
	samples.Delete("/:id", middlewares.AuthMiddleware(cfg), sampleController.DeleteSample)
}
