package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type returnVal struct {
	// Valid *bool   `json:"valid,omitempty"`
	CleanBody string `json:"cleaned_body"`
	Error     string `json:"error,omitempty"`
}

func HandlerValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type chirp struct {
		Body string `json:"body"`
	}

	param := chirp{}

	// type returnVal struct {
	// 	// Valid *bool   `json:"valid,omitempty"`
	// 	CleanBody string `json:"cleaned_body"`
	// 	Error     string `json:"error,omitempty"`
	// }

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(param.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanBody := cleanBodyMsg(param.Body)
	response := returnVal{
		CleanBody: cleanBody,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	response := returnVal{
		Error: msg,
	}

	respondWithJSON(w, code, response)
}

func cleanBodyMsg(msg string) string {
	profane := []string{"kerfuffle", "sharbert", "fornax"}

	splitMsg := strings.Split(msg, " ")

	cleanWords := make([]string, 0, len(splitMsg))

	for _, word := range splitMsg {
		isProfane := false
		lowerWord := strings.ToLower(word)

		for _, profaneWord := range profane {
			if lowerWord == profaneWord {
				isProfane = true
				break
			}
		}

		if isProfane {
			cleanWords = append(cleanWords, "****")
		} else {
			cleanWords = append(cleanWords, word)
		}
	}

	return strings.Join(cleanWords, " ")
}
