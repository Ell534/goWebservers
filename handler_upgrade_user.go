package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	newRequestBody := requestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json body", err)
		return
	}

	userID, err := uuid.Parse(newRequestBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse userID string to uuid", err)
		return
	}

	if newRequestBody.Event == "user.upgraded" {
		_, err := cfg.db.UpgradeUser(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "user not found", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, nil)
	}

	w.WriteHeader(http.StatusNoContent)

}
