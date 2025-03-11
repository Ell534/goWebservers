package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ell534/goWebservers/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	newRequestBody := requestBody{}
	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body", err)
	}
	queriedUser, err := cfg.db.GetUserByEmail(r.Context(), newRequestBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
	}

	err = auth.CheckPasswordHash(newRequestBody.Password, queriedUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	response := User{
		ID:        queriedUser.ID,
		CreatedAt: queriedUser.CreatedAt,
		UpdatedAt: queriedUser.UpdatedAt,
		Email:     queriedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, response)
}
