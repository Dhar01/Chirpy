package handlers

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodGet)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
