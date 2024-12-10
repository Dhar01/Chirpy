package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	err := srv.ListenAndServe()
	if err != nil {
		log.Printf(err.Error())
	}
}
