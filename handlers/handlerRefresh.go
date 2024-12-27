package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bearer := r.Header.Get("Authorization")
	refreshToken := strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer"))

	userID, err := cfg.Queries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(userID, cfg.SecretKey, time.Hour)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Service error", http.StatusInternalServerError)
		return
	}
}
