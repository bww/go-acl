package bearer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bww/go-acl/v1"
	"github.com/bww/go-acl/v1/provider"

	"github.com/bww/go-util/v1/crypto"
)

var errInvalidPolicyEffect = fmt.Errorf("Invalid policy effect")

const (
	methodBearer = "Bearer"
	sep          = "$"
)

type BearerProvider struct {
	key []byte
}

func NewProvider(key []byte) *BearerProvider {
	return &BearerProvider{key}
}

func (p BearerProvider) Validate(req *http.Request) error {
	header := req.Header.Get("Authorization")
	if header == "" {
		return provider.ErrUnauthorized
	}

	var method string
	if x := strings.Index(header, " "); x > 0 {
		method, header = header[:x], header[x+1:]
	} else {
		return provider.ErrMalformed
	}
	if !strings.EqualFold(method, methodBearer) {
		return provider.ErrUnsupported
	}

	parts := strings.SplitN(header, sep, 2) // <signature>$<base64 encoded message>
	if len(parts) != 2 {
		return provider.ErrMalformed
	}

	var policies []acl.ResourcePolicy
	err := crypto.VerifyMessage(p.key, crypto.SHA256, parts[0], &policies, parts[1])
	if err == crypto.ErrSignatureNotValid {
		return provider.ErrUnauthorized
	} else if err != nil {
		return provider.ErrMalformed
	}

	for _, e := range policies {
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

func (p BearerProvider) String() string {
	return "Bearer"
}
