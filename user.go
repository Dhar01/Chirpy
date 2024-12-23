package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type createUserRequest struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	person, err := cfg.queries.CreateUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "An error occurred when creating the user", http.StatusInternalServerError)
		return
	}

	user := User{
		ID: person.ID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
		Email: person.Email,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}
}
