package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.handlerCreateChirps(w, r)
	case http.MethodGet:
		cfg.handlerGetChirps(w, r)
	case http.MethodDelete:
		cfg.handlerDeleteChirps(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpApi struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Body        string    `json:"body"`
	UserID      uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	var req createChirpRequest

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find JWT", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate JWT", err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", err)
		return
	}

	// req.Body validation needed
	if len(req.Body) > 140 {
		http.Error(w, "Chirp is too long", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: req.Body,
		// UserID: req.UserID,
		UserID: id,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpApi{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *ApiConfig) handlerDeleteChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: token not found", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized", err)
		return
	}

	chirpStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID", err)
		return
	}

	authId, err := cfg.DB.GetIDofChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "somehting wrong", err)
		return
	}

	if authId == id {
		if err := cfg.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
			ID:     chirpID,
			UserID: id,
		}); err != nil {
			respondWithError(w, http.StatusNotFound, "chirp doesn't exist", err)
			return
		}
	} else {
		respondWithError(w, http.StatusForbidden, "what's going on", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
