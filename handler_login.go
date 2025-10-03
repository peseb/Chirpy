package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/peseb/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int    `json:"expires_in_seconds"`
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

	expiresIn := dto.ExpiresIn
	if expiresIn == 0 || expiresIn > 3600 {
		expiresIn = 3600
	}
	duration := time.Duration(expiresIn) * time.Second
	token, err := auth.MakeJWT(user.ID, cfg.authSecret, duration)
	if err != nil {
		respondWithError(rw, 401, "Incorrect email or password", err)
		return
	}

	respondWithJSON(rw, 200, UserDto{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
