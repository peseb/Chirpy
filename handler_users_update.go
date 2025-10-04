package main

import (
	"encoding/json"
	"net/http"

	"github.com/peseb/Chirpy/internal/auth"
	"github.com/peseb/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	decoder := json.NewDecoder(req.Body)
	dto := parameters{}
	err = decoder.Decode(&dto)
	if err != nil {
		respondWithError(rw, 500, "Unable to decode DTO", err)
		return
	}

	hashed_password, err := auth.HashPassword(dto.Password)
	if err != nil {
		respondWithError(rw, 500, "Unable to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          dto.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(rw, 500, "Unable to create user", err)
		return
	}

	respondWithJSON(rw, 200, UserDto{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
