package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) Register(req models.CreateUserRequest) (*models.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	var existingUser models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return nil, err
	}

	return &models.RegisterResponse{
		User: user.ToResponse(),
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		log.Printf("login database error for %s: %v", req.Email, err)
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		log.Printf("login blocked for inactive account: %s", user.Email)
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
	if err != nil {
		log.Printf("login token generation failed for %s: %v", user.Email, err)
		return nil, errors.New("invalid credentials")
	}

	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	if !utils.ValidateEmail(email) {
		return errors.New("invalid email format")
	}

	var user models.User
	if err := database.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			return nil
		}
		return errors.New("database error")
	}

	resetTokenValue, tokenHash, err := utils.GenerateResetToken(s.cfg.ResetTokenSecret)
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	// Invalidate previous tokens for this user
	database.GetDB().Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{})

	tokenRecord := models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	if err := database.GetDB().Create(&tokenRecord).Error; err != nil {
		return errors.New("failed to generate reset token")
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.FrontendURL, resetTokenValue)

	emailConfig := utils.EmailConfig{
		SMTPHost:     s.cfg.SMTPHost,
		SMTPPort:     s.cfg.SMTPPort,
		SMTPUsername: s.cfg.SMTPUsername,
		SMTPPassword: s.cfg.SMTPPassword,
		FromEmail:    s.cfg.FromEmail,
	}

	emailData := utils.EmailData{
		To:      user.Email,
		Subject: "Reset Your Password",
		Body:    utils.GenerateResetPasswordEmail(resetLink),
	}

	if err := utils.SendEmail(emailConfig, emailData); err != nil {
		return errors.New("failed to send reset email")
	}

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {

	if token == "" || newPassword == "" {
		return errors.New("token and new password are required")
	}

	if !utils.ValidatePassword(newPassword) {
		return errors.New("password must include upper, lower, number, special and be at least 8 characters")
	}

	rawToken, err := utils.VerifyResetToken(token, s.cfg.ResetTokenSecret)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	tokenHash := utils.HashResetToken(rawToken)

	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		var resetRecord models.PasswordResetToken
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", tokenHash).
			First(&resetRecord).Error; err != nil {
			return errors.New("invalid or expired reset token")
		}

		if resetRecord.Used || resetRecord.ExpiresAt.Before(time.Now()) {
			return errors.New("invalid or expired reset token")
		}

		var user models.User
		if err := tx.Where("id = ?", resetRecord.UserID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("invalid or expired reset token")
			}
			return errors.New("database error")
		}

		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			return errors.New("failed to hash password")
		}

		if err := tx.Model(&user).Update("password", hashedPassword).Error; err != nil {
			return errors.New("failed to update password")
		}

		if err := tx.Model(&resetRecord).Update("used", true).Error; err != nil {
			return errors.New("failed to update reset token")
		}

		emailConfig := utils.EmailConfig{
			SMTPHost:     s.cfg.SMTPHost,
			SMTPPort:     s.cfg.SMTPPort,
			SMTPUsername: s.cfg.SMTPUsername,
			SMTPPassword: s.cfg.SMTPPassword,
			FromEmail:    s.cfg.FromEmail,
		}

		emailData := utils.EmailData{
			To:      user.Email,
			Subject: "Password Reset Successful",
			Body:    utils.GeneratePasswordResetSuccessEmail(),
		}

		utils.SendEmail(emailConfig, emailData)

		return nil
	})
}

func validateRegisterRequest(req models.CreateUserRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return errors.New("all fields are required")
	}
	if !utils.ValidateEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if !utils.ValidatePassword(req.Password) {
		return errors.New("password must include upper, lower, number, special and be at least 8 characters")
	}
	return nil
}

func validateLoginRequest(req models.LoginRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}
	if !utils.ValidateEmail(req.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
