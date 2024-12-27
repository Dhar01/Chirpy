package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	hnd "github.com/Dhar01/Chirpy/handlers"
	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("can't connect to database: %v\n", err)
	}

	secretKey := os.Getenv("SECRET_KEY")

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	apiCfg := hnd.ApiConfig{
		Queries:   dbQueries,
		SecretKey: secretKey,
	}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir("./app")))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(handler))

	// handlerAssets := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./app/assets")))
	// mux.Handle("/app/assets/", apiCfg.middlewareMetricsInci(handlerAssets))

	mux.HandleFunc("/admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("/admin/reset", apiCfg.HandlerReset)

	mux.HandleFunc("/api/healthz", hnd.HandlerReadiness)
	mux.HandleFunc("/api/users", apiCfg.HandlerCreateUser)
	mux.HandleFunc("/api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("/api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("/api/revoke", apiCfg.HandlerRevoke)
	mux.HandleFunc("/api/chirps", apiCfg.HandlerChirps)
	mux.HandleFunc("/api/chirps/{chirpID}", apiCfg.HandlerChirps)

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
}
