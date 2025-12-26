package controllers

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(cfg),
	}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := ctrl.authService.Register(req)
	if err != nil {
		switch err.Error() {
		case "user with this email already exists",
			"all fields are required",
			"invalid email format",
			"password must include upper, lower, number, special and be at least 8 characters":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to process registration",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    response,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := ctrl.authService.Login(req)
	if err != nil {
		switch err.Error() {
		case "email and password are required", "invalid email format":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"data":    response,
	})
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req models.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := ctrl.authService.ForgotPassword(req.Email)
	if err != nil {
		if err.Error() == "email is required" || err.Error() == "invalid email format" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to process request",
		})
	}

	return c.JSON(fiber.Map{
		"message": "If the email exists, a reset link has been sent",
	})
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req models.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := ctrl.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		switch err.Error() {
		case "token and new password are required", "password must include upper, lower, number, special and be at least 8 characters":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "invalid or expired reset token":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "user not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to reset password",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}
