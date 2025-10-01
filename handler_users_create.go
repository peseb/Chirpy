package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/peseb/Chirpy/internal/auth"
	"github.com/peseb/Chirpy/internal/database"
)

type UserDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(rw http.ResponseWriter, req *http.Request) {
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
	hashed_password, err := auth.HashPassword(dto.Password)
	if err != nil {
		respondWithError(rw, 500, "Unable to hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          dto.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(rw, 500, "Unable to create user", err)
		return
	}

	respondWithJSON(rw, 201, UserDto{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
