package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ResetTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "reset_password"
	jwt.RegisteredClaims
}

// GenerateResetPasswordToken generates a JWT token for password reset
func GenerateResetPasswordToken(userID, email, secretKey string) (string, error) {
	claims := ResetTokenClaims{
		UserID: userID,
		Email:  email,
		Type:   "reset_password",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 hour expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateResetPasswordToken validates and parses a reset password token
func ValidateResetPasswordToken(tokenString, secretKey string) (*ResetTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ResetTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ResetTokenClaims); ok && token.Valid {
		// Check if token type is correct
		if claims.Type != "reset_password" {
			return nil, fmt.Errorf("invalid token type")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateRandomToken generates a random token for additional security
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

