package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	tokenSecret := "superSecret"
	userID := uuid.New()
	expiresIn := 5 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedID != userID {
		t.Errorf("expected userID %v, got %v", userID, validatedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	tokenSecret := "superSecret"
	userID := uuid.New()
	expiresIn := -1 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)

	if err == nil {
		t.Errorf("expected an error for expired token, but got none")
	}
}

func TestInvalidSecret(t *testing.T) {
	tokenSecret := "superSecret"
	wrongSecret := "wrongSecret"
	userID := uuid.New()
	expiresIn := 5 * time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)

	if err == nil {
		t.Errorf("expected an error for invalid secret, but got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid bearer token",
			authHeader:    "Bearer abc123xyz456",
			expectedToken: "abc123xyz456",
			expectError:   false,
		},
		{
			name:          "Bearer token with extra spaces",
			authHeader:    "Bearer   xyz789   ",
			expectedToken: "xyz789",
			expectError:   false,
		},
		{
			name:          "Missing authorization token",
			authHeader:    "",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "Invalid header format, no Bearer prefix",
			authHeader:    "Token abc123",
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(req.Header)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("expected token: %q, got: %q", tt.expectedToken, token)
				}
			}
		})
	}
}
