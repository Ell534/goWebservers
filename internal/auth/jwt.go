package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tokenClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	validatedClaims := &jwt.RegisteredClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, validatedClaims, keyFunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims type")
	}
	returnID := claims.Subject

	parsedUUID, err := uuid.Parse(returnID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID in token subject: %w", err)
	}

	return parsedUUID, nil
}
