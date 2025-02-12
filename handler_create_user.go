package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	newRequestBody := requestBody{}
	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
		return
	}

	newUser, err := cfg.db.CreateUser(r.Context(), newRequestBody.Email)

	response := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, response)
}
