package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	checkerMethod(w, r, http.MethodPost)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	value := fmt.Sprintf(
		`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`,
		cfg.FileserverHits.Load())

	w.Write([]byte(value))
}
