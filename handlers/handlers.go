package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/Dhar01/Chirpy/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Queries        *database.Queries
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
