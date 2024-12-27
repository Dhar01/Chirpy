package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirps", err)
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

	chirp, err := cfg.Queries.GetSingleChirp(r.Context(), id)
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
