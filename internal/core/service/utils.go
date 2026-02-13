package service

import (
	"github.com/badoux/checkmail"
	"github.com/golang-auth/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isValidEmail(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return domain.ErrInvalidEmail
	}
	if err := checkmail.ValidateHost(email); err != nil {
		return domain.ErrInvalidEmail
	}
	return nil
}
