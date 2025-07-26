package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	Port             string
	JWTSecret        string
	AllowedOrigins   string
	AllowCredentials bool
	
	// Email configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	
	// Frontend URL for reset password links
	FrontendURL string
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	allowCredentials := getEnv("CORS_ALLOW_CREDENTIALS", "false") == "true"

	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "admin"),
		DBName:           getEnv("DB_NAME", "go_fiber_db"),
		Port:             getEnv("PORT", "8000"),
		JWTSecret:        getEnv("JWT_SECRET", "default_secret"),
		AllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "*"),
		AllowCredentials: allowCredentials,
		
		// Email configuration
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@example.com"),
		
		// Frontend URL
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

