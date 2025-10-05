package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/peseb/Chirpy/internal/auth"
	"github.com/peseb/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	dto := parameters{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(rw, 500, "Unable to decode DTO", err)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), dto.Email)
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	passwordIsCorrect, err := auth.CheckPasswordHash(dto.Password, user.HashedPassword)
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	if !passwordIsCorrect {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	duration := 1 * time.Hour
	access_token, err := auth.MakeJWT(user.ID, cfg.authSecret, duration)
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	db_token, err := cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	respondWithJSON(rw, 200, UserDto{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        access_token,
		RefreshToken: db_token.Token,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
