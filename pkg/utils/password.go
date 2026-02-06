package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimal 8 karakter
	if len(password) < 8 {
		return false
	}

	// Regex rules:
	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasNumber := regexp.MustCompile(`[0-9]`)
	hasSymbol := regexp.MustCompile(`[\W_]`)

	return hasUpper.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasNumber.MatchString(password) &&
		hasSymbol.MatchString(password)
}
