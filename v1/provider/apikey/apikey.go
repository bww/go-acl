package apikey

import (
	"fmt"
	"net/http"

	"github.com/bww/go-acl/v1"
	"github.com/bww/go-acl/v1/provider"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-util/rand"
)

var errInvalidPolicyEffect = fmt.Errorf("Invalid policy effect")

func GenerateKey() (string, string) {
	return rand.RandomString(32), rand.RandomString(128)
}

type DataSource interface {
	AuthorizationForKeyAndSecret(key, secret string) (*acl.Authorization, error)
}

type APIKeyProvider struct {
	src DataSource
}

func NewProvider(src DataSource) *APIKeyProvider {
	return &APIKeyProvider{src}
}

func (p APIKeyProvider) Validate(req *http.Request) error {

	key, secret, err := acl.BasicAuthCredential(req)
	if err == acl.ErrNoAuthorization {
		return provider.ErrUnauthorized
	} else if err == acl.ErrUnsupportedMethod {
		return provider.ErrUnauthorized
	} else if err == acl.ErrMalformedRequest {
		return provider.ErrMalformed
	} else if err != nil {
		return err
	}

	a, err := p.src.AuthorizationForKeyAndSecret(key, secret)
	if err == dbx.ErrNotFound {
		return provider.ErrUnauthorized
	} else if err != nil {
		return err
	}

	if !a.Active {
		return provider.ErrForbidden
	}

	for _, e := range a.Policies {
		effect, err := e.Eval(req)
		if err != nil {
			return err
		}
		switch effect {
		case acl.Deny:
			return provider.ErrForbidden
		case acl.Allow:
			return nil
		case acl.Pass:
			continue
		default:
			return errInvalidPolicyEffect
		}
	}

	return provider.ErrForbidden // no policies matched
}

func (p APIKeyProvider) String() string {
	return "APIKey"
}
