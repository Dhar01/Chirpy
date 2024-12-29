package handlers

import (
	"log"
	"net/http"

	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if chirpID != "" {
		cfg.getSingleChirp(w, r, chirpID)
	} else if authorID != "" {
		cfg.getChirps(w, r, authorID, sortOrder)
	} else {
		cfg.getChirps(w, r, "", sortOrder)
	}
}

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request, authorID, sortOrder string) {
	var id uuid.UUID
	var err error

	if authorID == "" {
		id = uuid.Nil
	} else {
		id, err = uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid authorID", err)
			return
		}
	}

	// var chirps []database.Chirp
	// if sortOrder == "desc" {
	// 	chirps, err = cfg.DB.GetAllChirpsDESC(r.Context(), id)

	log.Printf("ID: %+v", id)
	// } else {

	chirps, err := cfg.sortedChirps(r, sortOrder, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find chirps", err)
		return
	}

	log.Printf("Number of chirps retrieved: %d", len(chirps))

	for i, chirp := range chirps {
		log.Printf("Chirp %d: %v", i, chirp)
	}

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

func (cfg *ApiConfig) sortedChirps(r *http.Request, sortOrder string, id uuid.UUID) ([]database.Chirp, error) {
	var chirps []database.Chirp
	var err error

	log.Printf("Sort order: %s, ID: %v", sortOrder, id)

	if sortOrder == "desc" {
		chirps, err = cfg.DB.GetAllChirpsDESC(r.Context(), id)
	} else {
		// handles both "asc" and empty/missing sort parameter
		chirps, err = cfg.DB.GetAllChirpsASC(r.Context(), id)
	}

	log.Printf("Retrieved chirps: %+v", chirps)

	if err != nil {
		return nil, err
	}
	return chirps, err
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
