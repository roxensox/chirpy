package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	// Finds the API key in the header, if provided

	// Gets the strings located at Authorization in header
	authorization := headers.Get("Authorization")

	// Returns error if nothing is found
	if authorization == "" {
		return "", fmt.Errorf("No authorization found")
	}

	// Splits result into its fields
	authParts := strings.Fields(authorization)

	// If the title isn't correct, returns error
	if authParts[0] != "ApiKey" {
		return "", fmt.Errorf("Incorrect authorization: %s should be ApiKey", authParts[0])
	}

	// Returns the key if it's found, otherwise returns error
	if len(authParts) > 1 {
		return authParts[1], nil
	} else {
		return "", fmt.Errorf("No API Key provided")
	}
}
