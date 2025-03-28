package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ell534/goWebservers/internal/auth"
	"github.com/Ell534/goWebservers/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not extract token from headers", err)
		return
	}

	accessTokenUserID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	newRequestBody := requestBody{}
	err = decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(newRequestBody.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}
	userCredentials := database.UpdateUserCredentialsParams{
		Email:          newRequestBody.Email,
		HashedPassword: hashedPassword,
		ID:             accessTokenUserID,
	}
	updatedUser, err := cfg.db.UpdateUserCredentials(r.Context(), userCredentials)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user credentials in database", err)
		return
	}

	response := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, response)
}
