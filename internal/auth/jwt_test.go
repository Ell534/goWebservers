package auth

import (
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
