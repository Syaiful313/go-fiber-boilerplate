package services

import (
	"errors"
	"fmt"
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) Register(req models.CreateUserRequest) (*models.RegisterResponse, error) {

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

	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("database error")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *AuthService) ForgotPassword(email string) error {

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

	resetToken, err := utils.GenerateResetPasswordToken(fmt.Sprintf("%d", user.ID), user.Email, s.cfg.JWTSecret)
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.FrontendURL, resetToken)

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

	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	claims, err := utils.ValidateResetPasswordToken(token, s.cfg.JWTSecret)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	var user models.User
	if err := database.GetDB().Where("id = ? AND email = ?", claims.UserID, claims.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return errors.New("database error")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := database.GetDB().Model(&user).Update("password", hashedPassword).Error; err != nil {
		return errors.New("failed to update password")
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
}
