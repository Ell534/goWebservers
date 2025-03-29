package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	apiPrefix := "ApiKey "
	if !strings.HasPrefix(authHeader, apiPrefix) {
		return "", fmt.Errorf("Invalid authorization header format")
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, apiPrefix))

	return apiKey, nil
}
