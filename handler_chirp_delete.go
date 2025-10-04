package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/peseb/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(rw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.authSecret)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	id := req.PathValue("id")
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

	if chirp.UserID != userId {
		respondWithError(rw, 403, "Unauthorized", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(rw, 500, "Unable to delete chirp", err)
		return
	}

	respondWithJSON(rw, 204, nil)
}
