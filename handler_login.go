package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ell534/goWebservers/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	const oneHourInSeconds = 3600
	var expiresIn int

	decoder := json.NewDecoder(r.Body)

	newRequestBody := requestBody{}

	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body", err)
	}

	if newRequestBody.ExpiresInSeconds == nil {
		expiresIn = oneHourInSeconds
	} else {
		if *newRequestBody.ExpiresInSeconds > oneHourInSeconds {
			expiresIn = oneHourInSeconds
		} else {
			expiresIn = oneHourInSeconds
		}
	}

	queriedUser, err := cfg.db.GetUserByEmail(r.Context(), newRequestBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(newRequestBody.Password, queriedUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	usersToken, err := auth.MakeJWT(queriedUser.ID, cfg.jwtSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create JWT", err)
		return
	}

	response := User{
		ID:        queriedUser.ID,
		CreatedAt: queriedUser.CreatedAt,
		UpdatedAt: queriedUser.UpdatedAt,
		Email:     queriedUser.Email,
		Token:     usersToken,
	}

	respondWithJSON(w, http.StatusOK, response)
}
