package main

import (
	"net/http"

	"github.com/Ell534/goWebservers/internal/auth"
)

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not extract token from headers", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke session", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
