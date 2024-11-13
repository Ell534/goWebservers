package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ell534/goWebservers/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.Queries
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	mux.HandleFunc("GET /api/reset", apiCfg.handleResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handleValidate)

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Listening on port 8080...")
	log.Fatal(myServer.ListenAndServe())
}
