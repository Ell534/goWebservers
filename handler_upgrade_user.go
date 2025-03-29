package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ell534/goWebservers/internal/auth"
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

	reqAPIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization header provided", err)
	}

	if reqAPIKey != cfg.polkaAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newRequestBody)
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
