package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	newRequestBody := requestBody{}
	err := decoder.Decode(&newRequestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
	}
	// newRequestBody contains the email that will be used to query the database to retrieve the user
	// once the user is retrieved, it will have their hashed password which will be used along with
	// the password from the request body to plug into the checkpasswordhash function to determine if
	// the correct password has been supplied.
	// If the email or password are wrong just return 401 error, incorrect email/password.
}
