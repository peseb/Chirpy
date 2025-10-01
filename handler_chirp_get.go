package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(rw http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	fmt.Printf("ID: %s\n", id)
	parsedId, err := uuid.Parse(id)

	if err != nil {
		respondWithError(rw, 400, "Invalid id. Must be a valid UUID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), parsedId)
	if err != nil {
		respondWithError(rw, 404, "Unable to find chirp", err)
		return
	}

	respondWithJSON(rw, 200, ChirpDto{
		Id:        chirp.ID,
		UserId:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
	})
}
