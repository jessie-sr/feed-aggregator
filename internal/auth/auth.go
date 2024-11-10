package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extract the API key from the headers of an http request
// Example format:
// Authorization: Apikey {insert apikey here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header: the first part should be 'ApiKey'")
	}

	return vals[1], nil
}
