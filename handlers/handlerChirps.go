package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.handlerCreateChirps(w, r)
	case http.MethodGet:
		cfg.handlerGetChirps(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpApi struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
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
		UserID: req.UserID,
	})

	if err != nil {
		http.Error(w, "can't create chirp", http.StatusInternalServerError)
		return
	}

	chirp := ChirpApi{
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

func (cfg *ApiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	if chirpID != "" {
		cfg.getSingleChirp(w, r, chirpID)
	} else {
		cfg.getAllChirps(w, r)
	}
}

func (cfg *ApiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Queries.GetAllChirps(r.Context())

	if err != nil {
		http.Error(w, "can't retrieve chirps", http.StatusInternalServerError)
		return
	}

	// for i, chrip := range chirps {
	// 	log.Printf("Chirp %d: %v", i, chrip)
	// }

	chirpApis := make([]ChirpApi, len(chirps))
	for i, chirp := range chirps {
		chirpApis[i] = ChirpApi{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(chirpApis); err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}

func (cfg *ApiConfig) getSingleChirp(w http.ResponseWriter, r *http.Request, chirpID string) {
	id, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.Queries.GetSingleChirp(r.Context(), id)
	if err != nil {
		http.Error(w, "chirp not exist!", http.StatusNotFound)
		return
	}

	info := ChirpApi{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "can't write to console", http.StatusInternalServerError)
		return
	}
}
