package main

import (
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpId")

	if chirpIdStr == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID not found in path", nil)
	}

	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing chirp id from path into UUID type", err)
	}

	retrievedChirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error retrieving chirp from database", err)
	}

	chirpPayload := ValidChirp{
		ID:        retrievedChirp.ID,
		CreatedAt: retrievedChirp.CreatedAt,
		UpdatedAt: retrievedChirp.UpdatedAt,
		Body:      retrievedChirp.Body,
		UserID:    retrievedChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpPayload)

}
