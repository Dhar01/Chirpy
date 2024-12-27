package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	expireTime, err := expireLimitSet(req.ExpiresAt)
	if err != nil {
		http.Error(w, "Bad expiration value", http.StatusBadRequest)
		return
	}

	user, err := cfg.Queries.GetUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "unauthorized attempt", http.StatusUnauthorized)
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		http.Error(w, "Unauthorized pass", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.SecretKey, expireTime)
	if err != nil {
		http.Error(w, "Malformed Token", http.StatusUnauthorized)
		return
	}

	person := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "can't encode the data", http.StatusInternalServerError)
		return
	}
}

func expireLimitSet(expire string) (time.Duration, error) {
	if expire == "" {
		return time.Hour, nil
	}

	seconds, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return 0, err
	}

	duration := time.Duration(seconds) * time.Second

	if duration < time.Hour {
		return duration, nil
	}

	return time.Hour, nil
}
