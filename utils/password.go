package utils

import "golang.org/x/crypto/bcrypt"

func ValidatePassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, r := range pw {
		switch {
		case 'A' <= r && r <= 'Z':
			hasUpper = true
		case 'a' <= r && r <= 'z':
			hasLower = true
		case '0' <= r && r <= '9':
			hasNumber = true
		case isSpecial(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func isSpecial(r rune) bool {
	switch {
	case r >= '!' && r <= '/':
		return true
	case r >= ':' && r <= '@':
		return true
	case r >= '[' && r <= '`':
		return true
	case r >= '{' && r <= '~':
		return true
	default:
		return false
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
