package handlers

import (
	"encoding/json"
	"net/http"

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

	user, err := cfg.Queries.GetUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "unauthorized attempt", http.StatusUnauthorized)
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	person := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "can't encode the data", http.StatusInternalServerError)
		return
	}
}
