package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Listening on port 8080...")
	myServer.ListenAndServe()
}
