package main

import (
	"net/http"

	"github.com/Ell534/goWebservers/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "chirp ID not found in path", nil)
		return
	}
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing chirp id from path into UUID type", err)
		return
	}

	// get bearer token and retrieve userID from token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no access token in header found", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate access token", err)
		return
	}

	// retrive the chirp from the database using the provided chirpID
	// if the chirp is not found then return 404
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp with provided chirpID not found in database", err)
		return
	}

	// if tokenUserID and chirpUserID do not match, return 403
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "userID in bearer token does not match userID of chirp to delete", err)
		return
	}

	// if userIDs match then delete the chirp, need to write query
	// if deletion successful then return a 204
}
