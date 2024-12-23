package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

type createChirpRequest struct {
	Body    string    `json:"body"`
	User_ID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) HandlerChirps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createChirpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "can't decode the body", http.StatusBadRequest)
		return
	}

	if len(req.Body) > 140 {
		http.Error(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	info, err := cfg.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   req.Body,
		UserID: req.User_ID,
	})

	if err != nil {
		http.Error(w, "can't create chirp", http.StatusInternalServerError)
		return
	}

	chirp := Chirp{
		ID:        info.ID,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		Body:      info.Body,
		UserID:    info.UserID,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(chirp); err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}
