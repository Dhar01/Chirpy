package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if token == "" {
		t.Fatalf("Expected a token, but got an empty string")
	}
}

func TestValidateJWT(t *testing.T) {

}
