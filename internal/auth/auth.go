package auth

import (
	"errors"
	"net/http"
	"strings"
)

// get api key extracts an api from the headers of an http request
// e.g Authorization: ApiKey {insert api key here}
func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("no authentication info found")
	}

	value := strings.Split(authorization, " ")
	if len(value) != 2 || value[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return value[1], nil
}
