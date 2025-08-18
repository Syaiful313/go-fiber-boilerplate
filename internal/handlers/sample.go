package handlers

import (
	"strconv"

	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SampleHandler struct{}

func NewSampleHandler() *SampleHandler {
	return &SampleHandler{}
}

func (h *SampleHandler) GetSamples(c *fiber.Ctx) error {
	var samples []models.Sample

	// Get query parameters for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Get samples with user information
	if err := database.GetDB().
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&samples).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch samples",
		})
	}

	// Convert to response format
	var responses []models.SampleResponse
	for _, sample := range samples {
		responses = append(responses, sample.ToResponse())
	}

	// Get total count
	var total int64
	database.GetDB().Model(&models.Sample{}).Count(&total)

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sample ID",
		})
	}

	var sample models.Sample
	if err := database.GetDB().
		Preload("User").
		First(&sample, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sample not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.JSON(fiber.Map{
		"data": sample.ToResponse(),
	})
}

func (h *SampleHandler) CreateSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req models.CreateSampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	sample := models.Sample{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}

	if err := database.GetDB().Create(&sample).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create sample",
		})
	}

	// Load user information
	database.GetDB().Preload("User").First(&sample, sample.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sample created successfully",
		"data":    sample.ToResponse(),
	})
}

func (h *SampleHandler) UpdateSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sample ID",
		})
	}

	var req models.UpdateSampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var sample models.Sample
	if err := database.GetDB().First(&sample, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sample not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check if user owns the sample
	if sample.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only update your own samples",
		})
	}

	// Update fields
	if req.Title != "" {
		sample.Title = req.Title
	}
	if req.Description != "" {
		sample.Description = req.Description
	}

	if err := database.GetDB().Save(&sample).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update sample",
		})
	}

	// Load user information
	database.GetDB().Preload("User").First(&sample, sample.ID)

	return c.JSON(fiber.Map{
		"message": "Sample updated successfully",
		"data":    sample.ToResponse(),
	})
}

func (h *SampleHandler) DeleteSample(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sample ID",
		})
	}

	var sample models.Sample
	if err := database.GetDB().First(&sample, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sample not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check if user owns the sample
	if sample.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own samples",
		})
	}

	if err := database.GetDB().Delete(&sample).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sample",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sample deleted successfully",
	})
}

