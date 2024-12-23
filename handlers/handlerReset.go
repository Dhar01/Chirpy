package handlers

import (
	"net/http"
	"os"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if platform := os.Getenv("PLATFORM"); platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := cfg.Queries.DeleteAllUsers(r.Context()); err != nil {
		http.Error(w, "An error occurred while resetting users", http.StatusInternalServerError)
		return
	}

	// cfg.fileserverHits.Store(0)

	w.WriteHeader(http.StatusOK)
}
