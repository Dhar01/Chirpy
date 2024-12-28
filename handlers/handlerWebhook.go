package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dhar01/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type eventHook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *ApiConfig) HandlerWebhooks(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	var hook eventHook
	upgrade := "user.upgraded"

	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode the response", err)
		return
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil || key != cfg.PaymentKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized access request", err)
		return
	}


	userID, err := uuid.Parse(hook.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "bad ID form", err)
		return
	}

	if hook.Event != upgrade {
		http.Error(w, "wrong event", http.StatusNoContent)
		return
	} else {
		if err := cfg.DB.SetMemberShip(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusNotFound, "couldn't update membership", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
