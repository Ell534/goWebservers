package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ell534/goWebservers/internal/auth"
	"github.com/Ell534/goWebservers/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	const JWTExpiration = time.Hour * 1
	refreshTokenExpiration := time.Now().AddDate(0, 0, 60)

	decoder := json.NewDecoder(r.Body)

	newRequestBody := requestBody{}

	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body", err)
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

	usersToken, err := auth.MakeJWT(queriedUser.ID, cfg.jwtSecret, time.Duration(JWTExpiration))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
	}

	newRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    queriedUser.ID,
		ExpiresAt: refreshTokenExpiration,
	}

	newRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), newRefreshTokenParams)

	response := User{
		ID:           queriedUser.ID,
		CreatedAt:    queriedUser.CreatedAt,
		UpdatedAt:    queriedUser.UpdatedAt,
		Email:        queriedUser.Email,
		Token:        usersToken,
		RefreshToken: newRefreshToken.Token,
	}

	respondWithJSON(w, http.StatusOK, response)
}
