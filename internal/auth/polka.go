package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	reqToken := headers.Get("Authorization")
	if reqToken == "" {
		return "", errors.New("authorization header does not exist")
	}

	token := strings.Split(reqToken, "ApiKey")
	if len(token) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if token[0] != "" {
		return "", errors.New("authorization header must start with 'ApiKey'")
	}

	return strings.TrimSpace(token[1]), nil
}
