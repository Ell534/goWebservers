package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ell534/goWebservers/internal/auth"
	"github.com/Ell534/goWebservers/internal/database"
	"github.com/google/uuid"
)

type ValidChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	newChirp := parameters{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header format", err)
		return
	}

	bearerTokenUserID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	const maxChirpLength = 140
	if len(newChirp.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp length exceeds 140 characters", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := cleanChirp(newChirp.Body, badWords)

	cleanedChirp := database.CreateChirpParams{
		Body:   cleaned,
		UserID: bearerTokenUserID,
	}

	validChirp, err := cfg.db.CreateChirp(r.Context(), cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp", err)
		return
	}

	response := ValidChirp{
		ID:        validChirp.ID,
		CreatedAt: validChirp.CreatedAt,
		UpdatedAt: validChirp.UpdatedAt,
		Body:      validChirp.Body,
		UserID:    validChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func cleanChirp(chirp string, badWords map[string]struct{}) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
