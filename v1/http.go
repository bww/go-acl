package acl

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNoAuthorization   = fmt.Errorf("No authorization")
	ErrMalformedRequest  = fmt.Errorf("Authorization is malformed")
	ErrUnsupportedMethod = fmt.Errorf("Authorization method is not supported")
)

// Obtain the basic auth credential for the specified request
func BasicAuthCredential(req *http.Request) (string, string, error) {

	auth := req.Header.Get("Authorization")
	if auth != "" {
		return BasicAuthCredentialFromAuthorization(auth)
	}

	auth = req.URL.Query().Get("auth")
	if auth != "" {
		return BasicAuthCredentialFromAuthorizationData(auth)
	}

	return "", "", ErrNoAuthorization
}

// Obtain the basic auth credential for the specified authorization header
func BasicAuthCredentialFromAuthorization(auth string) (string, string, error) {

	header := strings.Fields(auth)
	if len(header) != 2 {
		return "", "", ErrMalformedRequest
	}

	method := strings.ToLower(header[0])
	if method != "basic" {
		return "", "", ErrUnsupportedMethod
	}

	return BasicAuthCredentialFromAuthorizationData(header[1])
}

// Obtain the basic auth credential for the specified authorization header data
func BasicAuthCredentialFromAuthorizationData(auth string) (string, string, error) {

	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", "", ErrMalformedRequest
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return "", "", ErrMalformedRequest
	}

	return parts[0], parts[1], nil
}
