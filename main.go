package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	apiCfg := apiConfig{fileServerHits: 0}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handleResetMetrics)

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Listening on port 8080...")
	log.Fatal(myServer.ListenAndServe())
}
