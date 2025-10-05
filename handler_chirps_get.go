package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/peseb/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")

	chirps, err := cfg.getChirps(req, authorId)
	if err != nil {
		respondWithError(rw, 500, "Unable to get chirps from database", err)
		return
	}

	chirpList := []ChirpDto{}
	for _, c := range chirps {
		chirpList = append(chirpList, ChirpDto{
			Id:        c.ID,
			UserId:    c.UserID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
		})
	}
	respondWithJSON(rw, 200, chirpList)
}

func (cfg *apiConfig) getChirps(req *http.Request, authorId string) ([]database.Chirp, error) {
	if authorId == "" {
		return cfg.db.GetChirps(req.Context())
	}

	userId, err := uuid.Parse(authorId)
	if err != nil {
		return nil, err
	}

	return cfg.db.GetChirpsByAuthor(req.Context(), userId)
}
