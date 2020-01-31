package provider

import (
	"errors"
	"net/http"
)

var (
	ErrUnsupported         = errors.New("Authorization method is not supported")
	ErrMalformed           = errors.New("Authorization is malformed")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrExpired             = errors.New("Expired")
	ErrInsufficientEntropy = errors.New("Insufficient entropy")
	ErrForbidden           = errors.New("ErrForbidden")
)

type Provider interface {
	Validate(*http.Request) error
}
