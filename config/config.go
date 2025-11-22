package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	Port                string
	JWTSecret           string
	ResetTokenSecret    string
	AllowedOrigins      string
	AllowCredentials    bool
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	FromEmail           string
	FrontendURL         string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var err error
	cfg := &Config{}

	if cfg.DBHost, err = getRequiredEnv("DB_HOST"); err != nil {
		return nil, err
	}
	if cfg.DBPort, err = getRequiredEnv("DB_PORT"); err != nil {
		return nil, err
	}
	if cfg.DBUser, err = getRequiredEnv("DB_USER"); err != nil {
		return nil, err
	}
	if cfg.DBPassword, err = getRequiredEnv("DB_PASSWORD"); err != nil {
		return nil, err
	}
	if cfg.DBName, err = getRequiredEnv("DB_NAME"); err != nil {
		return nil, err
	}
	if cfg.Port, err = getRequiredEnv("PORT"); err != nil {
		return nil, err
	}
	if cfg.JWTSecret, err = getRequiredEnv("JWT_SECRET"); err != nil {
		return nil, err
	}
	if cfg.ResetTokenSecret, err = getRequiredEnv("RESET_TOKEN_SECRET"); err != nil {
		return nil, err
	}
	if cfg.AllowedOrigins, err = getRequiredEnv("CORS_ALLOWED_ORIGINS"); err != nil {
		return nil, err
	}
	allowCredentials, err := getRequiredEnv("CORS_ALLOW_CREDENTIALS")
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(allowCredentials) {
	case "true":
		cfg.AllowCredentials = true
	case "false":
		cfg.AllowCredentials = false
	default:
		return nil, fmt.Errorf("CORS_ALLOW_CREDENTIALS must be 'true' or 'false'")
	}

	if cfg.SMTPHost, err = getRequiredEnv("SMTP_HOST"); err != nil {
		return nil, err
	}
	if cfg.SMTPPort, err = getRequiredEnv("SMTP_PORT"); err != nil {
		return nil, err
	}
	if cfg.SMTPUsername, err = getRequiredEnv("SMTP_USERNAME"); err != nil {
		return nil, err
	}
	if cfg.SMTPPassword, err = getRequiredEnv("SMTP_PASSWORD"); err != nil {
		return nil, err
	}
	if cfg.FromEmail, err = getRequiredEnv("FROM_EMAIL"); err != nil {
		return nil, err
	}
	if cfg.FrontendURL, err = getRequiredEnv("FRONTEND_URL"); err != nil {
		return nil, err
	}
	if cfg.CloudinaryCloudName, err = getRequiredEnv("CLOUDINARY_CLOUD_NAME"); err != nil {
		return nil, err
	}
	if cfg.CloudinaryAPIKey, err = getRequiredEnv("CLOUDINARY_API_KEY"); err != nil {
		return nil, err
	}
	if cfg.CloudinaryAPISecret, err = getRequiredEnv("CLOUDINARY_API_SECRET"); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "default_secret" {
		return errors.New("JWT_SECRET must not use insecure default")
	}
	if c.ResetTokenSecret == c.JWTSecret {
		return errors.New("RESET_TOKEN_SECRET must differ from JWT_SECRET")
	}
	if c.AllowCredentials && c.AllowedOrigins == "*" {
		return errors.New("CORS_ALLOWED_ORIGINS cannot be '*' when credentials are allowed")
	}
	return nil
}

func getRequiredEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("%s environment variable is required", key)
}
