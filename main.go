package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("can't connect to database: %v\n", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{
		queries: dbQueries,
	}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	// handlerAssets := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./app/assets")))
	// mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(handlerAssets))

	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/validate_chirp", handlerValidate)
	mux.HandleFunc("/api/users", apiCfg.createUserHandler)

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%v\n", err.Error())
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
