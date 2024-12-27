package handlers

import (
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find token", err)
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user from refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.SecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
