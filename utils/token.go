package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const resetTokenLength = 32

// GenerateResetToken generates a random token, returns the signed token for clients and its hash for storage.
func GenerateResetToken(secret string) (signedToken string, tokenHash string, err error) {
	rawToken, err := GenerateRandomToken(resetTokenLength)
	if err != nil {
		return "", "", err
	}

	signature := signToken(rawToken, secret)
	signedToken = rawToken + "." + signature
	tokenHash = hashToken(rawToken)
	return signedToken, tokenHash, nil
}

// VerifyResetToken verifies the signature and returns the raw token part.
func VerifyResetToken(signedToken, secret string) (string, error) {
	parts := strings.Split(signedToken, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid token format")
	}

	expected := signToken(parts[0], secret)
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return "", errors.New("invalid token signature")
	}

	return parts[0], nil
}

// HashResetToken hashes the raw token for safe persistence.
func HashResetToken(rawToken string) string {
	return hashToken(rawToken)
}

// GenerateRandomToken returns a random hex-encoded string of the given byte length.
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func signToken(token, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
