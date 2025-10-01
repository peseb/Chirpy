package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.db.GetChirps(req.Context())
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
