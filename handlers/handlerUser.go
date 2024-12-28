package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dhar01/Chirpy/internal/auth"
	"github.com/Dhar01/Chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.handlerCreateUser(w, r)
	case http.MethodPut:
		cfg.handlerUpdateUser(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (cfg *ApiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", err)
		return
	}

	hashPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashPasswd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	// expireTime, err := expireLimitSet(req.ExpiresAt)
	// if err != nil {
	// 	http.Error(w, "Bad expiration value", http.StatusBadRequest)
	// 	return
	// }

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyREd: user.IsChirpyRed.Bool,
	})
}

func (cfg *ApiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if accessToken == "" || err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find JWT token", err)
		return
	}

	id, err := auth.ValidateJWT(accessToken, cfg.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate JWT", err)
		return
	}

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request", err)
		return
	}

	hashPass, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash the pass", err)
		return
	}

	if err := cfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		HashedPassword: hashPass,
		ID:             id,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update password", err)
		return
	}

	if err := cfg.DB.UpdateEmail(r.Context(), database.UpdateEmailParams{
		Email: req.Email,
		ID:    id,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update email", err)
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find user with email", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyREd: user.IsChirpyRed.Bool,
	})
}
