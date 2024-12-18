package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	// handlerAssets := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./app/assets")))
	// mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(handlerAssets))

	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/validate_chirp", handlerValidate)

	err := srv.ListenAndServe()
	if err != nil {
		log.Printf(err.Error())
	}
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	value := fmt.Sprintf(
		`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`,
		cfg.fileserverHits.Load())

	w.Write([]byte(value))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type chirp struct {
		Body string `json:"body"`
	}

	param := chirp{}

	type returnVal struct {
		// Valid *bool   `json:"valid,omitempty"`
		CleanBody string `json:"cleaned_body"`
		Error     string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		writeJSON(w, http.StatusInternalServerError, returnVal{Error: "Something went wrong"})
		return
	}

	if len(param.Body) > 140 {
		writeJSON(w, http.StatusBadRequest, returnVal{Error: "Chirp is too long"})
		return
	}

	validTrue := true
	writeJSON(w, http.StatusOK, returnVal{Valid: &validTrue})
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

}

func cleanBodyMsg(msg string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}

	msg = strings.ToLower(msg)
	splitMsg := strings.Split(msg, " ")


	for _, msg := range splitMsg {
		for _, fane := range profane {
			if msg == fane {
				strings.Join("****", " ")
			} else {
				strings.Join(msg, " ")
			}
		}
	}

	return msg
}