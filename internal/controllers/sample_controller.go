package controllers

import (
	"strconv"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/pagination"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SampleController struct {
	sampleService *services.SampleService
}

func NewSampleController(cfg *config.Config) *SampleController {
	return &SampleController{
		sampleService: services.NewSampleService(cfg),
	}
}

func (h *SampleController) GetSamples(c *fiber.Ctx) error {
	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	params := pagination.NewParams(queryParams)

	samples, meta, err := h.sampleService.GetSamples(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch samples",
		})
	}

	var responses []models.SampleResponse
	for _, sample := range samples {
		responses = append(responses, sample.ToResponse())
	}

	return c.JSON(fiber.Map{
		"data": responses,
		"meta": meta,
	})
}

func (h *SampleController) GetSampleById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sample ID"})
	}

	sample, err := h.sampleService.GetSampleById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sample not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{"data": sample.ToResponse()})
}

func (h *SampleController) CreateSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req models.CreateSampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	sample, err := h.sampleService.CreateSample(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create sample"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sample created successfully",
		"data":    sample.ToResponse(),
	})
}

func (h *SampleController) UpdateSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sample ID"})
	}

	var req models.UpdateSampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	sample, err := h.sampleService.UpdateSample(userID, id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sample not found"})
		}
		if err.Error() == "forbidden" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only update your own samples"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{
		"message": "Sample updated successfully",
		"data":    sample.ToResponse(),
	})
}

func (h *SampleController) DeleteSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sample ID"})
	}

	if err := h.sampleService.DeleteSample(userID, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sample not found"})
		}
		if err.Error() == "forbidden" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only delete your own samples"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete sample"})
	}

	return c.JSON(fiber.Map{"message": "Sample deleted successfully"})
}
