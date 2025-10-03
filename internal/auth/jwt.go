package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Makes a new JWT and returns it as a string

	// Builds the JWT with specified claims and signing algorithm
	newJWT := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  &jwt.NumericDate{time.Now().UTC()},
			ExpiresAt: &jwt.NumericDate{time.Now().UTC().Add(expiresIn)},
			Subject:   userID.String(),
		},
	)

	// Signs the JWT with the secret
	JWT, err := newJWT.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	// Returns the JWT as a string
	return JWT, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	// Validates an input token string and returns the token bearer's UUID if it's valid

	// Specifies the key function to pass the secret to the parser
	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	// Parses the token string
	out, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return uuid.UUID{}, err
	}

	// Checks if the token if valid
	if out.Valid {
		// Checks if the claims fit the structure we want (just jwt.RegisteredClaims)
		if clms, ok := out.Claims.(*jwt.RegisteredClaims); ok {
			// Gets the UUID from the token's claims
			uid, err := uuid.Parse(clms.Subject)
			if err != nil {
				return uuid.UUID{}, fmt.Errorf("Failed to parse UUID: %s", clms.Subject)
			}
			// Returns the UUID
			return uid, nil
		}
		return uuid.UUID{}, fmt.Errorf("Failed to assert claims type")
	}
	return uuid.UUID{}, fmt.Errorf("Invalid token")
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return token, fmt.Errorf("Token not found")
	}
	tknFields := strings.Fields(token)
	return tknFields[1], nil
}
