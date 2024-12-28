package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

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

func GetAPIKey(headers http.Header) (string, error) {
	key, err := getHeader(headers, "ApiKey")
	if err != nil {
		return "", err
	}

	return key, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token, err := getHeader(headers, "Bearer")
	if err != nil {
		return "", err
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func getHeader(headers http.Header, key string) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, key) {
		return "", errAuthHeaderNotFound
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, key))
	return token, nil
}
