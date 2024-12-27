package handlers

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	SecretKey      string
	Platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type response struct {
	User         User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type createUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	// ExpiresAt string `json:"expires_in_seconds"`
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func checkerMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) {
	if r.Method != allowedMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
