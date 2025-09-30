package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func respondWithError(rw http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(rw, code, errorResponse{Error: msg})
}

func respondWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(code)
	rw.Write(dat)
}

func cleanInput(input string) string {
	parts := strings.Split(input, " ")
	invalidWords := []string{"kerfuffle", "sharbert", "fornax"}

	result := ""
	for _, word := range parts {
		newWord := word
		lowered := strings.ToLower(word)
		if slices.Contains(invalidWords, lowered) {
			newWord = "****"
		}

		separator := " "
		if len(result) == 0 {
			separator = ""
		}
		result = strings.Join([]string{result, newWord}, separator)
	}

	return result
}
