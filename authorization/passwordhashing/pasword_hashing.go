package passwordhashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordDoesNotMatch = errors.New("password does not match")
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash with error: %w", err)
	}

	return string(hash), nil
}

func ValidatePasswordHash(password, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordDoesNotMatch
		}
		return err
	}
	return nil
}
