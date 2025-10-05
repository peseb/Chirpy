package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/peseb/Chirpy/internal/database"
)

type EventData struct {
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPolkaWebhook(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string    `json:"event"`
		Data  EventData `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	dto := parameters{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(rw, 500, "Unable to decode DTO", err)
		return
	}

	if dto.Event != "user.upgraded" {
		respondWithJSON(rw, 204, nil)
		return
	}

	_, err = cfg.db.SetIsChirpyRed(req.Context(), database.SetIsChirpyRedParams{
		IsChirpyRed: true,
		ID:          dto.Data.UserId,
	})
	if err != nil {
		respondWithError(rw, 404, "Unable to set membership for the given user", err)
		return
	}

	respondWithJSON(rw, 204, nil)
}
