package main

import (
	"net/http"
	"time"

	"github.com/peseb/Chirpy/internal/auth"
)

type RefreshDto struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(rw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	db_token, err := cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	if db_token.RevokedAt.Valid {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	access_token, err := auth.MakeJWT(db_token.UserID, cfg.authSecret, 1*time.Hour)
	if err != nil {
		respondWithError(rw, 401, "Unauthorized", err)
		return
	}

	respondWithJSON(rw, 200, RefreshDto{
		Token: access_token,
	})
}
