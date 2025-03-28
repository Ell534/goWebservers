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

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp with provided chirpID not found in database", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "userID in bearer token does not match userID of chirp to delete", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
