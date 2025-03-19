package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ell534/goWebservers/internal/auth"
	"github.com/Ell534/goWebservers/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	newRequestBody := requestBody{}
	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
	}

	hashedPass, err := auth.HashPassword(newRequestBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
	}

	newUserParams := database.CreateUserParams{
		Email:          newRequestBody.Email,
		HashedPassword: hashedPass,
	}
	newUser, err := cfg.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new user", err)
	}

	response := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, response)
}
