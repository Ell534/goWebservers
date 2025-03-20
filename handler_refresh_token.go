package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Ell534/goWebservers/internal/auth"
)

type refreshTokenResponse struct {
	RefreshToken string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not extract token from headers", err)
		return
	}

	queriedToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "token not found in database", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve token from database", err)
		return
	}
	if queriedToken.ExpiresAt.Before(time.Now()) || queriedToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token is expired or has been revoked", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), queriedToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve user from token", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access token", err)
		return
	}

	response := refreshTokenResponse{
		RefreshToken: newAccessToken,
	}

	respondWithJSON(w, http.StatusOK, response)

}
