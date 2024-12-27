package handlers

import (
	"net/http"
	"strings"
)

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bearer := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(bearer, "Bearer"))

	if err := cfg.Queries.RevokeRefreshToken(r.Context(), token); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
