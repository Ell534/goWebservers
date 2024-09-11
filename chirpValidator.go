package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// need to remove profanity
// lower case the string first
// split the string on whitespace, i.e " "
// check each word in the string slice
// ignore profanity when it is followed by punctuation, so check for exact matces ie current word == profanity word

func handleValidate(w http.ResponseWriter, req *http.Request) {
	type chirpFormat struct {
		Body string `json:"body"`
	}

	type validChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	chirp := chirpFormat{}
	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode chirp format")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "chirp is too long")
		return
	}

	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	chirpText := chirp.Body
	chirpTextSlice := strings.Split(chirpText, " ")
	fmt.Printf("before checking for profanity: %v", chirpTextSlice)
	for i, word := range chirpTextSlice {
		if slices.Contains(profanity, strings.ToLower(word)) {
			chirpTextSlice[i] = "****"
		}
	}
	fmt.Printf("after checking for profanity: %v", chirpTextSlice)
	cleanedChirp := strings.Join(chirpTextSlice, " ")

	respondWithJSON(w, http.StatusOK, validChirp{
		CleanedBody: cleanedChirp,
	})
}
