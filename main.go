package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	hnd "github.com/Dhar01/Chirpy/handlers"
	"github.com/Dhar01/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	platform := getEnvVariable("PLATFORM")
	secretKey := getEnvVariable("SECRET_KEY")
	dbURL := getEnvVariable("DB_URL")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("can't connect to database: %v\n", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := hnd.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		SecretKey:      secretKey,
		Platform:       platform,
	}

	mux := http.NewServeMux()

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

	port := "8080"

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	if err = srv.ListenAndServe(); err != nil {
		log.Printf("%v\n", err.Error())
	}
}

func getEnvVariable(env string) string {
	envVar := os.Getenv(env)
	if envVar == "" {
		log.Fatalf("%s must be set", env)
	}

	return envVar
}
