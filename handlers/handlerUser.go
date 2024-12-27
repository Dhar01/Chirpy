package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type createUserRequest struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	ExpiresAt string `json:"expires_in_seconds"`
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	hashPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "can't process password", http.StatusInternalServerError)
		return
	}

	person, err := cfg.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashPasswd,
	})
	if err != nil {
		http.Error(w, "An error occurred when creating the user", http.StatusInternalServerError)
		return
	}

	expireTime, err := expireLimitSet(req.ExpiresAt)
	if err != nil {
		http.Error(w, "Bad expiration value", http.StatusBadRequest)
		return
	}

	token, err := auth.MakeJWT(person.ID, cfg.SecretKey, expireTime)
	if err != nil {
		http.Error(w, "Bad expiration value", http.StatusBadRequest)
		return
	}

	user := User{
		ID:        person.ID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
		Email:     person.Email,
		Token:     token,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}
