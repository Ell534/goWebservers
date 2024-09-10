package main

import "net/http"

func (cfg *apiConfig) handleResetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0!"))
}
