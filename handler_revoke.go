package main

import (
	"net/http"

	"github.com/peseb/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(rw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	respondWithJSON(rw, 204, nil)
}
