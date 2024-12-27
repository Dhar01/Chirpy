package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := time.Hour

	t.Run("testing MakeJWT", func(t *testing.T) {
		token, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if token == "" {
			t.Fatalf("Expected a token, but got an empty string")
		}
	})
	t.Run("testing validateJWT", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		id, err := ValidateJWT(tokenString, tokenSecret)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if userID != id {
			t.Errorf("Expected UUIDs to be same, but got id = %v, original = %v", id, userID)
		}
	})
}
