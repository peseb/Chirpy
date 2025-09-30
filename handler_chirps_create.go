package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/peseb/Chirpy/internal/database"
)

type ChirpDto struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerCreateChirp(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	dto := parameters{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(rw, 500, "Unable to decode DTO", err)
		return
	}

	if len(dto.Body) > 140 {
		respondWithError(rw, 400, "Chirp is too long", nil)
		return
	}

	cleanedBody := cleanInput(dto.Body)
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: dto.UserId,
	})
	if err != nil {
		respondWithError(rw, 500, "Unable to create chirp", err)
		return
	}

	respondWithJSON(rw, 201, ChirpDto{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
	})
}

func cleanInput(input string) string {
	parts := strings.Split(input, " ")
	invalidWords := []string{"kerfuffle", "sharbert", "fornax"}

	for i, word := range parts {
		newWord := word

		lowered := strings.ToLower(word)
		if slices.Contains(invalidWords, lowered) {
			newWord = "****"
		}

		parts[i] = newWord
	}

	cleaned := strings.Join(parts, " ")
	return cleaned
}
