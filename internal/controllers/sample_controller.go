package controllers

import (
	"strconv"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SampleHandler struct {
	sampleService *services.SampleService
}

func NewSampleController(cfg *config.Config) *SampleHandler {
	return &SampleHandler{
		sampleService: services.NewSampleService(cfg),
	}
}

func (h *SampleHandler) GetSamples(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	samples, total, err := h.sampleService.GetSamples(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch samples"})
	}

	var responses []models.SampleResponse
	for _, sample := range samples {
		responses = append(responses, sample.ToResponse())
	}

	return c.JSON(fiber.Map{
		"data": responses,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *SampleHandler) GetSample(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sample ID"})
	}

	sample, err := h.sampleService.GetSample(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sample not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{"data": sample.ToResponse()})
}

func (h *SampleHandler) CreateSample(c *fiber.Ctx) error {
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

func (h *SampleHandler) UpdateSample(c *fiber.Ctx) error {
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

func (h *SampleHandler) DeleteSample(c *fiber.Ctx) error {
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
