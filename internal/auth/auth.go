package auth

import (
	"errors"
	"net/http"
)

// GetApiKey extract the api key from http request
func GetApiKey(r *http.Request) (string, error) {
	value := r.Header.Get("ApiKey")
	if value == "" {
		return "", errors.New("ApiKey not found")
	}

	return value, nil
}
