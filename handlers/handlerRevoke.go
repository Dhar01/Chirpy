package handlers

import (
	"net/http"

	"github.com/Dhar01/Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find token", err)
		return
	}

	if err := cfg.DB.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
