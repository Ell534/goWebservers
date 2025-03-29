package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")

	if authorIDStr != "" {
		userID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to parse author id string to uuid", err)
			return
		}

		chirpsByUser, err := cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "chirps by provided user not found in database", err)
			return
		}

		userChirps := []ValidChirp{}
		for _, chirp := range chirpsByUser {
			userChirps = append(userChirps, ValidChirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		respondWithJSON(w, http.StatusOK, userChirps)
	} else {

		allChirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve all chirps", err)
			return
		}

		chirps := []ValidChirp{}
		for _, chirp := range allChirps {
			chirps = append(chirps, ValidChirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		respondWithJSON(w, http.StatusOK, chirps)

	}
}
