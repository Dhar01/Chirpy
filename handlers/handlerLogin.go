package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dhar01/Chirpy/internal/auth"
	"github.com/Dhar01/Chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	// expireTime, err := expireLimitSet(req.ExpiresAt)
	// if err != nil {
	// 	http.Error(w, "Bad expiration value", http.StatusBadRequest)
	// 	return
	// }

	user, err := cfg.DB.GetUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.SecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
		return
	}

	if err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Refreshtoken: refreshToken,
		UserID:       user.ID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyREd:  user.IsChirpyRed.Bool,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

// func expireLimitSet(expire string) (time.Duration, error) {
// 	if expire == "" {
// 		return time.Hour, nil
// 	}
// 	seconds, err := strconv.ParseInt(expire, 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}
// 	duration := time.Duration(seconds) * time.Second
// 	if duration < time.Hour {
// 		return duration, nil
// 	}
// 	return time.Hour, nil
// }
