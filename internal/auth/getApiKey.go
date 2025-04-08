package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header not found")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Authorization header format must be 'ApiKey {token}'")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))
	if token == "" {
		return "", fmt.Errorf("Token not found in Authorization header")
	}

	return token, nil
}
