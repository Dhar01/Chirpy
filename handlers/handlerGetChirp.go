package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	authorID := r.URL.Query().Get("author_id")

	if chirpID != "" {
		cfg.getSingleChirp(w, r, chirpID)
	} else if authorID != "" {
		cfg.getChirps(w, r, authorID)
	} else {
		cfg.getChirps(w, r, "")
	}
}

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request, authorID string) {
	var id uuid.UUID
	var err error

	if authorID == "" {
		id = uuid.Nil
	} else {
		id, err = uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "couldn't decode authorID", err)
			return
		}
	}

	chirps, err := cfg.DB.GetAllChirps(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find chirps", err)
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

	respondWithJSON(w, http.StatusOK, chirpApis)
}

func (cfg *ApiConfig) getSingleChirp(w http.ResponseWriter, r *http.Request, chirpID string) {
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.DB.GetSingleChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't found chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, ChirpApi{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
