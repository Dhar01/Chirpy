package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	if password == "" || hash == "" {
		return fmt.Errorf("password not provided")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}
