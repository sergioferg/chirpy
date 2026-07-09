package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	apiString := headers.Get("Authorization")
	if apiString == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	apiSlice := strings.Fields(apiString)
	if len(apiSlice) < 2 || apiSlice[0] != "ApiKey" {
		return "", errors.New("Malformed authorization header")
	}

	apiKeyString := apiSlice[1]
	return apiKeyString, nil
}
