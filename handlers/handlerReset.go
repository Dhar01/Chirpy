package handlers

import (
	"net/http"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	if cfg.Platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		w.Write([]byte("Reset only allowed in dev environment"))
		return
	}

	cfg.fileserverHits.Store(0)
	if err := cfg.DB.Reset(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't reset database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state"))
}
